/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package splash

import (
  "golang.org/x/sys/windows"
)

var (
  kernel32          = windows.NewLazySystemDLL("kernel32.dll")
  pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
)

func getModuleHandle() (windows.Handle, error) {
  ret, _, err := pGetModuleHandleW.Call(uintptr(0))
  if ret == 0 {
    return 0, err
  }
  return windows.Handle(ret), nil
}

