/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package time

import (
  "time"
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/failure"
)

func Loader(L *lua.LState) int {
  mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "ToUnix": ToUnix,
    "ToIso8601": ToIso8601,
  })
  L.Push(mod)
  return 1
}

func ToUnix(L *lua.LState) int {
  timestamp := L.CheckString(1)

  t, err := time.Parse(time.RFC3339, timestamp)
  if err != nil {
    L.Push(lua.LNumber(0))
    L.Push(failure.LValue(L, "ERR_TIME_CONVERSION", err.Error()))
    return 1
  }

  L.Push(lua.LNumber(t.Unix()))
  return 1
}

func ToIso8601(L *lua.LState) int {
  unixTime := L.CheckInt64(1)

  t := time.Unix(unixTime, 0).UTC()
  iso8601 := t.Format(time.RFC3339)

  L.Push(lua.LString(iso8601))
  return 1
}