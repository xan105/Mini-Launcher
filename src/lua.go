/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import (
  "github.com/yuin/gopher-lua"
  "launcher/lua/global"
  "launcher/lua/regedit"
  "launcher/lua/random"
  "launcher/lua/file"
  "launcher/lua/user"
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
  
  //Module
  L.PreloadModule("regedit", regedit.Loader)
  L.PreloadModule("random", random.Loader)
  L.PreloadModule("file", file.Loader)
  L.PreloadModule("user", user.Loader)
  
  //Exec
  if err := L.DoFile(filePath); err != nil {
    panic("Lua", err.Error())
  }
}