/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package user

import (
  "os/user"
  "errors"
  "strings"
  "github.com/yuin/gopher-lua"
  "launcher/internal/locale"
)

func getUserName() (string, error) {
  user, err := user.Current()
  if err != nil {
    return "", err
  }
  parts := strings.Split(user.Username, "\\")
  last := len(parts)-1
  return parts[last], nil
}

func getUserLang() (string, string, string, error) {
  localeName, err := locale.GetUserLocale()
  if err != nil {
     return "", "", "", err
  }
  
  if !strings.Contains(localeName, "-") {
     return "", "", "", errors.New("Unexpected local ISO 639: \"" + localeName + "\"")
  }
  loc := strings.SplitN(localeName, "-", 2)
  if len(loc) != 2 { 
    return "", "", "", errors.New("Unexpected local ISO 639: \"" + localeName + "\"") 
  }
  
  code, region := loc[0], loc[1]
  lang, err := locale.GetLanguageFromLocale(localeName)
  if err != nil {
    return code, region, "", err
  }
  
  return code, region, lang, nil
}

func Loader(L *lua.LState) int {
  name, _ := getUserName()
  code, region, lang, _ := getUserLang()
  
  locale := L.NewTable()
  L.SetField(locale, "code", lua.LString(code))
  L.SetField(locale, "region", lua.LString(region))
  
  mod := L.NewTable()
  L.SetField(mod, "name", lua.LString(name))
  L.SetField(mod, "language", lua.LString(lang))
  L.SetField(mod, "locale", locale)
  L.Push(mod)
  return 1
}