/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package process

import (
  "os"
  "runtime"
  "github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {

  mod := L.NewTable()
  L.SetField(mod, "platform", lua.LString(runtime.GOOS))
  L.SetField(mod, "arch", lua.LString(runtime.GOARCH))
  L.SetField(mod, "pid", lua.LNumber(os.Getpid()))
  L.Push(mod)
  return 1
}