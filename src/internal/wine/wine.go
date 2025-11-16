/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package wine

import (
  "golang.org/x/sys/windows"
)

var (
  checked = false
  wine    = false
)

func IsWineOrProton() bool { //Check for wine_get_version() to detect Wine/Proton

  if checked {
    return wine
  }

  ntdll, err := windows.LoadLibrary("ntdll.dll")
  if err != nil {
    checked = true
    return wine
  }
  defer windows.FreeLibrary(ntdll)

  procAddr, err := windows.GetProcAddress(ntdll, "wine_get_version")
  if err == nil && procAddr != 0 {
    wine = true
  }
  
  checked = true
  return wine
}