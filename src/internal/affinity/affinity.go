/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package affinity

import (
  "golang.org/x/sys/windows"
)

var (
  kernel32                  = windows.NewLazySystemDLL("kernel32.dll")
  pSetProcessAffinityMask   = kernel32.NewProc("SetProcessAffinityMask")
)

func SetProcessAffinity(pid int, logicalCores []uint) error {

  hProcess, err := windows.OpenProcess(
    windows.PROCESS_SET_INFORMATION,
    false,
    uint32(pid),
  )
  if err != nil {
    return err
  }
  defer windows.CloseHandle(hProcess)
  
  var mask uintptr = 0;
  for _, core := range logicalCores {
    mask |= uintptr(1) << core
  }

  ret, _, err := pSetProcessAffinityMask.Call(
    uintptr(hProcess), 
    mask,
  )
  if ret == 0 {
    return err
  }
  return nil
}