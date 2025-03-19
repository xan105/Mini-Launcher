/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package ui

import (
  "os"
  "unsafe"
  "golang.org/x/sys/windows"
)

var (
  shell32           = windows.NewLazySystemDLL("shell32.dll")
  pExtractIcon      = shell32.NewProc("ExtractIconW")
)

func extractIcon() windows.Handle {
  exePath, err := os.Executable()
  if err != nil {
    return 0
  }
  hIcon, _, _ := pExtractIcon.Call(
    uintptr(0), 
    uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(exePath))), 
    0,
  )
  return windows.Handle(hIcon)
}