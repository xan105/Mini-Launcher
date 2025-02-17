/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package expand

import(
  "os"
  "os/user"
  "strings"
  "regexp"
  "path/filepath"
  "launcher/internal/regedit"
)

func getShellFolder(keys []string, root string) string {
  const path = "Software/Microsoft/Windows/CurrentVersion/Explorer/User Shell Folders"
  for _, key := range keys {
    value := regedit.QueryStringValue(root, path, key)
    if len(value) > 0 {
      return value
    }
  }
  return ""
}

func getUserShellFolder(keys []string) string {
  return getShellFolder(keys, "HKCU")
}

func getCommonShellFolder(keys []string) string {
  return getShellFolder(keys, "HKLM")
}

func ExpandVariables(input string) string {
  re := regexp.MustCompile(`%([^%]+)%`)
  
  return re.ReplaceAllStringFunc(input, func(match string) string {
    variable := strings.Trim(match, "%")
    switch variable { 
      case "APPDATA": {
        value:= getUserShellFolder([]string{"AppData"})
        if len(value) > 0 {
          return value
        }
        env:= os.Getenv("APPDATA")
        if len(env) > 0 {
          return env
        }
        return match
      }
      case "LOCALAPPDATA": {
        value:= getUserShellFolder([]string{"Local AppData"})
        if len(value) > 0 {
          return value
        }
        env:= os.Getenv("LOCALAPPDATA")
        if len(env) > 0 {
          return env
        }
        return match
      }
      case "PROGRAMDATA": {
        value:= getCommonShellFolder([]string{"Common AppData"})
        if len(value) > 0 {
          return value
        }
        env:= os.Getenv("PROGRAMDATA")
        if len(env) > 0 {
          return env
        }
        return match
      }
      case "DESKTOP": {
        value:= getUserShellFolder([]string{"Desktop"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
        }
        return match
      }
      case "DOCUMENTS": {
        value:= getUserShellFolder([]string{
                "{F42EE2D3-909F-4907-8871-4C22FC0BF756}", //win10
                "Personal"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Documents")
        }
        return match
      }
      case "MUSIC": {
        value:= getUserShellFolder([]string{
                "{A0C69A99-21C8-4671-8703-7934162FCF1D}", //win10
                "My Music"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Music")
        }
        return match
      }
      case "PICTURES": {
        value:= getUserShellFolder([]string{
                "{0DDD015D-B06C-45D5-8C4C-F59713854639}", //win10
                "My Pictures"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Pictures")
        }
        return match
      }
      case "VIDEOS": {
        value:= getUserShellFolder([]string{
                "{35286A68-3C57-41A1-BBB1-0EAE73D76C95}", //win10
                "My Video"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Videos")
        }
        return match
      }
      case "DOWNLOAD": {
        value:= getUserShellFolder([]string{
                "{7D83EE9B-2244-4E70-B1F5-5393042AF1E4}", //win10
                "{374DE290-123F-4565-9164-39C4925E467B}"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
        }
        return match
      }
      case "SAVEGAME": {
        value:= getUserShellFolder([]string{"{4C5C32FF-BB9D-43b0-B5B4-2D72E54EAAA4}"})
        if len(value) > 0 {
          return value
        }
        profile := os.Getenv("USERPROFILE")
        if len(profile) > 0 {
          return filepath.Join(os.Getenv("USERPROFILE"), "Saved Games")
        }
        return match
      }
      case "HOMEDIR":
      case "USERPROFILE": {
        value:= os.Getenv("USERPROFILE")
        if len(value) > 0 {
          return value
        }
        return match
      }
      case "PUBLIC": {
        value:= os.Getenv("PUBLIC")
        if len(value) > 0 {
          return value
        }
        return match
      }
      case "SYSTEMDIR":{
        variables := []string{ "SYSTEMROOT", "WINDIR" }
        for _, variable := range variables {
          value := os.Getenv(variable)
          if len(value) > 0 {
            return value
          }
        }
        return match
      }
      case "TEMP":
      case "TMP": {
        variables := []string{"TEMP", "TMP"}
        for _, variable := range variables {
          value := os.Getenv(variable)
          if len(value) > 0 {
            return value
          }
        }
        return match
      }
      case "CURRENTDIR": {
        cwd, err := os.Getwd()
        if err != nil {
          return match
        }
        return cwd
      }
      case "BINDIR": {
        process, err := os.Executable()
        if err != nil {
          return match
        }
        return filepath.Dir(process)
      }
      case "USERNAME": {
        user, err := user.Current()
        if err != nil {
          return match
        }
        return user.Username
      }
    }
    return match
  })
}