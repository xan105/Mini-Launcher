/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package version

import (
  "syscall"
  "unsafe"
  "slices"
  "errors"
  "path/filepath"
)

var (
  version                 = syscall.NewLazyDLL("Version.dll")
  getFileVersionInfoSize  = version.NewProc("GetFileVersionInfoSizeW")
  getFileVersionInfo      = version.NewProc("GetFileVersionInfoW")
  verQueryValue           = version.NewProc("VerQueryValueW")
)

type VSFixedFileInfo struct {
  DwSignature        uint32
  DwStrucVersion     uint32
  DwFileVersionMS    uint32
  DwFileVersionLS    uint32
  DwProductVersionMS uint32
  DwProductVersionLS uint32
  DwFileFlagsMask    uint32
  DwFileFlags        uint32
  DwFileOS           uint32
  DwFileType         uint32
  DwFileSubtype      uint32
  DwFileDateMS       uint32
  DwFileDateLS       uint32
}

type FileVersion struct {
  Major    uint16
  Minor    uint16
  Build    uint16
  Revision uint16
}

func FromFile(filePath string) (FileVersion, error) {

  if !slices.Contains([]string{".exe", ".dll"}, filepath.Ext(filePath)) {
    return FileVersion{}, errors.New("Query version information only from binary (exe/dll)")
  }

  utf16Path, err := syscall.UTF16PtrFromString(filePath)
  if err != nil {
    return FileVersion{}, err
  }

  size, _, err := getFileVersionInfoSize.Call(uintptr(unsafe.Pointer(utf16Path)), 0)
  if size == 0 {
    return FileVersion{}, err
  }

  buf := make([]byte, size)
  if res, _, err := getFileVersionInfo.Call(
    uintptr(unsafe.Pointer(utf16Path)), 
    0, 
    size, 
    uintptr(unsafe.Pointer(&buf[0])),
    ); res == 0 {
    return FileVersion{}, err
  }

  var versionInfo unsafe.Pointer
  var versionInfoSize uint32

  if res, _, err := verQueryValue.Call(
    uintptr(unsafe.Pointer(&buf[0])),
    uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(`\`))),
    uintptr(unsafe.Pointer(&versionInfo)),
    uintptr(unsafe.Pointer(&versionInfoSize)),
  ); res == 0 {
    return FileVersion{}, err
  }

  info := (*VSFixedFileInfo)(versionInfo)
  version := FileVersion{
    Major:    uint16(info.DwFileVersionMS >> 16),
    Minor:    uint16(info.DwFileVersionMS & 0xFFFF),
    Build:    uint16(info.DwFileVersionLS >> 16),
    Revision: uint16(info.DwFileVersionLS & 0xFFFF),
  }

  return version, nil
}