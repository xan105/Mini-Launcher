/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package http

import (
  "github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Fetch": Fetch,
  }
  
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}