/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package util

import (
  "github.com/yuin/gopher-lua"
)

func ToLuaValue(L *lua.LState, value any) lua.LValue {
// All Go numeric types are converted to lua.LNumber (float64)
// Large integers may lose precision
  switch v := value.(type) {
    case string:
      return lua.LString(v)
    case bool:
      return lua.LBool(v)
    case int:
      return lua.LNumber(v)
    case int8:
      return lua.LNumber(v)
    case int16:
      return lua.LNumber(v)
    case int32:
      return lua.LNumber(v)
    case int64:
      return lua.LNumber(v)
    case uint:
      return lua.LNumber(v)
    case uint8:
      return lua.LNumber(v)
    case uint16:
      return lua.LNumber(v)
    case uint32:
      return lua.LNumber(v)
    case uint64:
      return lua.LNumber(v)
    case float32:
      return lua.LNumber(v)
    case float64:
      return lua.LNumber(v)
    case map[string]any:
      return ToLuaTable(L, v)
    case []any:
      arrayTable := L.NewTable()
      for i, item := range v {
        arrayTable.RawSetInt(i+1, ToLuaValue(L, item))
      }
    return arrayTable
    default:
      return lua.LNil
  }
}

func ToLuaTable(L *lua.LState, data map[string]any) *lua.LTable {
  table := L.NewTable()
  for key, value := range data {
    switch v := value.(type) {
      case map[string]any:
        table.RawSetString(key, ToLuaTable(L, v))
      case []any:
        arrayTable := L.NewTable()
        for i, item := range v {
          arrayTable.RawSetInt(i+1, ToLuaValue(L, item))
        }
        table.RawSetString(key, arrayTable)
      default:
        table.RawSetString(key, ToLuaValue(L, v))
    }
  }
  return table
}

func ToGoMap(luaTable *lua.LTable) map[string]any {
  data := make(map[string]any)
  luaTable.ForEach(func(key lua.LValue, value lua.LValue) {
    switch v := value.(type) {
    case *lua.LTable:
      data[key.String()] = ToGoMap(v)
    case lua.LString:
        data[key.String()] = v.String()
    case lua.LNumber:
        data[key.String()] = float64(v)
    case lua.LBool:
        data[key.String()] = bool(v)
    default:
      data[key.String()] = v.String()
    }
  })
  return data
}