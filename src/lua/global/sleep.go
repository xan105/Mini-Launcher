/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package global

import (
  "time"
  "github.com/yuin/gopher-lua"
)

func Sleep(L *lua.LState) int {
  interval := L.CheckInt(1)
  
  wait := make(chan struct{})
  go func() {
      time.Sleep(time.Millisecond * time.Duration(interval))
      close(wait)
  }()
  <-wait

  return 0
}