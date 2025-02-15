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

type Config struct {
  Bin         string              `json:"bin"`
  Cwd         string              `json:"cwd"`
  Args        string              `json:"args"`
  Env         map[string]string   `json:"env"`
  Hide        bool                `json:"hide"`
  Keygen      string              `json:"keygen"`
  Addons      []Addon             `json:"addons"`
  Integrity   []File              `json:"integrity"`
  Splash      Splash              `json:"splash"`
}

func main(){

  args := parseArgs()

  cwd, err := os.Getwd() 
  if err != nil { panic(err.Error()) }

  configFile := args.Config
  if !filepath.IsAbs(args.Config) {
    configFile = filepath.Join(cwd, args.Config)
  }

  config, err := readJSON(configFile)
  if err != nil { panic("Parsing JSON failure!\n\n" + err.Error()) }
  
  binary := expand.ExpandVariables(config.Bin)
  if !filepath.IsAbs(config.Bin) {
    binary = filepath.Join(cwd, config.Bin)
  }
  
  if config.Integrity != nil && len(config.Integrity) > 0 {
    for _, file := range config.Integrity {
      
      target:= binary
      if len(file.Path) > 0 {
        target = expand.ExpandVariables(file.Path)
        if !filepath.IsAbs(file.Path) {
          target = filepath.Join(cwd, file.Path)
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
      sum, err := checkSum(target, algo)
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
    cmd.Dir = config.Cwd
  }

  cmd.Env = os.Environ()
  if len(config.Env) > 0 {
    for key, value := range config.Env {
      cmd.Env = append(cmd.Env, key + "=" + expand.ExpandVariables(value))
    }
  }

  if len(config.Keygen) > 0 {
    file := config.Keygen
    if !filepath.IsAbs(file) {
      file = filepath.Join(cwd, file)
    }
    loadLua(file)
    //Add WASI as well later on (?)
  }

  cmd.Stdin = nil
  cmd.Stdout = nil
  cmd.Stderr = nil

  err = cmd.Start()
    if err != nil { panic(err.Error()) }
  
  //Addons
  if config.Addons != nil && len(config.Addons) > 0 {
    for _, addon := range config.Addons {
          
      dylib := addon.Path
      if !filepath.IsAbs(dylib) {
        dylib = filepath.Join(cwd, dylib)
      }
            
      if fileExist(dylib){
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
    image = expand.ExpandVariables(image)
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