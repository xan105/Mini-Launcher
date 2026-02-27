/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package yaml

import (
  "gopkg.in/yaml.v3"
  "github.com/yuin/gopher-lua"
  "launcher/lua/util"
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
  yamlStr := L.CheckString(1)

  var data map[string]any
  if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_YAML_PARSE", err.Error()))
    return 2
  }

  luaTable := util.ToLuaTable(L, data)
  L.Push(luaTable)
  return 1
}

func Stringify(L *lua.LState) int {
  luaTable := L.CheckTable(1)

  data := util.ToGoMap(luaTable)
  yamlBytes, err := yaml.Marshal(data)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_YAML_PARSE", err.Error()))
    return 2
  }

  L.Push(lua.LString(yamlBytes))
  return 1
}