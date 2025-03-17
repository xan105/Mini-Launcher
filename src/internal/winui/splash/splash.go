/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package splash

import (
  "log/slog"
  "os"
  "unsafe"
  "golang.org/x/sys/windows"
)

const classNameGUID = "D2FF2B71-7532-4BA6-8025-4D044372B710"  //Random GUID

func createBrushFromBMP(splashImage string) (windows.Handle, BITMAP, error) {
  hbm, err := loadImage(splashImage)
  if err != nil {    
    return 0, BITMAP{}, err
  } 
   
  hbrush, err := createPatternBrush(hbm)
  if err != nil {    
    return 0, BITMAP{}, err
  }
  
  //Get Image dimension
  image, err:= getObject(hbm)
  if err != nil {
    return hbrush, image, err
  }

  return hbrush, image, nil
}

func getScreenResolution() (uint32, uint32, error){
  hDC, err := getDC(0)
  if err != nil {
    return 0, 0, err
  }
  defer releaseDC(0, hDC)
  
  width := getDeviceCaps(hDC, HORZRES)
  height := getDeviceCaps(hDC, VERTRES)
  return width, height, nil
}

func createSplash(splashImage string, waitEvent string, pid int, exit chan bool) { 
  slog.Info("Create Splash Window")

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
        exit <- true
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
          postQuitMessage(0)
          exit <- true
        }
      }
      default:
        ret := defWindowProc(hwnd, msg, wparam, lparam)
        return ret
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
      dispatchMessage(&msg)
    } else {
      break
    }
  }
  exit <- true
}

func Show(splashImage string, waitEvent string, pid int) chan bool {
  exit := make(chan bool)
  go createSplash(splashImage, waitEvent, pid, exit)
  return exit
}