/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package steamid

import (
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/steamid"
)

func Loader(L *lua.LState) int {
  steamid.RegisterType(L)
  L.Push(L.NewFunction(steamid.Constructor))
  return 1
}