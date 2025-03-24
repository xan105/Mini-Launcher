/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package process

import (
  "os"
  "runtime"
  "github.com/yuin/gopher-lua"
)

var EventRegistry = make(map[string]*lua.LFunction)

func Loader(L *lua.LState) int {

  mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "Cwd": Getwd,
    "ExecPath": ExecPath,
    "On": On,
  })
  
  L.SetField(mod, "platform", lua.LString(runtime.GOOS))
  L.SetField(mod, "arch", lua.LString(runtime.GOARCH))
  L.SetField(mod, "pid", lua.LNumber(os.Getpid()))
  L.Push(mod)
  return 1
}

func On(L *lua.LState) int {
  name     := L.CheckString(1)
  callback := L.CheckFunction(2)

  EventRegistry[name] = callback
  return 0
}

func Getwd(L *lua.LState) int {
  cwd, _ := os.Getwd()
  L.Push(lua.LString(cwd))
  return 1
}

func ExecPath(L *lua.LState) int {
  path, _ := os.Executable()
  L.Push(lua.LString(path))
  return 1
}