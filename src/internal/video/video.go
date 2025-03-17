/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//Ported to Go From https://github.com/xan105/node-win-screen-resolution (MIT)

package video

import (
  "unsafe"
  "errors"
  "golang.org/x/sys/windows"
)

var (
  user32                                     = windows.NewLazySystemDLL("user32.dll")
  gdi32                                      = windows.NewLazySystemDLL("gdi32.dll")
  shcore                                     = windows.NewLazySystemDLL("shcore.dll")
  pSetThreadDpiAwarenessContext              = user32.NewProc("SetThreadDpiAwarenessContext")
  pGetDC                                     = user32.NewProc("GetDC")
  pReleaseDC                                 = user32.NewProc("ReleaseDC")
  pMonitorFromPoint                          = user32.NewProc("MonitorFromPoint")
  pGetDeviceCaps                             = gdi32.NewProc("GetDeviceCaps")
  pGetDpiForMonitor                          = shcore.NewProc("GetDpiForMonitor")
)

const (
  HORZRES                                    = 8
  VERTRES                                    = 10
  VREFRESH                                   = 116
  MDT_EFFECTIVE_DPI                          = 0
  MONITOR_DEFAULTTOPRIMARY                   = 0x00000001
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
  Scale   uint64
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
  
  var dpiX, dpiY uint32 = 96, 96 // Default DPI
  hMonitor, _, err := pMonitorFromPoint.Call(
    0, 
    0, 
    uintptr(MONITOR_DEFAULTTOPRIMARY),
  )
  if hMonitor == 0 {
    return mode, err
  }
  
  pGetDpiForMonitor.Call(
    hMonitor, 
    uintptr(MDT_EFFECTIVE_DPI), 
    uintptr(unsafe.Pointer(&dpiX)), 
    uintptr(unsafe.Pointer(&dpiY)),
  )
  
  //NB: Microsoft states that dpiX == dpiY and just to pick one
  mode.Scale = uint64((float64(dpiX) / 96.0) * 100.0)
  if dpiX != dpiY {
    return mode, errors.New("DPI X should equal DPI Y !!")
  }

  return mode, nil
}