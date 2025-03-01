/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package splash

import (
  "os"
  "unsafe"
  "syscall"
  "log/slog"
)

var (
  kernel32 = syscall.NewLazyDLL("kernel32.dll")

  pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
)

func getModuleHandle() (syscall.Handle, error) {
  ret, _, err := pGetModuleHandleW.Call(uintptr(0))
  if ret == 0 {
    return 0, err
  }
  return syscall.Handle(ret), nil
}

var (
  gdi32 = syscall.NewLazyDLL("Gdi32.dll")
  
  pGetDeviceCaps      = gdi32.NewProc("GetDeviceCaps")
  pCreatePatternBrush = gdi32.NewProc("CreatePatternBrush")
  pGetObjectW         = gdi32.NewProc("GetObjectW")
)

const (
  HORZRES = 8
  VERTRES = 10
)

func getDeviceCaps(hDC syscall.Handle, index int32) uint32 {
  ret, _, _ := pGetDeviceCaps.Call(
    uintptr(hDC),
    uintptr(index),
  )

  return uint32(ret)
}

func createPatternBrush(hbm syscall.Handle) (syscall.Handle, error) {
  ret, _, err := pCreatePatternBrush.Call(uintptr(hbm))
  if ret == 0 {
    return 0, err
  }
  return syscall.Handle(ret), nil
}

type BITMAP struct {
  bmType        uint32
  bmWidth       int32
  bmHeight      int32
  bmWidthBytes  uint32
  bmPlanes      uint16
  bmBitsPixel   uint16
  bmBits        uintptr
}

func getObject(hBitmap syscall.Handle) (BITMAP, error) {
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

var (
  user32 = syscall.NewLazyDLL("user32.dll")

  pCreateWindowExW          = user32.NewProc("CreateWindowExW")
  pDefWindowProcW           = user32.NewProc("DefWindowProcW")
  pDestroyWindow            = user32.NewProc("DestroyWindow")
  pDispatchMessageW         = user32.NewProc("DispatchMessageW")
  pGetMessageW              = user32.NewProc("GetMessageW")
  pLoadCursorW              = user32.NewProc("LoadCursorW")
  pPostQuitMessage          = user32.NewProc("PostQuitMessage")
  pRegisterClassExW         = user32.NewProc("RegisterClassExW")
  pTranslateMessage         = user32.NewProc("TranslateMessage")
  pLoadImageW               = user32.NewProc("LoadImageW")
  pSetWinEventHook          = user32.NewProc("SetWinEventHook")
  pUnhookWinEvent           = user32.NewProc("UnhookWinEvent")
  pGetDC                    = user32.NewProc("GetDC")
  pReleaseDC                = user32.NewProc("ReleaseDC")
  pIsWindowVisible          = user32.NewProc("IsWindowVisible")
)

func isWindowVisible(hWnd syscall.Handle) bool {
  ret, _, _ := pIsWindowVisible.Call(uintptr(hWnd))
  return ret != 0
}

func getDC(hWnd syscall.Handle) (syscall.Handle, error) {
  ret, _, err := pGetDC.Call(
    uintptr(hWnd),
  )
  if ret == 0 {
    return 0, err
  }
  return syscall.Handle(ret), nil
}

func releaseDC (hWnd syscall.Handle, hDC syscall.Handle) bool {
  ret, _, _ := pReleaseDC.Call(
    uintptr(hWnd),
    uintptr(hDC),
  )
  return ret != 0
} 

const (
  WS_VISIBLE       = 0x10000000
  WS_EX_TOPMOST    = 0x00000008
  WS_POPUP         = 0x80000000
  WS_EX_TOOLWINDOW = 0x000000080
  WS_TABSTOP       = 0x00010000
)

func createWindow(className, windowName string, style, style_ext uint32, x, y, width, height uint32, parent, menu, instance syscall.Handle) (syscall.Handle, error) {
  ret, _, err := pCreateWindowExW.Call(
    uintptr(style),
    uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className))),
    uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))),
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
  return syscall.Handle(ret), nil
}

const (
    EVENT_SYSTEM_FOREGROUND     = 0x0003
    EVENT_OBJECT_CREATE         = 0x8000
    EVENT_OBJECT_SHOW           = 0x8002
    WINEVENT_OUTOFCONTEXT       = 0x0000
    WINEVENT_INCONTEXT          = 0x0004
    WINEVENT_SKIPOWNPROCESS     = 0x0002
    WINEVENT_SKIPOWNTHREAD      = 0x0001
    OBJID_WINDOW                = 0
    OBJID_CURSOR                = -9
    OBJID_CLIENT                = -4;
)

func setWinEventHook(eventMin, eventMax uint32, hmodWinEventProc syscall.Handle, pfnWinEventProc uintptr, idProcess int, idThread, dwFlags uint32) (syscall.Handle, error) {
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
  return syscall.Handle(ret), nil
}

func unhookWinEvent(hWinEventHook syscall.Handle) bool {
  ret, _, _ := pUnhookWinEvent.Call(
    uintptr(hWinEventHook),
  )
  return ret != 0
}

const (
  IMAGE_BITMAP     = 0x00
  LR_LOADFROMFILE  = 0x00000010
)

func loadImage(imagePath string) (syscall.Handle, error) {
  ret, _, err := pLoadImageW.Call(
    uintptr(0),
    uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(imagePath))),
    uintptr(IMAGE_BITMAP),
    uintptr(0),
    uintptr(0),
    uintptr(LR_LOADFROMFILE),
  )
  if ret == 0 {
    return 0, err
  }
  return syscall.Handle(ret), nil
}

const (
  WM_CREATE     = 0x0001
  WM_DESTROY    = 0x0002
  WM_SHOWWINDOW = 0x0018
)

func defWindowProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
  ret, _, _ := pDefWindowProcW.Call(
    uintptr(hwnd),
    uintptr(msg),
    uintptr(wparam),
    uintptr(lparam),
  )
  return uintptr(ret)
}

func destroyWindow(hwnd syscall.Handle) error {
  ret, _, err := pDestroyWindow.Call(uintptr(hwnd))
  if ret == 0 {
    return err
  }
  return nil
}

type tPOINT struct {
  x, y int32
}

type tMSG struct {
  hwnd    syscall.Handle
  message uint32
  wParam  uintptr
  lParam  uintptr
  time    uint32
  pt      tPOINT
}

func dispatchMessage(msg *tMSG) {
  pDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}

func getMessage(msg *tMSG, hwnd syscall.Handle, msgFilterMin, msgFilterMax uint32) (bool, error) {
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

type WNDCLASSEXW struct {
  size       uint32
  style      uint32
  wndProc    uintptr
  clsExtra   int32
  wndExtra   int32
  instance   syscall.Handle
  icon       syscall.Handle
  cursor     syscall.Handle
  background syscall.Handle
  menuName   *uint16
  className  *uint16
  iconSm     syscall.Handle
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

func translateMessage(msg *tMSG) {
  pTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}

func CreateWindow(exit chan bool, pid int, splashImage string, waitEvent string) {
  
  slog.Info("splash")

  className := "D2FF2B71-7532-4BA6-8025-4D044372B710" //Random GUID
  instance, err := getModuleHandle()
  if err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }

  var win syscall.Handle

  activeWinEventHook := func(hWinEventHook syscall.Handle, event uint32, hwnd syscall.Handle, idObject int32, idChild int32, idEventThread uint32, dwmsEventTime uint32) uintptr {

    if (waitEvent == "FOREGROUND" && (event == EVENT_SYSTEM_FOREGROUND && isWindowVisible(hwnd))) ||
       (waitEvent == "WINDOW" && (event == EVENT_OBJECT_SHOW && idObject == OBJID_WINDOW)) ||
       (waitEvent == "CURSOR" && (event == EVENT_OBJECT_SHOW && idObject == OBJID_CURSOR)) {
      slog.Info("Splash bye bye")
      destroyWindow(win)
      unhookWinEvent(hWinEventHook) 
    }

    return 0
  }

  lpfnWndProc := func(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
    switch msg {
    case WM_DESTROY:
      postQuitMessage(0)
      exit <- true
    case WM_SHOWWINDOW: {
      _, err := setWinEventHook(
        EVENT_SYSTEM_FOREGROUND,
        EVENT_OBJECT_SHOW,
        0, 
        syscall.NewCallback(activeWinEventHook), 
        pid,
        0, 
        WINEVENT_OUTOFCONTEXT | WINEVENT_SKIPOWNPROCESS,
      )
      if err != nil {
        slog.Error(err.Error())
        exit <- true
      }
      
    }
    default:
      ret := defWindowProc(hwnd, msg, wparam, lparam)
      return ret
    }
    return 0
  }
  
  
  //load bitmap
  hbm , err := loadImage(splashImage)
    if err != nil {    
      slog.Error(err.Error())
      exit <- true
      return
    } 
  //and create an HBRUSH around it using CreatePatternBrush(), and then assign that to the WNDCLASS::hbrBackground    
  hbrush, err := createPatternBrush(hbm)
    if err != nil {    
      slog.Error(err.Error())
      exit <- true
      return
    }
  
  wcx := WNDCLASSEXW{
    wndProc:    syscall.NewCallback(lpfnWndProc),
    instance:   instance,
    background: hbrush,
    className:  syscall.StringToUTF16Ptr(className),
  }
  wcx.size = uint32(unsafe.Sizeof(wcx))

  if _, err = registerClassEx(&wcx); err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }

  //Get screen resolution
  hDC, err := getDC(0)
  if err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }
  defer releaseDC(0, hDC)
  screenWidth := getDeviceCaps(hDC, HORZRES)
  screenHeight := getDeviceCaps(hDC, VERTRES)

  //Get Image dimension
  image, err:= getObject(hbm)
  if err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }

  //check process hasn't crashed since we started it
  if _, err = os.FindProcess(pid); err != nil { 
    slog.Error(err.Error())
    exit <- true
    return
  }

  //create window
  win, err = createWindow(
    className,
    "Launcher",
    WS_EX_TOOLWINDOW | WS_EX_TOPMOST,
    WS_VISIBLE | WS_POPUP | WS_TABSTOP,
    (screenWidth - uint32(image.bmWidth)) / 2, //center X
    (screenHeight - uint32(image.bmHeight)) / 2, //center Y
    uint32(image.bmWidth),
    uint32(image.bmHeight),
    0,
    0,
    instance,
  )
  if err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }

  for {
    msg := tMSG{}
    gotMessage, err := getMessage(&msg, 0, 0, 0)
    if err != nil {
      slog.Error(err.Error())
      exit <- true
      return
    }

    if gotMessage {
      translateMessage(&msg)
      dispatchMessage(&msg)
    } else {
      break
    }
  }
  exit <- true
}