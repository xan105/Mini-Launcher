/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package regedit

import (
  "github.com/yuin/gopher-lua"
  "launcher/internal/regedit"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "QueryStringValue": QueryStringValue,
    "WriteStringValue": WriteStringValue,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func QueryStringValue(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  key  := L.ToString(3)  

  value:= regedit.QueryStringValue(root, path, key)          
  L.Push(lua.LString(value))
  return 1
}

func WriteStringValue(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)
  key   := L.ToString(3)
  value := L.ToString(4)

  regedit.WriteStringValue(root, path, key, value)
    
  return 0
}