/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package json

import (
  "encoding/json"
  "github.com/yuin/gopher-lua"
)

func toLuaValue(L *lua.LState, value interface{}) lua.LValue {
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

func toLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
  table := L.NewTable()
  for key, value := range data {
    switch v := value.(type) {
    case map[string]interface{}:
      table.RawSetString(key, toLuaTable(L, v))
    case []interface{}:
      arrayTable := L.NewTable()
      for i, item := range v {
        arrayTable.RawSetInt(i+1, toLuaValue(L, item))
      }
      table.RawSetString(key, arrayTable)
    default:
      table.RawSetString(key, toLuaValue(L, v))
    }
  }
  return table
}

func toGoMap(luaTable *lua.LTable) map[string]interface{} {
  data := make(map[string]interface{})
  luaTable.ForEach(func(key lua.LValue, value lua.LValue) {
    switch v := value.(type) {
    case *lua.LTable:
      data[key.String()] = toGoMap(v)
    default:
      data[key.String()] = v.String()
    }
  })
  return data
}

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Parse": Parse,
    "Stringify": Stringify,
  }
  
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Parse(L *lua.LState) int {
  jsonStr := L.CheckString(1)

  var data map[string]interface{}
  if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  luaTable := toLuaTable(L, data)
  L.Push(luaTable)
  return 1
}

func Stringify(L *lua.LState) int {
  luaTable := L.CheckTable(1)
  pretty := true
  if L.GetTop() > 1 {
    pretty = L.CheckBool(2)
  }
  
  indent := ""
  if pretty {
    indent = "  "
  }
  
  data := toGoMap(luaTable)
  jsonBytes, err := json.MarshalIndent(data, "", indent)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  L.Push(lua.LString(jsonBytes))
  return 1
}