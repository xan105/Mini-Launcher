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

type Addon struct {
  Path      string                `json:"path"`
  Required  bool                  `json:"required"`
}

type File struct {
  Path      string                `json:"file"`
  SRI       string                `json:"sri"` 
  Size      int64                 `json:"size"`
}

type Splash struct {
  Show        bool                `json:"show"`
  Images      []string            `json:"image"`
  Timeout     uint                `json:"timeout"`
}

type Link struct {
  Origin      string              `json:"from"`
  Destination string              `json:"to"`
}

type Config struct {
  Bin         string              `json:"bin"`
  Cwd         string              `json:"cwd"`
  Args        string              `json:"args"`
  Env         map[string]string   `json:"env"`
  Hide        bool                `json:"hide"`
  Script      string              `json:"script"`
  Addons      []Addon             `json:"addons"`
  Integrity   []File              `json:"integrity"`
  Splash      Splash              `json:"splash"`
  Symlink     []Link              `json:"symlink"`
}

func main(){

  args := parseArgs()

  cwd, err := os.Getwd() 
  if err != nil { panic(err.Error()) }

  configFile := filepath.FromSlash(args.Config)
  if !filepath.IsAbs(args.Config) {
    configFile = filepath.Join(cwd, args.Config)
  }

  config, err := fs.ReadJSON[Config](configFile)
  if err != nil { panic("Parsing JSON failure!\n\n" + err.Error()) }
  
  binary := filepath.FromSlash(expand.ExpandVariables(config.Bin))
  if !filepath.IsAbs(binary) {
    binary = filepath.Join(cwd, binary)
  }
  
  if config.Integrity != nil && len(config.Integrity) > 0 {
    for _, file := range config.Integrity {
      
      target:= binary
      if len(file.Path) > 0 {
        target = filepath.FromSlash(expand.ExpandVariables(file.Path))
        if !filepath.IsAbs(target) {
          target = filepath.Join(cwd, target)
        } 
      }
      
      stats, err := os.Stat(target)
      if err != nil { 
        if errors.Is(err, os.ErrNotExist) {
          panic("Integrity failure!\n\n" + 
                "File does not exist: \"" + target + "\"") 
        }
        panic(err.Error())  
      }
      if file.Size > 0 {
        if stats.Size() != file.Size { 
          panic("Integrity failure!\n\n" + 
                "Size mismatch: \"" + target + "\"") 
        }
      }
      
      SRI := strings.SplitN(file.SRI, "-", 2)
      if len(SRI) != 2 {
        panic("Integrity failure!\n\n" + 
              "Unexpected SRI format: \"" + file.SRI + "\"") 
      }
      
      algo, expected := SRI[0], SRI[1]
      sum, err := fs.CheckSum(target, algo)
      if err != nil { panic(err.Error()) }
      if sum != expected { 
        panic("Integrity failure!\n\n" + 
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
  cmd.SysProcAttr = &syscall.SysProcAttr{ CmdLine: strings.Join(argv, " ") } //verbatim arguments

  if config.Hide == true {
    cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  }

  cmd.Dir = filepath.Dir(binary)
  if len(config.Cwd) > 0 {
    cmd.Dir = filepath.FromSlash(config.Cwd)
  }

  cmd.Env = os.Environ()
  if len(config.Env) > 0 {
    for key, value := range config.Env {
      cmd.Env = append(cmd.Env, key + "=" + expand.ExpandVariables(value))
    }
  }

  if config.Symlink != nil && len(config.Symlink) > 0 {
    for _, link := range config.Symlink {
      origin := filepath.FromSlash(expand.ExpandVariables(link.Origin))
      if !filepath.IsAbs(origin) {
        origin = filepath.Join(cwd, origin)
      }
      
      destination := filepath.FromSlash(expand.ExpandVariables(link.Destination))
      if !filepath.IsAbs(destination) {
        destination = filepath.Join(cwd, destination)
      }

      err:= fs.CreateFolderSymlink(origin, destination)
      if err != nil { 
        panic("Symlink failure!\n\n" + err.Error()) 
      }
    }
  } 

  if len(config.Script) > 0 {
    file := filepath.FromSlash(config.Script)
    if !filepath.IsAbs(file) {
      file = filepath.Join(cwd, file)
    }
    loadLua(file)
    //Add WASI as well later on (?)
  }

  cmd.Stdin = nil
  cmd.Stdout = nil
  cmd.Stderr = nil

  if (args.DryRun) {
    os.Exit(0)
  }
  
  err = cmd.Start()
    if err != nil { panic(err.Error()) }
  
  //Addons
  if config.Addons != nil && len(config.Addons) > 0 {
    for _, addon := range config.Addons {
          
      dylib := filepath.FromSlash(addon.Path)
      if !filepath.IsAbs(dylib) {
        dylib = filepath.Join(cwd, dylib)
      }
            
      if fs.FileExist(dylib){
        err = hook.CreateRemoteThread(uintptr(cmd.Process.Pid), dylib)
        if err != nil {
          if addon.Required {
            cmd.Process.Kill()
            panic(err.Error())
          }
        }
      }
    }
  }
  
  //splash screen
  if config.Splash.Show {
    image := config.Splash.Images[rand.Intn(len(config.Splash.Images))]
    image = filepath.FromSlash(expand.ExpandVariables(image))
    if !filepath.IsAbs(image) {
      image = filepath.Join(cwd, image)
    }
  
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