/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package failure

import (
  "github.com/yuin/gopher-lua"
)

const Name = "Failure"

type Failure struct {
  Code    string
  Message string
}

func LValue(L *lua.LState, code string, message string) lua.LValue {
  ud := L.NewUserData()
  ud.Value = &Failure{Code: code, Message: message}
  L.SetMetatable(ud, L.GetTypeMetatable(Name))
  return ud
}

func LCheckFailure(L *lua.LState, n int) *Failure {
  ud := L.CheckUserData(n)
  if err, ok := ud.Value.(*Failure); ok {
    return err
  }
  L.ArgError(n, "UserData<Failure> type expected!")
  return nil
}

func constructor(L *lua.LState) int {
  code := L.ToString(1)
  if len(code) == 0 { code = "ERR_UNKNOWN" }
  
  message := L.ToString(2)
  if len(message) == 0 { message = "An unknown error occurred" }

  L.Push(LValue(L, code, message))
  return 1
}

func index(L *lua.LState) int {
  err := LCheckFailure(L, 1)
  key := L.ToString(2)
  switch key {
    case "code":
      L.Push(lua.LString(err.Code))
    case "message":
      L.Push(lua.LString(err.Message))
    default:
      L.Push(lua.LNil)
  }
  return 1
}

func tostring(L *lua.LState) int {
  err := LCheckFailure(L, 1)
  L.Push(lua.LString("[" + err.Code + "]: " + err.Message))
  return 1
}

func RegisterType(L *lua.LState) {
  mt := L.NewTypeMetatable(Name)
  L.SetField(mt, "__call", L.NewFunction(constructor))
  L.SetField(mt, "__index", L.NewFunction(index))
  L.SetField(mt, "__tostring", L.NewFunction(tostring))
  L.SetField(mt, "__metatable", lua.LString("Protected metatable!"))
  L.SetGlobal(Name, mt)
  L.SetMetatable(mt, mt) //https://github.com/yuin/gopher-lua/issues/36#issuecomment-113885402
}