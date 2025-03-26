/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package lua

import (
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/failure"
  "launcher/lua/global"
  "launcher/lua/global/array"
  "launcher/lua/module/regedit"
  "launcher/lua/module/random"
  "launcher/lua/module/file"
  "launcher/lua/module/archive"
  "launcher/lua/module/user"
  "launcher/lua/module/video"
  "launcher/lua/module/http"
  "launcher/lua/module/config/json"
  "launcher/lua/module/config/ini"
  "launcher/lua/module/config/toml"
  "launcher/lua/module/config/yaml"
  "launcher/lua/module/config/xml"
  "launcher/lua/module/process"
  "launcher/lua/module/shell"
  "launcher/lua/module/time"
)

type Permissions struct {
  Fs    bool  //Filesystem
  Net   bool  //Network request
  Reg   bool  //Windows registry
  Exec  bool  //Exec shell command
}

var L *lua.LState

var EventRegistry = map[string]map[string]*lua.LFunction{ 
  "process": process.EventRegistry,
}

func LoadLua(filePath string, perm Permissions) error {

  if L != nil { return nil }
  
  L = lua.NewState(lua.Options{ SkipOpenLibs: true })
  
  //Opening a subset of built-in modules
  for _, builtin := range []struct {
    name string
    function lua.LGFunction
  }{
    { lua.LoadLibName, lua.OpenPackage }, //Must be first
    { lua.BaseLibName, lua.OpenBase },
    { lua.TabLibName, lua.OpenTable },
    { lua.StringLibName, lua.OpenString },
    { lua.MathLibName, lua.OpenMath },
    { lua.CoroutineLibName, lua.OpenCoroutine },
  } {
    if err := L.CallByParam(lua.P{
      Fn:      L.NewFunction(builtin.function),
      NRet:    0,
      Protect: true,
    }, lua.LString(builtin.name)); err != nil {
      return err
    }
  }
  
  //Custom Type
  failure.RegisterType(L)
  
  //Globals
  L.SetGlobal("sleep", L.NewFunction(global.Sleep))
  L.SetGlobal("print", L.NewFunction(global.Log)) //override built-in and alias it to console.log
  L.SetGlobal("console", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "log": global.Log,
    "warn": global.Warn,
    "error": global.Error,
  }))
  L.SetGlobal("Array", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "find": array.Find,
    "some": array.Some,
    "includes": array.Includes,
  }))

  //Module
  if perm.Reg {
    L.PreloadModule("regedit", regedit.Loader)
  } else {
    L.PreloadModule("regedit", permissionStub)
  }
  if perm.Fs {
    L.PreloadModule("file", file.Loader)
    L.PreloadModule("archive", archive.Loader)
  } else {
    L.PreloadModule("file", permissionStub)
    L.PreloadModule("archive", permissionStub)
  }
  if perm.Net {
    L.PreloadModule("http", http.Loader)
  } else {
    L.PreloadModule("http", permissionStub)
  }
  if perm.Exec {
    L.PreloadModule("shell", shell.Loader)
  } else {
    L.PreloadModule("shell", permissionStub)
  }
  L.PreloadModule("random", random.Loader)
  L.PreloadModule("user", user.Loader)
  L.PreloadModule("video", video.Loader)
  L.PreloadModule("config/json", json.Loader)
  L.PreloadModule("config/ini", ini.Loader)
  L.PreloadModule("config/toml", toml.Loader)
  L.PreloadModule("config/yaml", yaml.Loader)
  L.PreloadModule("config/xml", xml.Loader)
  L.PreloadModule("process", process.Loader)
  L.PreloadModule("time", time.Loader)
  
  //Exec
  return L.DoFile(filePath);
}

func CloseLua() {
  if L != nil { 
    L.Close()
  }
}

func TriggerEvent(module string, event string) error {
  if L == nil { return nil }
  
  events, exists := EventRegistry[module]
  if !exists { return nil }
  
  callback, exists := events[event]
  if !exists { return nil }
  
  L.Push(callback)
  return L.PCall(0, 0, nil)
}

func permissionStub(L *lua.LState) int {
  L.RaiseError("Module unavailable due to lack of permission !")
  return 0
}