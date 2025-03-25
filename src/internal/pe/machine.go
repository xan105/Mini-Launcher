/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package pe

import (
  "os"
  "errors"
  "encoding/binary"
)

func GetArchFromMachineType(path string) (string, error) {
  
  const location = 0x3C
  
  const (
    IMAGE_FILE_MACHINE_UNKNOWN  uint16 = 0x0000
    IMAGE_FILE_MACHINE_I386     uint16 = 0x014c
    IMAGE_FILE_MACHINE_AMD64    uint16 = 0x8664
    IMAGE_FILE_MACHINE_ARM64    uint16 = 0xAA64
    IMAGE_FILE_MACHINE_ARM      uint16 = 0x01C0 
  )
  
  var machineTypes = map[uint16]string{
    IMAGE_FILE_MACHINE_I386:  "386",
    IMAGE_FILE_MACHINE_AMD64: "amd64",
    IMAGE_FILE_MACHINE_ARM64: "arm64",
  }
  
  file, err := os.Open(path)
  if err != nil {
    return "", err
  }
  defer file.Close()

  header := make([]byte, 4)
  _, err = file.ReadAt(header, location)
  if err != nil {
    return "", err
  }
  offset := int64(binary.LittleEndian.Uint32(header))

  machine := make([]byte, 2)
  _, err = file.ReadAt(machine, offset + 4)
  if err != nil {
    return "", err
  }
  machineType := binary.LittleEndian.Uint16(machine)

  arch, found := machineTypes[machineType]
  if !found {
    return "", errors.New("Unsupported machine type !")
  }

  return arch, nil
}