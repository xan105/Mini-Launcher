/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "golang.org/x/sys/windows"
)

func alert(message string){
  windows.MessageBox(
    windows.HWND(uintptr(0)),
    windows.StringToUTF16Ptr(message),
    windows.StringToUTF16Ptr("Launcher"),
    windows.MB_OK,
  )
}

func panic(message string){
  windows.MessageBox(
    windows.HWND(uintptr(0)),
    windows.StringToUTF16Ptr(message),
    windows.StringToUTF16Ptr("Launcher"),
    windows.MB_OK | windows.MB_ICONERROR,
  )
  os.Exit(1)
}