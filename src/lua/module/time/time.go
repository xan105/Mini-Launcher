/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package time

import (
  "fmt"
  "time"
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/failure"
)

func Loader(L *lua.LState) int {
  mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
    "Current": Current,
    "HumanizeDuration": HumanizeDuration,
    "ToUnix": ToUnix,
    "ToIso8601": ToIso8601,
  })
  L.Push(mod)
  return 1
}

func Current(L *lua.LState) int {
  L.Push(lua.LNumber(time.Now().Unix()))
  return 1
}

func HumanizeDuration(L *lua.LState) int {
  seconds := L.CheckInt64(1)
  duration := time.Duration(seconds) * time.Second
  switch {
    case duration < time.Minute:
      L.Push(lua.LString(fmt.Sprintf("%d seconds", seconds)))
    case duration < time.Hour:
      L.Push(lua.LString(fmt.Sprintf("%d minutes", seconds/60)))
    case duration < time.Hour*24:
      L.Push(lua.LString(fmt.Sprintf("%d hours", seconds/3600)))
    case duration < time.Hour*24*30:
      L.Push(lua.LString(fmt.Sprintf("%d days", seconds/86400)))
    case duration < time.Hour*24*365:
      L.Push(lua.LString(fmt.Sprintf("%d months", seconds/(86400*30))))
    default:
      L.Push(lua.LString(fmt.Sprintf("%d years", seconds/(86400*365))))
  }
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