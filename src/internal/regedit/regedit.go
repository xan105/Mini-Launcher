/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package regedit

import (
  "path/filepath"
  "golang.org/x/sys/windows/registry"
)

func getRootKey(root string) registry.Key {
  switch root {
    case "HKCU":
      return registry.CURRENT_USER
    case "HKLM":
      return registry.LOCAL_MACHINE
    case "HKU":
      return registry.USERS
    case "HKCC":
      return registry.CURRENT_CONFIG
    case "HKCR":
      return registry.CLASSES_ROOT
    default:
      return 0
  }
}

func QueryStringValue(root string, path string, key string) string { 

  var result string
  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  result, keyType, _ := k.GetStringValue(key)
  
  if keyType == registry.EXPAND_SZ {
    expanded, err := registry.ExpandString(result)
    if err == nil {
      result = expanded
    }
  }
 
  return result
}

func WriteStringValue(root string, path string, key string, value string) {
  
  HKEY := getRootKey(root)
  var buf []byte;
  
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS)
  defer k.Close()
  _, keyType, _ := k.GetValue(key, buf)
  
  if keyType == registry.EXPAND_SZ {
    k.SetExpandStringValue(key, value)
  } else {
    k.SetStringValue(key, value)
  }
}

func DeleteKeyValue (root string, path string, key string) {

  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.DeleteValue(key)
}