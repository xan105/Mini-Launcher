/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package pe

const (
  IMAGE_FILE_MACHINE_I386         uint16 = 0x014c
  IMAGE_FILE_MACHINE_AMD64        uint16 = 0x8664
  IMAGE_FILE_MACHINE_ARM64        uint16 = 0xAA64
  IMAGE_DOS_SIGNATURE             uint16 = 0x5A4D
  IMAGE_FILE_LARGE_ADDRESS_AWARE  uint16 = 0x0020
  IMAGE_NT_SIGNATURE              uint32 = 0x00004550
)

type ImageDosHeader struct {
  E_magic    uint16
  E_cblp     uint16
  E_cp       uint16
  E_crlc     uint16
  E_cparhdr  uint16
  E_minalloc uint16
  E_maxalloc uint16
  E_ss       uint16
  E_sp       uint16
  E_csum     uint16
  E_ip       uint16
  E_cs       uint16
  E_lfarlc   uint16
  E_ovno     uint16
  E_res      [4]uint16
  E_oemid    uint16
  E_oeminfo  uint16
  E_res2     [10]uint16
  E_lfanew   int32 // file address of new exe header
}

type ImageFileHeader struct {
  Machine              uint16
  NumberOfSections     uint16
  TimeDateStamp        uint32
  PointerToSymbolTable uint32
  NumberOfSymbols      uint32
  SizeOfOptionalHeader uint16
  Characteristics      uint16
}