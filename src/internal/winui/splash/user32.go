/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package splash

import (
  "unsafe"
  "golang.org/x/sys/windows"
)

var (
  user32                    = windows.NewLazySystemDLL("user32.dll")
  pCreateWindowExW          = user32.NewProc("CreateWindowExW")
  pDefWindowProcW           = user32.NewProc("DefWindowProcW")
  pDestroyWindow            = user32.NewProc("DestroyWindow")
  pDispatchMessageW         = user32.NewProc("DispatchMessageW")
  pGetMessageW              = user32.NewProc("GetMessageW")
  pPostQuitMessage          = user32.NewProc("PostQuitMessage")
  pRegisterClassExW         = user32.NewProc("RegisterClassExW")
  pTranslateMessage         = user32.NewProc("TranslateMessage")
  pLoadImageW               = user32.NewProc("LoadImageW")
  pSetWinEventHook          = user32.NewProc("SetWinEventHook")
  pUnhookWinEvent           = user32.NewProc("UnhookWinEvent")
  pGetDC                    = user32.NewProc("GetDC")
  pReleaseDC                = user32.NewProc("ReleaseDC")
)

const (
  WM_CREATE                 = 0x0001
  WM_DESTROY                = 0x0002
  WM_SHOWWINDOW             = 0x0018
  WS_VISIBLE                = 0x10000000
  WS_EX_TOPMOST             = 0x00000008
  WS_POPUP                  = 0x80000000
  WS_EX_TOOLWINDOW          = 0x000000080
  WS_TABSTOP                = 0x00010000
  EVENT_SYSTEM_FOREGROUND   = 0x0003
  EVENT_OBJECT_CREATE       = 0x8000
  EVENT_OBJECT_SHOW         = 0x8002
  WINEVENT_OUTOFCONTEXT     = 0x0000
  WINEVENT_INCONTEXT        = 0x0004
  WINEVENT_SKIPOWNPROCESS   = 0x0002
  WINEVENT_SKIPOWNTHREAD    = 0x0001
  OBJID_WINDOW              = 0
  OBJID_CURSOR              = -9
  OBJID_CLIENT              = -4;
  IMAGE_BITMAP              = 0x00
  LR_LOADFROMFILE           = 0x00000010
)

type POINT struct {
  x                         int32
  y                         int32
}

type MSG struct {
  hwnd                      windows.Handle
  message                   uint32
  wParam                    uintptr
  lParam                    uintptr
  time                      uint32
  pt                        POINT
}

type WNDCLASSEXW struct {
  size                      uint32
  style                     uint32
  wndProc                   uintptr
  clsExtra                  int32
  wndExtra                  int32
  instance                  windows.Handle
  icon                      windows.Handle
  cursor                    windows.Handle
  background                windows.Handle
  menuName                  *uint16
  className                 *uint16
  iconSm                    windows.Handle
}

func getDC(hWnd windows.Handle) (windows.Handle, error) {
  ret, _, err := pGetDC.Call(
    uintptr(hWnd),
  )
  if ret == 0 {
    return 0, err
  }
  return windows.Handle(ret), nil
}

func releaseDC (hWnd windows.Handle, hDC windows.Handle) bool {
  ret, _, _ := pReleaseDC.Call(
    uintptr(hWnd),
    uintptr(hDC),
  )
  return ret != 0
} 

func createWindow(className string, windowName string, style, style_ext uint32, x, y, width, height uint32, parent, menu, instance windows.Handle) (windows.Handle, error) {
  ret, _, err := pCreateWindowExW.Call(
    uintptr(style),
    uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(className))),
    uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(windowName))),
    uintptr(style_ext),
    uintptr(x),
    uintptr(y),
    uintptr(width),
    uintptr(height),
    uintptr(parent),
    uintptr(menu),
    uintptr(instance),
    uintptr(0),
  )
  if ret == 0 {
    return 0, err
  }
  return windows.Handle(ret), nil
}

func setWinEventHook(eventMin uint32, eventMax uint32, hmodWinEventProc windows.Handle, pfnWinEventProc uintptr, idProcess int, idThread, dwFlags uint32) (windows.Handle, error) {
  ret, _, err := pSetWinEventHook.Call(
    uintptr(eventMin),
    uintptr(eventMax),
    uintptr(hmodWinEventProc),
    pfnWinEventProc,
    uintptr(idProcess),
    uintptr(idThread),
    uintptr(dwFlags),
  )
  if ret == 0 {
    return 0, err
  }
  return windows.Handle(ret), nil
}

func unhookWinEvent(hWinEventHook windows.Handle) bool {
  ret, _, _ := pUnhookWinEvent.Call(
    uintptr(hWinEventHook),
  )
  return ret != 0
}

func loadImage(imagePath string) (windows.Handle, error) {
  ret, _, err := pLoadImageW.Call(
    uintptr(0),
    uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(imagePath))),
    uintptr(IMAGE_BITMAP),
    uintptr(0),
    uintptr(0),
    uintptr(LR_LOADFROMFILE),
  )
  if ret == 0 {
    return 0, err
  }
  return windows.Handle(ret), nil
}

func defWindowProc(hwnd windows.Handle, msg uint32, wparam, lparam uintptr) uintptr {
  ret, _, _ := pDefWindowProcW.Call(
    uintptr(hwnd),
    uintptr(msg),
    uintptr(wparam),
    uintptr(lparam),
  )
  return uintptr(ret)
}

func destroyWindow(hwnd windows.Handle) error {
  ret, _, err := pDestroyWindow.Call(uintptr(hwnd))
  if ret == 0 {
    return err
  }
  return nil
}

func registerClassEx(wcx *WNDCLASSEXW) (uint16, error) {
  ret, _, err := pRegisterClassExW.Call(
    uintptr(unsafe.Pointer(wcx)),
  )
  if ret == 0 {
    return 0, err
  }
  return uint16(ret), nil
}

func dispatchMessage(msg *MSG) {
  pDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}

func getMessage(msg *MSG, hwnd windows.Handle, msgFilterMin, msgFilterMax uint32) (bool, error) {
  ret, _, err := pGetMessageW.Call(
    uintptr(unsafe.Pointer(msg)),
    uintptr(hwnd),
    uintptr(msgFilterMin),
    uintptr(msgFilterMax),
  )
  if int32(ret) == -1 {
    return false, err
  }
  return int32(ret) != 0, nil
}

func postQuitMessage(exitCode int32) {
  pPostQuitMessage.Call(uintptr(exitCode))
}

func translateMessage(msg *MSG) {
  pTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}