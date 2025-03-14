/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

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

func loadLua(filePath string){
  L := lua.NewState(lua.Options{ SkipOpenLibs: true })
  defer L.Close()
  
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
      panic("Lua", err.Error())
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
  L.PreloadModule("regedit", regedit.Loader)
  L.PreloadModule("random", random.Loader)
  L.PreloadModule("file", file.Loader)
  L.PreloadModule("archive", archive.Loader)
  L.PreloadModule("user", user.Loader)
  L.PreloadModule("video", video.Loader)
  L.PreloadModule("http", http.Loader)
  L.PreloadModule("config/json", json.Loader)
  L.PreloadModule("config/ini", ini.Loader)
  L.PreloadModule("config/toml", toml.Loader)
  L.PreloadModule("config/yaml", yaml.Loader)
  L.PreloadModule("config/xml", xml.Loader)
  L.PreloadModule("process", process.Loader)
  
  //Exec
  if err := L.DoFile(filePath); err != nil {
    panic("Lua", err.Error())
  }
}