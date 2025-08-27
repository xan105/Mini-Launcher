/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package ui

import (
  "log/slog"
  "os"
  "unsafe"
  "runtime"
  "golang.org/x/sys/windows"
)

func createSplashWindow(splashImage string, waitEvent string, pid int, exit chan bool) { 
  runtime.LockOSThread() //GetMessageW() must be called in the same thread
  defer runtime.UnlockOSThread()
  
  slog.Info("Create Splash Window")
  
  const classNameGUID = "D2FF2B71-7532-4BA6-8025-4D044372B710"  //Random GUID

  var win windows.Handle

  activeWinEventHook := func(hWinEventHook windows.Handle, event uint32, hwnd windows.HWND, idObject int32, idChild int32, idEventThread uint32, dwmsEventTime uint32) uintptr {
    if (waitEvent == "FOREGROUND" && (event == EVENT_SYSTEM_FOREGROUND && windows.IsWindowVisible(hwnd))) ||
       (waitEvent == "WINDOW" && (event == EVENT_OBJECT_SHOW && idObject == OBJID_WINDOW)) ||
       (waitEvent == "CURSOR" && (event == EVENT_OBJECT_SHOW && idObject == OBJID_CURSOR)) {
        slog.Info("Splash bye bye")
        destroyWindow(win)
        unhookWinEvent(hWinEventHook) 
    }
    return 0
  }

  lpfnWndProc := func(hwnd windows.Handle, msg uint32, wparam uintptr, lparam uintptr) uintptr {
    switch msg {
      case WM_DESTROY:
        postQuitMessage(0)
      case WM_SHOWWINDOW: {
        _, err := setWinEventHook(
          EVENT_SYSTEM_FOREGROUND,
          EVENT_OBJECT_SHOW,
          0, 
          windows.NewCallback(activeWinEventHook), 
          pid,
          0, 
          WINEVENT_OUTOFCONTEXT | WINEVENT_SKIPOWNPROCESS,
        )
        if err != nil {
          slog.Error(err.Error())
          destroyWindow(hwnd)
        }
      }
      default:
        return defWindowProc(hwnd, msg, wparam, lparam)
    }
    return 0
  }
  
  instance, err := getModuleHandle()
  if err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }
  
  hbrush, image, err := createBrushFromBMP(splashImage)
  if err != nil {    
    slog.Error(err.Error())
    exit <- true
    return
  }
  
  wcx := WNDCLASSEXW{
    wndProc:    windows.NewCallback(lpfnWndProc),
    instance:   instance,
    background: hbrush,
    className:  windows.StringToUTF16Ptr(classNameGUID),
  }
  wcx.size = uint32(unsafe.Sizeof(wcx))
  
  if _, err := registerClassEx(&wcx); err != nil {
    slog.Error(err.Error())
    exit <- true
    return
  }
  
  screenWidth, screenHeight, err := getScreenResolution()
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
  
  win, err = createWindow(
    classNameGUID,
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
    msg := MSG{}
    gotMessage, err := getMessage(&msg, 0, 0, 0)
    if err != nil {
      slog.Error(err.Error())
      exit <- true
      return
    }

    if gotMessage {
      translateMessage(&msg)
      if msg.message == WM_QUIT {
        break
      }
      dispatchMessage(&msg)
    } else {
      break
    }
  }
  exit <- true
  return
}

func Splash(splashImage string, waitEvent string, pid int) chan bool {
  exit := make(chan bool)
  go createSplashWindow(splashImage, waitEvent, pid, exit)
  return exit
}