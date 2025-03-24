/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package json

import (
  "encoding/json"
  "github.com/yuin/gopher-lua"
  "launcher/lua/module/config"
  "launcher/lua/type/failure"
)

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
    L.Push(failure.LValue(L, "ERR_JSON_PARSE", err.Error()))
    return 2
  }

  luaTable := config.ToLuaTable(L, data)
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
  
  data := config.ToGoMap(luaTable)
  jsonBytes, err := json.MarshalIndent(data, "", indent)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_JSON_PARSE", err.Error()))
    return 2
  }

  L.Push(lua.LString(jsonBytes))
  return 1
}