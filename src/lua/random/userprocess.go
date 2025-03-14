/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package random

import (
  "math/rand"
  "unsafe"
  "time"
  "golang.org/x/sys/windows"
  "github.com/yuin/gopher-lua"
)

func getCurrentUserSID() (string, error) {

  hProcess, err := windows.GetCurrentProcess()
  if err != nil {
		return "", err
	}

	var token windows.Token
	if err := windows.OpenProcessToken(hProcess, windows.TOKEN_QUERY, &token); err != nil {
		return "", err
	}
	defer token.Close()

  user, err := token.GetTokenUser()
  if err != nil {
    return "", err
  }
  
  return user.User.Sid.String(), nil
}

func getProcessOwnerSID(pid uint32) (string, error) {
  hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
  if err != nil {
    return "", err
  }
  defer windows.CloseHandle(hProcess)

  var token windows.Token
  err = windows.OpenProcessToken(hProcess, windows.TOKEN_QUERY, &token)
  if err != nil {
    return "", err
  }
  defer token.Close()

  user, err := token.GetTokenUser()
  if err != nil {
    return "", err
  }

  return user.User.Sid.String(), nil
}

func listUserProcesses() ([]uint32, error) {
  handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
  if err != nil {
    return nil, err
  }
  defer windows.CloseHandle(handle)

  var entry windows.ProcessEntry32
  entry.Size = uint32(unsafe.Sizeof(entry))

  currentSID, err := getCurrentUserSID()
  if err != nil {
    return nil, err
  }

  var userPIDs []uint32
  if err := windows.Process32First(handle, &entry); err != nil {
    return nil, err
  }

  for {
    processSID, err := getProcessOwnerSID(entry.ProcessID)
    if err == nil && processSID == currentSID {
      userPIDs = append(userPIDs, entry.ProcessID)
    }
    if err := windows.Process32Next(handle, &entry); err != nil {
      break
    }
  }

  return userPIDs, nil
}

func GetRandomUserPID(L *lua.LState) int {
  pids, err := listUserProcesses()
  if err != nil {
    L.Push(lua.LNumber(0))
    return 1
  }
  
  if len(pids) == 0 {
    L.Push(lua.LNumber(0))
    return 1
  }

  rand.Seed(time.Now().UnixNano())
  pid := pids[rand.Intn(len(pids))]
  
  L.Push(lua.LNumber(pid))
  return 1
}