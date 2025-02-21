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

func Log(L *lua.LState) int {
  val := L.CheckAny(1)

  switch v := val.(type) {
  case *lua.LTable:
    msg := []string{"{"}
    L.ForEach(v, func(key lua.LValue, value lua.LValue) {
      msg = append(msg, "  " + value.String() + ",")
    })
    msg = append(msg, "}")
    slog.Info(strings.Join(msg, "\n"))
  default:
    slog.Info(val.String()) 
  }

  return 0
}

func Warn(L *lua.LState) int {
  val := L.CheckAny(1)

  switch v := val.(type) {
  case *lua.LTable:
    msg := []string{"{"}
    L.ForEach(v, func(key lua.LValue, value lua.LValue) {
      msg = append(msg, "  " + value.String() + ",")
    })
    msg = append(msg, "}")
    slog.Warn(strings.Join(msg, "\n"))
  default:
    slog.Warn(val.String()) 
  }

  return 0
}

func Error(L *lua.LState) int {
  val := L.CheckAny(1)

  switch v := val.(type) {
  case *lua.LTable:
    msg := []string{"{"}
    L.ForEach(v, func(key lua.LValue, value lua.LValue) {
      msg = append(msg, "  " + value.String() + ",")
    })
    msg = append(msg, "}")
    slog.Error(strings.Join(msg, "\n"))
  default:
    slog.Error(val.String()) 
  }

  return 0
}