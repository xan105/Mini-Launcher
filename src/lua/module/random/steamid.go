/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package random

import (
  "time"
  "math/rand"
  "strconv"
  "github.com/yuin/gopher-lua"
  "launcher/internal/steam"
)

func SteamID(L *lua.LState) int {
  rand.Seed(time.Now().UnixNano())
  accountid := strconv.FormatUint(uint64(rand.Uint32()), 10)
  sid, _ := steam.ParseSteamID("[U:1:" + accountid +"]")
  L.Push(lua.LString(sid.AsSteam64()))
  return 1
}


