/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//ref: https://github.com/xan105/node-cgo-regodit

package regedit

import (
  "strconv"
  "encoding/hex"
  "github.com/yuin/gopher-lua"
  "launcher/internal/regedit"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "KeyExists": KeyExists,
    "ListAllSubkeys": ListAllSubkeys,
    "ListAllValues": ListAllValues,
    "QueryValueType": QueryValueType,
    "QueryStringValue": QueryStringValue,
    "QueryMultiStringValue": QueryMultiStringValue,
    "QueryBinaryValue": QueryBinaryValue,
    "QueryIntegerValue": QueryIntegerValue,
    "WriteKey": WriteKey,
    "DeleteKey": DeleteKey,
    "WriteStringValue": WriteStringValue,
    "WriteMultiStringValue": WriteMultiStringValue,
    "WriteBinaryValue": WriteBinaryValue,
    "WriteDwordValue": WriteDwordValue,
    "WriteQwordValue": WriteQwordValue,
    "DeleteKeyValue": DeleteKeyValue,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func KeyExists(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)

  L.Push(lua.LBool(regedit.KeyExists(root, path)))
  return 1
}

func ListAllSubkeys(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  
  values := regedit.ListAllSubkeys(root, path)
  table := L.NewTable()
  for _, value := range values {
    table.Append(lua.LString(value))
  }
  
  L.Push(table)
  return 1
}

func ListAllValues(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  
  values := regedit.ListAllValues(root, path)
  table := L.NewTable()
  for _, value := range values {
    table.Append(lua.LString(value))
  }
  
  L.Push(table)
  return 1
}

func QueryValueType(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  key  := L.ToString(3)  

  value := regedit.QueryValueType(root, path, key)          
  L.Push(lua.LString(value))
  return 1
}

func QueryStringValue(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  key  := L.ToString(3)  

  value := regedit.QueryStringValue(root, path, key)          
  L.Push(lua.LString(value))
  return 1
}

func QueryMultiStringValue(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  key  := L.ToString(3)  

  values := regedit.QueryMultiStringValue(root, path, key)   
  table := L.NewTable()
  for _, value := range values {
    table.Append(lua.LString(value))
  }
  
  L.Push(table)
  return 1
}

func QueryBinaryValue(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  key  := L.ToString(3)  

  value := regedit.QueryBinaryValue(root, path, key)          
  L.Push(lua.LString(hex.EncodeToString(value)))
  return 1
}

func QueryIntegerValue(L *lua.LState) int {
  root := L.ToString(1)  
  path := L.ToString(2)
  key  := L.ToString(3)  

  value := regedit.QueryIntegerValue(root, path, key)
  L.Push(lua.LString(strconv.FormatUint(value, 10)))
  return 1
}

func WriteKey(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)

  regedit.WriteKey(root, path)   
  return 0
}

func DeleteKey(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)

  regedit.DeleteKey(root, path)   
  return 0
}

func WriteStringValue(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)
  key   := L.ToString(3)
  value := L.ToString(4)

  regedit.WriteStringValue(root, path, key, value)   
  return 0
}

func WriteMultiStringValue(L *lua.LState) int {
  root   := L.ToString(1)  
  path   := L.ToString(2)
  key    := L.ToString(3)
  table  := L.ToTable(4)
  
  var values []string
  table.ForEach(func(_, value lua.LValue) {
    if str, ok := value.(lua.LString); ok {
      values = append(values, string(str))
    }
  })

  regedit.WriteMultiStringValue(root, path, key, values)   
  return 0
}

func WriteBinaryValue(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)
  key   := L.ToString(3)
  value := L.ToString(4)

  x, _ := hex.DecodeString(value)
  regedit.WriteBinaryValue(root, path, key, x)   
  return 0
}

func WriteDwordValue(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)
  key   := L.ToString(3)
  value := L.ToString(4)

  i, _ := strconv.ParseUint(value, 10, 32)
  regedit.WriteDwordValue(root, path, key, uint32(i))   
  return 0
}

func WriteQwordValue(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)
  key   := L.ToString(3)
  value := L.ToString(4)

  i, _ := strconv.ParseUint(value, 10, 64)
  regedit.WriteQwordValue(root, path, key, i)   
  return 0
}

func DeleteKeyValue(L *lua.LState) int {
  root  := L.ToString(1)  
  path  := L.ToString(2)
  key   := L.ToString(3)

  regedit.DeleteKeyValue(root, path, key)   
  return 0
}