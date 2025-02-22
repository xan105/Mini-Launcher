/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package array

import (
  "github.com/yuin/gopher-lua"
)

func Find(L *lua.LState) int {
  table := L.CheckTable(1)
  fn := L.CheckFunction(2)

  for i := 1; i <= table.Len(); i++ {
    value := table.RawGetInt(i)

    L.Push(fn)
    L.Push(value)
    err := L.PCall(1, 1, nil)
    if err != nil {
      L.RaiseError("Error calling function: %v", err)
    }

    if L.ToBool(-1) {
      L.Push(value)
      return 1
    }

    L.Pop(1)
  }

  L.Push(lua.LNil)
  return 1
}

func Some(L *lua.LState) int {
  table := L.CheckTable(1)
  fn := L.CheckFunction(2)

  for i := 1; i <= table.Len(); i++ {
    value := table.RawGetInt(i)

    L.Push(fn)
    L.Push(value)
    err := L.PCall(1, 1, nil)
    if err != nil {
      L.RaiseError("Error calling function: %v", err)
    }

    if L.ToBool(-1) {
      L.Push(lua.LTrue)
      return 1
    }

    L.Pop(1)
  }

  L.Push(lua.LFalse)
  return 1
}

func Includes(L *lua.LState) int {
  table := L.CheckTable(1)
  searchElement := L.CheckAny(2)

  for i := 1; i <= table.Len(); i++ {
    value := table.RawGetInt(i)

    if value == searchElement {
      L.Push(lua.LTrue)
      return 1
    }
  }

  L.Push(lua.LFalse)
  return 1
}