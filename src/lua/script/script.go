/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package script

import (
  "embed"
  "path/filepath"
  "strings"
  "github.com/yuin/gopher-lua"
)

//go:embed lua_modules/*.lua
var luaFS embed.FS

func loader(src string) lua.LGFunction {
  return func(L *lua.LState) int {
    if err := L.DoString(src); err != nil {
      L.RaiseError(err.Error())
    }
    return 1
  }
}

func ImportEmbeddedLuaScript(L *lua.LState) error {
  entries, err := luaFS.ReadDir("lua_modules")
  if err != nil {
    return err
  }
  
  for _, entry := range entries {
    if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".lua") {
      continue
    }

    ext := filepath.Ext(entry.Name())
    name := strings.TrimSuffix(entry.Name(), ext)

    data, err := luaFS.ReadFile("lua_modules/" + entry.Name())
    if err != nil {
      return err
    }
    L.PreloadModule(name, loader(string(data)))
  }
  
  return nil
}