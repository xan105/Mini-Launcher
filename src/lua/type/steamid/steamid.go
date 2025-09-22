/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package steamid

import (
  "launcher/internal/steam"
  "github.com/yuin/gopher-lua"
)

const Name = "SteamID"

func LValue(L *lua.LState, sid *steam.SteamID) lua.LValue {
  ud := L.NewUserData()
  ud.Value = sid
  L.SetMetatable(ud, L.GetTypeMetatable(Name))
  return ud
}

func LCheckSteamID(L *lua.LState, n int) *steam.SteamID {
  ud := L.CheckUserData(n)
  if v, ok := ud.Value.(*steam.SteamID); ok {
    return v
  }
  L.ArgError(n, "UserData<SteamID> type expected!")
  return nil
}

func Constructor(L *lua.LState) int {
  val := L.Get(1)
  id, ok := val.(lua.LString)
  if !ok {
    L.RaiseError("A SteamID requires an explicit String type to avoid 64bits precision loss")
  }
  sid, err := steam.ParseSteamID(string(id))
  if err != nil {
    L.RaiseError(err.Error())
    return 0
  }
  L.Push(LValue(L, sid))
  return 1
}

func index(L *lua.LState) int {
  sid := LCheckSteamID(L, 1)
  key := L.ToString(2)
  switch key {
    case "universe":
      L.Push(lua.LNumber(uint32(sid.Universe)))
    case "type":
      L.Push(lua.LNumber(uint32(sid.Type)))
    case "instance":
      L.Push(lua.LNumber(uint32(sid.Instance)))
    case "accountid":
      L.Push(lua.LNumber(uint32(sid.AccountID)))
    case "asSteam2":
      L.Push(L.NewFunction(func(L *lua.LState) int {
        sid := LCheckSteamID(L, 1)
        L.Push(lua.LString(sid.AsSteam2()))
        return 1
      }))
    case "asSteam3":
      L.Push(L.NewFunction(func(L *lua.LState) int {
        sid := LCheckSteamID(L, 1)
        L.Push(lua.LString(sid.AsSteam3()))
        return 1
      }))
    case "asSteam64":
      L.Push(L.NewFunction(func(L *lua.LState) int {
        sid := LCheckSteamID(L, 1)
        L.Push(lua.LString(sid.AsSteam64()))
        return 1
      }))
    default:
      L.Push(lua.LNil)
  }
  return 1
}

func tostring(L *lua.LState) int {
  sid := LCheckSteamID(L, 1)
  L.Push(lua.LString(sid.AsSteam2()))
  return 1
}

func RegisterType(L *lua.LState) {
  mt := L.NewTypeMetatable(Name)
  L.SetField(mt, "__index", L.NewFunction(index))
  L.SetField(mt, "__tostring", L.NewFunction(tostring))
  L.SetField(mt, "__metatable", lua.LString("Protected metatable!"))
  L.SetMetatable(mt, mt)
}