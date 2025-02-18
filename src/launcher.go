/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "path/filepath"
  "strings"
  "os/exec"
  "syscall"
  "errors"
  "time"
  "math/rand"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/hook"
  "launcher/internal/splash"
)

func main(){

  cmdLine := parseArgs()

  config, err := fs.ReadJSON[Config](fs.Resolve(cmdLine.ConfigPath))
  if err != nil { panic("JSON Parser", err.Error()) }
  
  binary := fs.Resolve(expand.ExpandVariables(config.Bin))

  //File Integrity
  if config.Integrity != nil && len(config.Integrity) > 0 {
    for _, file := range config.Integrity {
      
      target := binary
      if len(file.Path) > 0 {
        target = fs.Resolve(expand.ExpandVariables(file.Path))
      }
      
      stats, err := os.Stat(target)
      if err != nil { 
        if errors.Is(err, os.ErrNotExist) {
          panic("Integrity failure", "File does not exist: \"" + target + "\"") 
        }
        panic("Integrity failure", err.Error())  
      }
      if file.Size > 0 {
        if stats.Size() != file.Size { 
          panic("Integrity failure", "Size mismatch: \"" + target + "\"") 
        }
      }
      
      SRI := strings.SplitN(file.SRI, "-", 2)
      if len(SRI) != 2 {
        panic("Integrity failure", "Unexpected SRI format: \"" + file.SRI + "\"")
      }
      
      algo, expected := SRI[0], SRI[1]
      sum, err := fs.CheckSum(target, algo)
      if err != nil { panic("Integrity failure", err.Error()) }
      if sum != expected { 
        panic("Integrity failure", 
              "Hash mismatch: \"" + target + "\"\n" +
              "SRI: " + algo + "-" + sum)
      }
    }
  }
  
  cmd := exec.Command(binary)
  argv := []string{ "\"" + binary + "\"" }
  if len(config.Args) > 0 {
    argv = append(argv, expand.ExpandVariables(config.Args))
  }
  cmd.SysProcAttr = &syscall.SysProcAttr{ 
    CmdLine: strings.Join(argv, " "), //verbatim arguments
    HideWindow: config.Hide,
  }

  cmd.Dir = filepath.Dir(binary)
  if len(config.Cwd) > 0 {
    cmd.Dir = fs.Resolve(config.Cwd)
  }

  cmd.Env = os.Environ()
  if len(config.Env) > 0 {
    for key, value := range config.Env {
      if len(key) == 0 || len(value) == 0 {
        continue
      }
      cmd.Env = append(cmd.Env, key + "=" + expand.ExpandVariables(value))
    }
  }

  //Symlink
  if config.Symlink != nil && len(config.Symlink) > 0 {
    for _, link := range config.Symlink {
      if len(link.Origin) == 0 || len(link.Destination) == 0 {
        continue
      }
      origin := fs.Resolve(expand.ExpandVariables(link.Origin))
      destination := fs.Resolve(expand.ExpandVariables(link.Destination))
      err := fs.CreateFolderSymlink(origin, destination)
      if err != nil { 
        panic("Symlink", err.Error()) 
      }
    }
  } 

  //Lua Scripting
  if len(config.Script) > 0 {
    file := fs.Resolve(expand.ExpandVariables(config.Script))
    loadLua(file)
    //Add WASI as well later on (?)
  }

  cmd.Stdin = nil
  cmd.Stdout = nil
  cmd.Stderr = nil

  if cmdLine.DryRun { os.Exit(0) }
  err = cmd.Start()
  if err != nil { panic("Launcher", err.Error()) }
  
  //Addons
  if config.Addons != nil && len(config.Addons) > 0 {
    for _, addon := range config.Addons {
      if len(addon.Path) == 0 {
        continue
      }
      dylib := fs.Resolve(expand.ExpandVariables(addon.Path))
      if fs.FileExist(dylib){
        err = hook.CreateRemoteThread(cmd.Process.Pid, dylib)
        if err != nil {
          if addon.Required {
            cmd.Process.Kill()
            panic("Remote Thread", err.Error())
          }
        }
      }
    }
  }
  
  //Splash screen
  if config.Splash.Show && config.Splash.Images != nil && len(config.Splash.Images) > 0 {
    image := config.Splash.Images[rand.Intn(len(config.Splash.Images))]
    if len(image) > 0 {
      image = fs.Resolve(expand.ExpandVariables(image))

      var timeout uint = 10
      if config.Splash.Timeout > 0 {
        timeout = config.Splash.Timeout
      }
    
      exit := make(chan bool)
      go splash.CreateWindow(exit, cmd.Process.Pid, image)

      select {
        case <-exit:
          return
        case <-time.After(time.Second * time.Duration(timeout)):
          return
      }
    }
  }
}