/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package global

import (
  "log/slog"
  "strings"
  "github.com/yuin/gopher-lua"
)

func format(L *lua.LState, val lua.LValue, depth int) string {
  switch v := val.(type) {
  case *lua.LTable:
    msg := []string{"{"}
    indent := strings.Repeat("  ", depth + 1)
    L.ForEach(v, func(key lua.LValue, value lua.LValue) {
      msg = append(msg, indent + key.String() + ": "+ format(L, value, depth + 1) + ",")
    })
    last := len(msg)-1
    msg[last] = strings.TrimRight(msg[last], ",")
    msg = append(msg, strings.Repeat("  ", depth) + "}")
    return strings.Join(msg, "\n")
  default:
    return val.String()
  }
}

func Log(L *lua.LState) int {
  val := L.CheckAny(1)
  slog.Info(format(L, val, 0))
  return 0
}

func Warn(L *lua.LState) int {
  val := L.CheckAny(1)
  slog.Warn(format(L, val, 0))
  return 0
}

func Error(L *lua.LState) int {
  val := L.CheckAny(1)
  slog.Error(format(L, val, 0))
  return 0
}