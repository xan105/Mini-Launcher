/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//ref: https://github.com/xan105/node-cgo-regodit

package regedit

import (
  "strings"
  "path/filepath"
  "golang.org/x/sys/windows/registry"
)

var rootKeys = map[string]registry.Key{
	"HKCU": registry.CURRENT_USER,
	"HKLM": registry.LOCAL_MACHINE,
	"HKU":  registry.USERS,
	"HKCC": registry.CURRENT_CONFIG,
	"HKCR": registry.CLASSES_ROOT,
}

func KeyExists(root string, path string) bool {
  k, err := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  
  if err != nil {
    return false
  } else {
    return true
  }
}

func ListAllSubkeys(root string, path string) []string {
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS)
  defer k.Close()
  
  list, _ := k.ReadSubKeyNames(-1)
  return list
}

func ListAllValues(root string, path string) []string {
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS)
  defer k.Close()
  
  list, _ := k.ReadValueNames(-1)
  return list
}

func QueryValueType(root string, path string, key string) string {
  var buf []byte;
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE)
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
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE)
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
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  list, _, _ := k.GetStringsValue(key)
  
  return list
}

func QueryBinaryValue(root string, path string, key string) []byte { //REG_BINARY
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  bytes, _, _ := k.GetBinaryValue(key)

  return bytes
}

func QueryIntegerValue(root string, path string, key string) uint64 { //REG_DWORD & REG_QWORD
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.QUERY_VALUE)
  defer k.Close()
  i, _, _ := k.GetIntegerValue(key)
 
  return i
}

func Create(root string, path string) {
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
}

func Delete(root string, path string) {
  subkeys := ListAllSubkeys(root, path)
  for _, subkey := range subkeys {
    Delete(root, subkey)
  }
  registry.DeleteKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path))
}

func WriteStringValue(root string, path string, key string, value string) { //REG_SZ
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS)
  defer k.Close()
  k.SetStringValue(key, value)
}

func WriteExpandStringValue(root string, path string, key string, value string) { //REG_EXPAND_SZ
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS)
  defer k.Close()
  k.SetExpandStringValue(key, value)
}

func WriteMultiStringValue(root string, path string, key string, value []string) { //REG_MULTI_SZ
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetStringsValue(key, value)
}

func WriteBinaryValue(root string, path string, key string, value []byte) { //REG_BINARY
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetBinaryValue(key, value)
}

func WriteDwordValue(root string, path string, key string, value uint32) { //REG_DWORD
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetDWordValue(key, value)
}

func WriteQwordValue(root string, path string, key string, value uint64) { //REG_QWORD
  k, _, _ := registry.CreateKey(rootKeys[strings.ToUpper(root)], filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.SetQWordValue(key, value)
}

func DeleteValue (root string, path string, key string) {
  k, _ := registry.OpenKey(rootKeys[strings.ToUpper(root)] , filepath.FromSlash(path), registry.ALL_ACCESS) 
  defer k.Close()
  k.DeleteValue(key)
}