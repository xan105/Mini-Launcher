/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//ref: https://github.com/xan105/node-cgo-regodit

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

func KeyExists(root string, path string) bool {

  HKEY := getRootKey(root)

  k, err := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  
  if err != nil {
    return false
  } else {
    return true
  }
}

func ListAllSubkeys(root string, path string) []string {
  
  HKEY := getRootKey(root)
  
  k, _ := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS)
  defer k.Close()
  
  list, _ := k.ReadSubKeyNames(-1)
  return list
}

func ListAllValues(root string, path string) []string {
  
  HKEY := getRootKey(root)
  
  k, _ := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS)
  defer k.Close()
  
  list, _ := k.ReadValueNames(-1)
  return list
}

func QueryValueType(root string, path string, key string) string {

  var buf []byte;
  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  _, valtype, _ := k.GetValue(key, buf)
 
  switch valtype {
    case 0: return "NONE"
    case 1: return "SZ"
    case 2: return "EXPAND_SZ"
    case 3: return "BINARY"
    case 4: return "DWORD"
    case 5: return "DWORD_BIG_ENDIAN"
    case 6: return "LINK"
    case 7: return "MULTI_SZ"
    case 8: return "RESOURCE_LIST"
    case 9: return "FULL_RESOURCE_DESCRIPTOR"
    case 10: return "RESOURCE_REQUIREMENTS_LIST"
    case 11: return "QWORD"
    default: return "NONE"
  }
}

func QueryStringValue(root string, path string, key string) string { //REG_SZ & REG_EXPAND_SZ

  var result string
  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY, filepath.FromSlash(path), registry.QUERY_VALUE)
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

func QueryMultiStringValue(root string, path string, key string) []string { //REG_MULTI_SZ

  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY, filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  list, _, _ := k.GetStringsValue(key)
  
  return list
}

func QueryBinaryValue(root string, path string, key string) []byte { //REG_BINARY

  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY, filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  bytes, _, _ := k.GetBinaryValue(key)

  return bytes
}

func QueryIntegerValue(root string, path string, key string) uint64 { //REG_DWORD & REG_QWORD

  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY, filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  i, _, _ := k.GetIntegerValue(key)
 
  return i
}

func WriteKey (root string, path string) {
  HKEY := getRootKey(root)
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
}

func DeleteKey (root string, path string) {
  HKEY := getRootKey(root)
  registry.DeleteKey(HKEY, filepath.FromSlash(path)) 
}

func WriteStringValue(root string, path string, key string, value string) { //REG_SZ
  
  HKEY := getRootKey(root)
  
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS)
  defer k.Close()
  k.SetStringValue(key, value)
}

func WriteExpandStringValue(root string, path string, key string, value string) { //REG_EXPAND_SZ
  
  HKEY := getRootKey(root)
  
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS)
  defer k.Close()
  k.SetExpandStringValue(key, value)
}

func WriteMultiStringValue(root string, path string, key string, value []string) { //REG_MULTI_SZ

  HKEY := getRootKey(root)

  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetStringsValue(key, value)
}

func WriteBinaryValue(root string, path string, key string, value []byte) { //REG_BINARY

  HKEY := getRootKey(root)
  
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetBinaryValue(key, value)
}

func WriteDwordValue(root string, path string, key string, value uint32) { //REG_DWORD

  HKEY := getRootKey(root)
  
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetDWordValue(key, value)
}

func WriteQwordValue(root string, path string, key string, value uint64) { //REG_QWORD

  HKEY := getRootKey(root)
  
  k, _, _ := registry.CreateKey(HKEY, filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetQWordValue(key, value)
}

func DeleteKeyValue (root string, path string, key string) {

  HKEY := getRootKey(root)

  k, _ := registry.OpenKey(HKEY , filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.DeleteValue(key)
}