/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package video

import (
  "launcher/internal/video"
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/failure"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Current": Current,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Current(L *lua.LState) int {

  display, err := video.GetCurrentDisplayMode()
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_WIN32_API", err.Error()))
    return 2
  }

  displayMode := L.NewTable()
  L.SetField(displayMode, "width", lua.LNumber(display.Width))
  L.SetField(displayMode, "height", lua.LNumber(display.Height))
  L.SetField(displayMode, "hz", lua.LNumber(display.Hz))
  L.SetField(displayMode, "scale", lua.LNumber(display.Scale))
  L.Push(displayMode)
  return 1
}