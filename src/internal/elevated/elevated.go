/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package elevated

import (
  "os"
  "strings"
  "unsafe"
  "golang.org/x/sys/windows"
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
  verb, _ := windows.UTF16PtrFromString("runas")
  exe, _ := windows.UTF16PtrFromString(exePath)
  args, _ := windows.UTF16PtrFromString(strings.Join(os.Args[1:], " "))
      
  windows.ShellExecute(
    0,
    verb, 
    exe, 
    args, 
    nil,
    windows.SW_SHOWDEFAULT,
  )
  
  os.Exit(0)
}