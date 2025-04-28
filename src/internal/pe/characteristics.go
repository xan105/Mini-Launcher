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

func PatchLargeAddress(filename string, aware bool) error {
  file, err := os.OpenFile(filename, os.O_RDWR, 0)
  if err != nil {
    return err
  }
  defer file.Close()

  var dosHeader ImageDosHeader
  err = binary.Read(file, binary.LittleEndian, &dosHeader)
  if err != nil {
    return err
  }
  if dosHeader.E_magic != IMAGE_DOS_SIGNATURE {
    return errors.New("Unexpected DOS signature")
  }

  //Seek to PE header
  _, err = file.Seek(int64(dosHeader.E_lfanew), 0)
  if err != nil {
    return err
  }

  //Read and check PE signature
  var peSignature uint32
  err = binary.Read(file, binary.LittleEndian, &peSignature)
  if err != nil {
    return err
  }
  if peSignature != IMAGE_NT_SIGNATURE {
    return errors.New("Unexpected PE signature")
  }

  //Read File Header
  var fileHeader ImageFileHeader
  err = binary.Read(file, binary.LittleEndian, &fileHeader)
  if err != nil {
    return err
  }
  if fileHeader.Machine != IMAGE_FILE_MACHINE_I386 {
    return nil //Not a x86 binary
  }

  if current := (fileHeader.Characteristics & IMAGE_FILE_LARGE_ADDRESS_AWARE) != 0; aware == current {
    return nil //Nothing to do
  }

  if aware {
    fileHeader.Characteristics |= IMAGE_FILE_LARGE_ADDRESS_AWARE
  } else {
    fileHeader.Characteristics &= ^IMAGE_FILE_LARGE_ADDRESS_AWARE
  }

  //Seek back to File Header position
  _, err = file.Seek(int64(dosHeader.E_lfanew)+4, 0) // +4 for PE signature
  if err != nil {
    return err
  }

  //Write the updated header
  err = binary.Write(file, binary.LittleEndian, &fileHeader)
  if err != nil {
    return err
  }

  return nil
}