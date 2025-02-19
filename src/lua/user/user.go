/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package user

import (
  "os/user"
  "github.com/yuin/gopher-lua"
  "launcher/internal/locale"
)

func getUserName() (string, error) {
  user, err := user.Current()
  if err != nil {
    return "", err
  }
  return user.Username, nil
}

func getUserLang() (string, string, error) {
  code, err := locale.GetUserLocale()
  if err != nil {
    return "", "", err
  }
  lang, err := locale.GetLanguageFromLocale(code)
  if err != nil {
    return code, "", err
  }
  return code, lang, nil
}

func Loader(L *lua.LState) int {
  mod := L.NewTable()
  
  name, _ := getUserName()
  code, lang, _ := getUserLang()
  
  L.SetField(mod, "name", lua.LString(name))
  L.SetField(mod, "locale", lua.LString(code))
  L.SetField(mod, "language", lua.LString(lang))
  L.Push(mod)
  return 1
}