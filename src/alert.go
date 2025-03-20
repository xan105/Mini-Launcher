/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "log/slog"
  "golang.org/x/sys/windows"
  "launcher/lua"
)

func alert(title string, message string){
  slog.Warn(message)
  windows.MessageBox(
    windows.HWND(uintptr(0)),
    windows.StringToUTF16Ptr(message),
    windows.StringToUTF16Ptr(title),
    windows.MB_OK,
  )
}

func panic(title string, message string){
  slog.Error(message)
  windows.MessageBox(
    windows.HWND(uintptr(0)),
    windows.StringToUTF16Ptr(message),
    windows.StringToUTF16Ptr(title),
    windows.MB_OK | windows.MB_ICONERROR,
  )
  lua.CloseLua()
  os.Exit(1)
}