/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package locale

import (
  "unsafe"
  "strings"
  "golang.org/x/sys/windows"
)

var (
  kernel32                    = windows.NewLazySystemDLL("kernel32.dll")
  pGetUserDefaultLocaleName   = kernel32.NewProc("GetUserDefaultLocaleName")
  pGetLocaleInfoEx            = kernel32.NewProc("GetLocaleInfoEx")
)

const (
  LOCALE_NAME_MAX_LENGTH      = 85
  LOCALE_SENGLISHLANGUAGENAME = 0x1001
)

func GetUserLocale() (string, error) {
  buffer := make([]uint16, LOCALE_NAME_MAX_LENGTH)
  ret, _, err := pGetUserDefaultLocaleName.Call(
    uintptr(unsafe.Pointer(&buffer[0])),
    uintptr(LOCALE_NAME_MAX_LENGTH),
  )
  if ret == 0 {
    return "", err
  }
  return windows.UTF16ToString(buffer), nil
}

func GetLanguageFromLocale(locale string) (string, error) {
  localePtr, _ := windows.UTF16PtrFromString(locale)
  buffer := make([]uint16, 100)

  ret, _, err := pGetLocaleInfoEx.Call(
    uintptr(unsafe.Pointer(localePtr)),
    uintptr(LOCALE_SENGLISHLANGUAGENAME),
    uintptr(unsafe.Pointer(&buffer[0])),
    uintptr(len(buffer)),
  )
  if ret == 0 {
    return "", err
  }
  return strings.ToLower(windows.UTF16ToString(buffer)), nil
}
