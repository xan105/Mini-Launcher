/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package lua

import (
  "github.com/yuin/gopher-lua"
  "launcher/lua/global"
  "launcher/lua/global/array"
  "launcher/lua/regedit"
  "launcher/lua/random"
  "launcher/lua/file"
  "launcher/lua/archive"
  "launcher/lua/user"
  "launcher/lua/video"
  "launcher/lua/http"
  "launcher/lua/config/json"
  "launcher/lua/config/ini"
  "launcher/lua/config/toml"
  "launcher/lua/config/yaml"
  "launcher/lua/config/xml"
  "launcher/lua/process"
)

type Permissions struct {
  Fs    bool  //Filesystem
  Net   bool  //Network request
  Reg   bool  //Windows registry
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
    {lua.LoadLibName, lua.OpenPackage}, //Must be first
    {lua.BaseLibName, lua.OpenBase},
    {lua.TabLibName, lua.OpenTable},
    {lua.StringLibName, lua.OpenString},
    {lua.MathLibName, lua.OpenMath},
  } {
    if err := L.CallByParam(lua.P{
      Fn:      L.NewFunction(builtin.function),
      NRet:    0,
      Protect: true,
    }, lua.LString(builtin.name)); err != nil {
      return err
    }
  }
  
  //Globals
  L.SetGlobal("sleep", L.NewFunction(global.Sleep))
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
    L.PreloadModule("file", permissionStub
    L.PreloadModule("archive", permissionStub)
  }
  if perm.Net {
    L.PreloadModule("http", http.Loader)
  } else {
    L.PreloadModule("http", permissionStub)
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