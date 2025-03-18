/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "os/exec"
  "strings"
  "syscall"
  "path/filepath"
  "launcher/lua"
  "launcher/internal/fs"
  "launcher/internal/expand"
)

func buildCommand(binary string, config Config) *exec.Cmd {
  
  var cmd *exec.Cmd
  
  if config.Shell {
    shell := os.Getenv("COMSPEC")
    if len(shell) == 0 {
      shell = filepath.Join(os.Getenv("WINDIR") + "System32/cmd.exe")
    }
    cmd = exec.Command(shell)
    argv := []string{ "\"" + shell + "\"" } //argv0
    argv = append(argv, "/D", "/C", "\"\"" + binary + "\"")
    if len(config.Args) > 0 {
      argv = append(argv, expand.ExpandVariables(config.Args))
    }
    last := len(argv)-1
    argv[last] = argv[last] + "\""
    
    cmd.SysProcAttr = &syscall.SysProcAttr{ 
      CmdLine: strings.Join(argv, " "),
      HideWindow: config.Hide,
    }
  } else {
    cmd = exec.Command(binary)
    argv := []string{ "\"" + binary + "\"" } //argv0
    if len(config.Args) > 0 {
      argv = append(argv, expand.ExpandVariables(config.Args))
    }
    cmd.SysProcAttr = &syscall.SysProcAttr{ 
      CmdLine: strings.Join(argv, " "), //verbatim arguments
      HideWindow: config.Hide,
    }
  }
  
  cmd.Dir = filepath.Dir(binary)
  if len(config.Cwd) > 0 {
    cmd.Dir = fs.Resolve(config.Cwd)
  }

  cmd.Env = os.Environ()
  if len(config.Env) > 0 {
    for key, value := range config.Env {
      if len(key) > 0 && len(value) > 0 {
        cmd.Env = append(cmd.Env, key + "=" + expand.ExpandVariables(value))
      }
    }
  }
  
  cmd.Stdin = nil
  cmd.Stdout = nil
  cmd.Stderr = nil
  
  return cmd
}

func main(){

  cmdLine := parseArgs()
  configFile := fs.Resolve(cmdLine.ConfigPath)
  config, err := fs.ReadJSON[Config](configFile)
  if err != nil { panic("JSON Parser", err.Error()) }

  if path := displayMenuOverride(config.Menu, cmdLine.ConfigPath); len(path) > 0 {
    overrideFile := fs.Resolve(path)
    if (configFile != overrideFile) {
      override, err := fs.ReadJSON[Config](overrideFile)
      if err != nil { panic("JSON Parser", err.Error()) }
      mergeConfig(&config, &override)
    }
  }

  binary := fs.Resolve(expand.ExpandVariables(config.Bin))
  cmd := buildCommand(binary, config)
  
  verifyIntegrity(binary, config.Integrity)
  makeLink(config.Symlink)
  applyFileAttributes(config.Attrib)
  setCompatFlags(binary, config.Compatibility)
  updatePrefixSettings(config.Prefix)

  if len(config.Script) > 0 {
    script := fs.Resolve(expand.ExpandVariables(config.Script))
    ext := filepath.Ext(script)
    switch ext {
      case ".lua": {
        if err := lua.LoadLua(script); err != nil {
          panic("Lua", err.Error())
        }
      }
      default: {
        panic("Launcher", "Unsupported script: \""+ ext +"\"")  
      }
    } 
  }

  if cmdLine.DryRun { 
    lua.CloseLua()
    os.Exit(0) 
  }
  
  if err := cmd.Start(); err != nil { 
    panic("Launcher", err.Error()) 
  }
  
  loadAddons(cmd.Process, config.Addons)
  displaySplash(cmd.Process.Pid, config.Splash)
  
  if config.Wait {
    cmd.Wait()
  }
  
  if err := lua.TriggerEvent("process", "will-quit"); err != nil {
    panic("Lua", err.Error())
  }
  
  lua.CloseLua()
}