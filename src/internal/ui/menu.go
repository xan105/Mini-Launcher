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

func createMenuWindow(labels []string, button chan int) { 
  runtime.LockOSThread() //GetMessageW() must be called in the same thread
  slog.Info("Create Menu Window")
  
  const (
    classNameGUID   = "3949D435-2861-42B9-9E47-C721A71F85E9"  //Random GUID
    buttonClass     = "BUTTON"  //Windows predefined class
    titleBar        = 32
    charWidthMin    = 10
    borderOffset    = 20
    buttonHeight    = 50
    separator       = 20
    paddingY        = 60
    paddingX        = 100
  )

  lpfnWndProc := func(hwnd windows.Handle, msg uint32, wparam uintptr, lparam uintptr) uintptr {
    switch msg {
      case WM_DESTROY: 
        postQuitMessage(0)
      case WM_CLOSE:
        destroyWindow(hwnd)
        os.Exit(0)
      case WM_COMMAND:
        id := uint16(wparam)
        index := int(id)
        button <- index
        destroyWindow(hwnd)
      default:
        return defWindowProc(hwnd, msg, wparam, lparam)
    }
    return 0
  }
  
  instance, err := getModuleHandle()
  if err != nil {
    slog.Error(err.Error())
    button <- -1
    return
  }
  
  hIcon := extractIcon()
  wcx := WNDCLASSEXW{
    wndProc:    windows.NewCallback(lpfnWndProc),
    instance:   instance,
    background: COLOR_WINDOW,
    className:  windows.StringToUTF16Ptr(classNameGUID),
    icon:      windows.Handle(hIcon),
    iconSm:    windows.Handle(hIcon),
  }
  wcx.size = uint32(unsafe.Sizeof(wcx))
  
  if _, err := registerClassEx(&wcx); err != nil {
    slog.Error(err.Error())
    button <- -1
    return
  }
  
  screenWidth, screenHeight, err := getScreenResolution()
  if err != nil {
    slog.Error(err.Error())
    button <- -1
    return
  }

  charLen := 0
  for _, label := range labels {
    if len(label) > charLen {
      charLen = len(label)
    }
  }
  buttonWidth     := (charLen * charWidthMin) + borderOffset
  menuWidth       := buttonWidth + paddingX
  menuHeight      := titleBar + ((buttonHeight + separator) * len(labels)) + paddingY 
  buttonPosX      := paddingX / 2
  buttonPosY      := paddingY / 2 
  
  win, err := createWindow(
    classNameGUID,
    "Launcher",
    0,
    (WS_SYSMENU | WS_VISIBLE | WS_TABSTOP) &^ WS_MAXIMIZEBOX,
    (screenWidth - uint32(menuWidth)) / 2, //center X
    (screenHeight - uint32(menuHeight)) / 2, //center Y
    uint32(menuWidth),
    uint32(menuHeight),
    0,
    0,
    instance,
  )
  if err != nil {
    slog.Error(err.Error())
    button <- -1
    return
  }

  for i, label := range labels{
    if _, err := createWindow( 
      buttonClass,
      label,
      0,
      WS_TABSTOP | WS_VISIBLE | WS_CHILD | BS_DEFPUSHBUTTON,
      uint32(buttonPosX),
      uint32(buttonPosY),
      uint32(buttonWidth),
      uint32(buttonHeight),
      win,
      windows.Handle(uint16(i)),
      instance,
    ); err != nil {
      slog.Error(err.Error())
      button <- -1
      return
    }
    buttonPosY += buttonHeight + separator
  }

  for {
    msg := MSG{}
    gotMessage, err := getMessage(&msg, 0, 0, 0)
    if err != nil {
      slog.Error(err.Error())
      button <- -1
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
  button <- -1
  return
}

func Menu(labels []string) chan int {
  button := make(chan int)
  go createMenuWindow(labels, button)
  return button
}