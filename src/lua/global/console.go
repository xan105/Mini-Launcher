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

//cf: https://nodejs.org/api/util.html#customizing-utilinspect-colors
const (
  reset     = "\033[0m"
  yellow    = "\033[33m"  // Bigint, Boolean, Number
  magenta   = "\033[35m"  // Date
  underline = "\033[4m"   // Module
  green     = "\033[32m"  // String, Symbol
  cyan      = "\033[36m"  // Special (e.g., Proxies)
  red       = "\033[31m"  // RegExp
  bold      = "\033[1m"   // Null
  grey      = "\033[90m"  // Undefined
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
  case lua.LString:
    return green + "\"" + v.String() + "\"" + reset
  case lua.LNumber:
    return yellow + v.String() + reset
  case lua.LBool:
    return yellow + v.String() + reset
  case *lua.LNilType:
    return bold + "nil" + reset
  case *lua.LFunction, *lua.LState, *lua.LChannel:
    return cyan + v.String() + reset
  default:
    return grey + val.String() + reset
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