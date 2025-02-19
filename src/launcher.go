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
  "launcher/internal/fs"
  "launcher/internal/expand"
)

func buildCommand(binary string, args string, cwd string, env map[string]string, hide bool) *exec.Cmd {
  cmd := exec.Command(binary)
  
  argv := []string{ "\"" + binary + "\"" }
  if len(args) > 0 {
    argv = append(argv, expand.ExpandVariables(args))
  }
  cmd.SysProcAttr = &syscall.SysProcAttr{ 
    CmdLine: strings.Join(argv, " "), //verbatim arguments
    HideWindow: hide,
  }

  cmd.Dir = filepath.Dir(binary)
  if len(cwd) > 0 {
    cmd.Dir = fs.Resolve(cwd)
  }

  cmd.Env = os.Environ()
  if len(env) > 0 {
    for key, value := range env {
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
  config, err := fs.ReadJSON[Config](fs.Resolve(cmdLine.ConfigPath))
  if err != nil { panic("JSON Parser", err.Error()) }
  
  binary := fs.Resolve(expand.ExpandVariables(config.Bin))
  cmd := buildCommand(binary, config.Args, config.Cwd, config.Env, config.Hide)
  
  verifyIntegrity(binary, config.Integrity)
  setCompatFlags(binary, config.Compatibility)
  makeLink(config.Symlink)

  if len(config.Script) > 0 {
    loadLua(fs.Resolve(expand.ExpandVariables(config.Script)))
  }

  if cmdLine.DryRun { os.Exit(0) }
  err = cmd.Start()
  if err != nil { panic("Launcher", err.Error()) }
  
  inject(cmd.Process, config.Addons)
  displaySplash(cmd.Process.Pid, config.Splash)
}