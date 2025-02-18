/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package video

import (
  "syscall"
)

var (
  user32          = syscall.NewLazyDLL("user32.dll")
  gdi32           = syscall.NewLazyDLL("gdi32.dll")
  
  pSetThreadDpiAwarenessContext = user32.NewProc("SetThreadDpiAwarenessContext")
  pGetDC           = user32.NewProc("GetDC")
  pReleaseDC       = user32.NewProc("ReleaseDC")
  pGetDeviceCaps   = gdi32.NewProc("GetDeviceCaps")
)

const (
  HORZRES   = 8
  VERTRES   = 10
  VREFRESH  = 116
  // DPI awareness
  USER_DEFAULT_SCREEN_DPI                    = 96
  DPI_AWARENESS_CONTEXT_UNAWARE              = ^uintptr(1) + 1
  DPI_AWARENESS_CONTEXT_SYSTEM_AWARE         = ^uintptr(2) + 1
  DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE    = ^uintptr(3) + 1
  DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = ^uintptr(4) + 1
  DPI_AWARENESS_CONTEXT_UNAWARE_GDISCALED    = ^uintptr(5) + 1
)

type VideoMode struct {
  Width   uint64
  Height  uint64
  Hz      uint64
}

func GetCurrentDisplayMode() (VideoMode, error) {

  pSetThreadDpiAwarenessContext.Call(uintptr(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)) //>=Win10

  var mode VideoMode

  hdc, _, err := pGetDC.Call(0)
  if hdc == 0 {
    return mode, err
  }
  defer pReleaseDC.Call(0, hdc)

  width, _, _  := pGetDeviceCaps.Call(hdc, HORZRES)
  height, _, _ := pGetDeviceCaps.Call(hdc, VERTRES)
  hz, _, _     := pGetDeviceCaps.Call(hdc, VREFRESH)

  mode.Width = uint64(width)
  mode.Height = uint64(height)
  mode.Hz = uint64(hz)
  
  return mode, nil
}