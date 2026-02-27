/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package ini

import (
  "github.com/yuin/gopher-lua"
  "launcher/internal/ini"
  "launcher/lua/util"
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
  iniStr := L.CheckString(1)

  options := ini.ParserOptions{
    Filter:  []string{},
    Global:  true,
    Unquote: true,
    Boolean: true,
    Number:  true,
  }

  if L.GetTop() >= 2 {
    table := L.CheckTable(2)
    
    table.ForEach(func(key lua.LValue, value lua.LValue) {
      switch key.String() {
      case "filter":
        if arr, ok := value.(*lua.LTable); ok {
          filter := []string{}
          arr.ForEach(func(_, v lua.LValue) {
            filter = append(filter, v.String())
          })
          options.Filter = filter
        }
      case "global":
        if b, ok := value.(lua.LBool); ok {
          options.Global = bool(b)
        }
      case "unquote":
        if b, ok := value.(lua.LBool); ok {
          options.Unquote = bool(b)
        }
      case "boolean":
        if b, ok := value.(lua.LBool); ok {
          options.Boolean = bool(b)
        }
      case "number":
        if b, ok := value.(lua.LBool); ok {
          options.Number = bool(b)
        }
      }
    })
  }

  data := ini.Parse(iniStr, &options)
  luaTable := util.ToLuaTable(L, data)
  L.Push(luaTable)
  return 1
}

func Stringify(L *lua.LState) int {
  luaTable := L.CheckTable(1)
  
  options := ini.StringifyOptions{
    Whitespace: true,
    BlankLine: false,
    Quote: false,
  }
  
  if L.GetTop() >= 2 {
    table := L.CheckTable(2)
    
    table.ForEach(func(key lua.LValue, value lua.LValue) {
      switch key.String() {
      case "whitespace":
        if b, ok := value.(lua.LBool); ok {
          options.Whitespace = bool(b)
        }
      case "blankLine":
        if b, ok := value.(lua.LBool); ok {
          options.BlankLine = bool(b)
        }
      case "quote":
        if b, ok := value.(lua.LBool); ok {
          options.Quote = bool(b)
        }
      }
    })
  }

  data := util.ToGoMap(luaTable)
  iniStr := ini.Stringify(data, &options)
  L.Push(lua.LString(iniStr))
  return 1
}