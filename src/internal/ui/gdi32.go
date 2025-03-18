/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package ui

import (
  "unsafe"
  "golang.org/x/sys/windows"
)

var (
  gdi32               = windows.NewLazySystemDLL("Gdi32.dll")
  pGetDeviceCaps      = gdi32.NewProc("GetDeviceCaps")
  pCreatePatternBrush = gdi32.NewProc("CreatePatternBrush")
  pGetObjectW         = gdi32.NewProc("GetObjectW")
)

const (
  HORZRES             = 8
  VERTRES             = 10
)

type BITMAP struct {
  bmType              uint32
  bmWidth             int32
  bmHeight            int32
  bmWidthBytes        uint32
  bmPlanes            uint16
  bmBitsPixel         uint16
  bmBits              uintptr
}

func getDeviceCaps(hDC windows.Handle, index int32) uint32 {
  ret, _, _ := pGetDeviceCaps.Call(
    uintptr(hDC),
    uintptr(index),
  )

  return uint32(ret)
}

func createPatternBrush(hbm windows.Handle) (windows.Handle, error) {
  ret, _, err := pCreatePatternBrush.Call(uintptr(hbm))
  if ret == 0 {
    return 0, err
  }
  return windows.Handle(ret), nil
}

func getObject(hBitmap windows.Handle) (BITMAP, error) {
  var bmp BITMAP
  
  ret, _, err := pGetObjectW.Call(
    uintptr(hBitmap), 
    uintptr(unsafe.Sizeof(bmp)), 
    uintptr(unsafe.Pointer(&bmp)),
  )
  if ret == 0 {
    return bmp, err
  }
  return bmp, nil
}