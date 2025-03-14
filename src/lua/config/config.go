/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package config

import (
  "github.com/yuin/gopher-lua"
)

func ToLuaValue(L *lua.LState, value interface{}) lua.LValue {
  switch v := value.(type) {
    case string:
      return lua.LString(v)
    case bool:
      return lua.LBool(v)
    case float64:
      return lua.LNumber(v)
    default:
      return lua.LNil
  }
}

func ToLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
  table := L.NewTable()
  for key, value := range data {
    switch v := value.(type) {
    case map[string]interface{}:
      table.RawSetString(key, ToLuaTable(L, v))
    case []interface{}:
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

func ToGoMap(luaTable *lua.LTable) map[string]interface{} {
  data := make(map[string]interface{})
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