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
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Write(L *lua.LState) int {
  //get argument
  filename  := L.ToString(1)  
  data      := L.ToString(2)
  format    := L.ToString(3)

  if len(format) == 0 {
    format = "utf8"
  } 

  fs.WriteFile(filename, data, format)
    
  return 0 //number of result(s)
}