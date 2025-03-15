/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package elevated

import (
  "os"
  "strings"
  "syscall"
  "unsafe"
  "golang.org/x/sys/windows"
)

var (
  modShell32      = syscall.NewLazyDLL("shell32.dll")
  pShellExecuteW  = modShell32.NewProc("ShellExecuteW")
)

func IsElevated() bool {

  hProcess, err := windows.GetCurrentProcess()
  if err != nil {
    return false
  }

  var token windows.Token
  if err := windows.OpenProcessToken(hProcess, windows.TOKEN_QUERY, &token); err != nil {
    return false
  }
  defer token.Close()

  var elevation uint32
  var returned uint32
  if err := windows.GetTokenInformation(
    token,
    windows.TokenElevation,
    (*byte)(unsafe.Pointer(&elevation)),
    uint32(unsafe.Sizeof(elevation)),
    &returned,
  ); err != nil {
    return false
  }

  return elevation > 0
}

func RestartElevated() {

  exePath, _ := os.Executable()
  verb, _ := syscall.UTF16PtrFromString("runas")
  exe, _ := syscall.UTF16PtrFromString(exePath)
  args, _ := syscall.UTF16PtrFromString(strings.Join(os.Args[1:], " "))
      
  pShellExecuteW.Call(
    0,
    uintptr(unsafe.Pointer(verb)), 
    uintptr(unsafe.Pointer(exe)), 
    uintptr(unsafe.Pointer(args)), 
    0,
    10,
  )
  
  os.Exit(0)
}