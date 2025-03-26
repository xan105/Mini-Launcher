/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package version

import (
  "unsafe"
  "golang.org/x/sys/windows"
)

type FileVersion struct {
  Major                   uint16
  Minor                   uint16
  Build                   uint16
  Revision                uint16
}

func FromFile(filePath string) (FileVersion, error) {

  size, err := windows.GetFileVersionInfoSize(filePath, nil)
  if err != nil {
    return FileVersion{}, err
  }

  buf := make([]byte, size)
  if err := windows.GetFileVersionInfo(
    filePath, 
    0,
    size,
    unsafe.Pointer(&buf[0]),
  ); err != nil {
    return FileVersion{}, err
  }

  var versionInfo *windows.VS_FIXEDFILEINFO
  var versionInfoSize uint32
  if err := windows.VerQueryValue(unsafe.Pointer(
    &buf[0]), 
    "\\", 
    unsafe.Pointer(&versionInfo), 
    &versionInfoSize,
  ); err != nil {
    return FileVersion{}, err
  }

  version := FileVersion{
    Major:    uint16(versionInfo.FileVersionMS >> 16),
    Minor:    uint16(versionInfo.FileVersionMS & 0xFFFF),
    Build:    uint16(versionInfo.FileVersionLS >> 16),
    Revision: uint16(versionInfo.FileVersionLS & 0xFFFF),
  }

  return version, nil
}