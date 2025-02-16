/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package random

import (
  "github.com/yuin/gopher-lua"
  "math/rand"
)

func randAlphaNumString(length int) string {
  //cf: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
  
  const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  const IdxBits = 6
  const IdxMask = 1<<IdxBits - 1
  const IdxMax = 63 / IdxBits

  bytes := make([]byte, length)
  for i, cache, remain := length-1, rand.Int63(), IdxMax; i >= 0; {
    if remain == 0 {
      cache, remain = rand.Int63(), IdxMax
    }
    if idx := int(cache & IdxMask); idx < len(charset) {
      bytes[i] = charset[idx]
      i--
    }
    cache >>= IdxBits
    remain--
  }
  return string(bytes)
}

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "AlphaNumString": AlphaNumString,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func AlphaNumString(L *lua.LState) int {
  //get argument
  length := L.ToInt(1)  

  //push result
  value:= randAlphaNumString(length)
  L.Push(lua.LString(value))
  return 1 //number of result(s)
}