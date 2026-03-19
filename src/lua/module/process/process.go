/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package process

import (
  "os"
  "os/exec"
  "strings"
  "runtime"
  "path/filepath"
  "launcher/lua/util"
  "launcher/internal/wine"
  "github.com/yuin/gopher-lua"
)

var EventRegistry = map[string][]*lua.LFunction{}

func parseEnviron(L *lua.LState, environ []string) *lua.LTable {
  env := L.NewTable()
  for _, entry := range environ {
    if key, value, ok := strings.Cut(entry, "="); ok && len(key) > 0 {
      L.SetField(env, key, lua.LString(value))
    }
  }
  return env
}

func Loader(L *lua.LState, targetProcess *exec.Cmd, argv []string) int {

  mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "On": On,
  })
  
  L.SetField(mod, "platform", lua.LString(runtime.GOOS))
  L.SetField(mod, "arch", lua.LString(runtime.GOARCH))
  L.SetField(mod, "pid", lua.LNumber(os.Getpid()))
  L.SetField(mod, "wine", lua.LBool(wine.IsWineOrProton()))
  
  execPath, _ := os.Executable()
  cwd, _ := os.Getwd()
  L.SetField(mod, "path", lua.LString(execPath))
  L.SetField(mod, "bin", lua.LString(filepath.Base(execPath)))
  L.SetField(mod, "dir", lua.LString(filepath.Dir(execPath)))
  L.SetField(mod, "cwd", lua.LString(cwd))
  L.SetField(mod, "args", util.ToLuaValue(L, os.Args[1:]))
  L.SetField(mod, "env", parseEnviron(L, os.Environ()))
  
  target := L.NewTable()
  L.SetField(target, "path", lua.LString(targetProcess.Path))
  L.SetField(target, "bin", lua.LString(filepath.Base(targetProcess.Path)))
  L.SetField(target, "dir", lua.LString(filepath.Dir(targetProcess.Path)))
  L.SetField(target, "cwd", lua.LString(targetProcess.Dir))
  L.SetField(target, "argv", util.ToLuaValue(L, argv))
  L.SetField(target, "env", parseEnviron(L, targetProcess.Env))
  L.SetField(mod, "target", target)

  L.Push(mod)
  return 1
}

func On(L *lua.LState) int {
  name     := L.CheckString(1)
  callback := L.CheckFunction(2)

  EventRegistry[name] = append(EventRegistry[name], callback)
  return 0
}