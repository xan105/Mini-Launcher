/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package file

import (
  "github.com/yuin/gopher-lua"
  "launcher/internal/fs"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Write": Write,
    "Read": Read,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Write(L *lua.LState) int {
  filename := L.ToString(1)  
  data     := L.ToString(2)
  format   := L.ToString(3)

  if len(format) == 0 {
    format = "utf8"
  } 

  err := fs.WriteFile(filename, data, format)
  if err != nil {
    L.RaiseError(err.Error());
  }
    
  return 0
}

func Read(L *lua.LState) int {
  filename := L.ToString(1)  
  format   := L.ToString(2)

  if len(format) == 0 {
    format = "utf8"
  } 
  
  data, err := fs.ReadFile(filename, format)
  if err != nil {
    L.RaiseError(err.Error());
  }

  L.Push(lua.LString(data))
  return 1
}