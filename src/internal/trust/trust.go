/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//cf: https://github.com/xan105/node-win-verify-trust

package trust

import (
  "unsafe"
  "golang.org/x/sys/windows"
)

const (
  WTD_REVOKE_NONE                   = 0
  WTD_CHOICE_FILE                   = 1
  WTD_UI_NONE                       = 2
  WTD_STATEACTION_VERIFY            = 0x00000001
  WTD_STATEACTION_CLOSE             = 0x00000002
)

type WINTRUST_FILE_INFO struct {
  CbStruct                          uint32
  PcwszFilePath                     *uint16
  HFile                             windows.Handle
  PgKnownSubject                    *windows.GUID
}

type WINTRUST_DATA struct {
  CbStruct                          uint32
  PPolicyCallbackData               uintptr
  PSIPClientData                    uintptr
  DwUIChoice                        uint32
  FdwRevocationChecks               uint32
  DwUnionChoice                     uint32
  PFile                             *WINTRUST_FILE_INFO
  DwStateAction                     uint32
  H_WVTStateData                    windows.Handle
  PwszURLReference                  *uint16
  DwUIContext                       uint32
}

var (
  wintrust                          = windows.NewLazySystemDLL("wintrust.dll")
  pWinVerifyTrust                   = wintrust.NewProc("WinVerifyTrust")
  WINTRUST_ACTION_GENERIC_VERIFY_V2 = windows.GUID{
    Data1: 0xaac56b,
    Data2: 0xcd44,
    Data3: 0x11d0,
    Data4: [8]byte{0x8c, 0xc2, 0x0, 0xc0, 0x4f, 0xc2, 0x95, 0xee},
  }
)

func VerifySignature(filePath string) (bool, error) {
  filePathPtr, err := windows.UTF16PtrFromString(filePath)
  if err != nil {
    return false, err
  }

  fileInfo := WINTRUST_FILE_INFO{
    CbStruct:      uint32(unsafe.Sizeof(WINTRUST_FILE_INFO{})),
    PcwszFilePath: filePathPtr,
    HFile:         0,
    PgKnownSubject: nil,
  }

  winTrustData := WINTRUST_DATA{
    CbStruct:            uint32(unsafe.Sizeof(WINTRUST_DATA{})),
    PPolicyCallbackData: 0,
    PSIPClientData:      0,
    DwUIChoice:          WTD_UI_NONE,
    FdwRevocationChecks: WTD_REVOKE_NONE,
    DwUnionChoice:       WTD_CHOICE_FILE,
    PFile:               &fileInfo,
    DwStateAction:       WTD_STATEACTION_VERIFY,
    H_WVTStateData:      0,
    PwszURLReference:    nil,
    DwUIContext:         0,
  }

  ret, _, err := pWinVerifyTrust.Call(
    0, 
    uintptr(unsafe.Pointer(&WINTRUST_ACTION_GENERIC_VERIFY_V2)), 
    uintptr(unsafe.Pointer(&winTrustData)),
  )
  
  // Any H_WVTStateData must be released by a call with close.
  winTrustData.DwStateAction = WTD_STATEACTION_CLOSE
  pWinVerifyTrust.Call(
    0, 
    uintptr(unsafe.Pointer(&WINTRUST_ACTION_GENERIC_VERIFY_V2)), 
    uintptr(unsafe.Pointer(&winTrustData)),
  )

  if ret == 0 {
    return true, nil
  }
  
  return false, err
}