/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package shell

import (
  "os"
  "os/exec"
  "syscall"
  "bytes"
  "path/filepath"
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/failure"
)

func Loader(L *lua.LState) int {
  mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "Run": Run,
  })
  L.Push(mod)
  return 1
}

func Run(L *lua.LState) int {
  command := L.CheckString(1)
  
  shell := os.Getenv("COMSPEC")
  if len(shell) == 0 {
    shell = filepath.Join(os.Getenv("WINDIR") + "System32/cmd.exe")
  }
  cmd := exec.Command(shell, "/C", command)
  cmd.SysProcAttr = &syscall.SysProcAttr{ HideWindow: true }
  
  stdout := bytes.Buffer{}
  stderr := bytes.Buffer{}
  cmd.Stdout = &stdout
  cmd.Stderr = &stderr
  
  result := L.NewTable()
  if err := cmd.Start(); err != nil {
    L.Push(result)
    L.Push(failure.LValue(L, "ERR_SPAWN_PROCESS", err.Error()))
    return 2
  }
  
  wait := make(chan error)
  go func(){
    wait <- cmd.Wait()
  }()
  <- wait
  
  L.SetField(result, "stdout", lua.LString(stdout.String()))
  L.SetField(result, "stderr", lua.LString(stderr.String()))
  L.Push(result)
  return 1
}