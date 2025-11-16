/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package thread

import (
  "unsafe"
  "golang.org/x/sys/windows"
)

var (
  kernel32            = windows.NewLazySystemDLL("kernel32.dll")
  pVirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
  pVirtualFreeEx      = kernel32.NewProc("VirtualFreeEx")
  pCreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
)

func CreateRemoteThread(pid int, path string) error {

  //Opens a handle to the target process with the needed permissions
  hProcess, err := windows.OpenProcess(
    windows.PROCESS_CREATE_THREAD | 
    windows.PROCESS_VM_OPERATION | 
    windows.PROCESS_VM_WRITE | 
    windows.PROCESS_VM_READ |
    windows.PROCESS_QUERY_INFORMATION,
    false,
    uint32(pid),
  )
  if err != nil {
    return err
  }
  defer windows.CloseHandle(hProcess)

 //Allocates virtual memory for the file path
  lpBaseAddress, _, err := pVirtualAllocEx.Call(
    uintptr(hProcess), 
    0, 
    uintptr((len(path) + 1) * 2),
    windows.MEM_RESERVE | windows.MEM_COMMIT, 
    windows.PAGE_EXECUTE_READWRITE,
  )
 
  //Converts the file path to type LPCWSTR
  lpBuffer, err := windows.UTF16PtrFromString(path)
  if err != nil {
    return err
  }
 
 //Writes the filename to the previously allocated space
  lpNumberOfBytesWritten:= uintptr(0)
  err = windows.WriteProcessMemory(
    hProcess, 
    lpBaseAddress, 
    (*byte)(unsafe.Pointer(lpBuffer)),
    uintptr((len(path) + 1) * 2),
    &lpNumberOfBytesWritten,
  )
  if err != nil {
    return err
  }
 
 //Gets a pointer to the LoadLibrary function
  LoadLibAddr, err := windows.GetProcAddress(
    windows.Handle(kernel32.Handle()), 
    "LoadLibraryW",
  )
  if err != nil {
    return err
  }
 
 //Creates a remote thread that loads the DLL triggering it
  hThread, _, err := pCreateRemoteThread.Call(
    uintptr(hProcess), 
    0, 
    0, 
    LoadLibAddr, 
    lpBaseAddress, 
    0, 
    0,
  )
  if hThread == 0 {
    return err
  }
  defer windows.CloseHandle(windows.Handle(hThread))

  windows.WaitForSingleObject(windows.Handle(hThread), windows.INFINITE)

  pVirtualFreeEx.Call(
    uintptr(hProcess), 
    lpBaseAddress, 
    0,
    windows.MEM_RELEASE,
  )

  return nil
}

func ResumeThread(tid uint32) error {

  hThread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, tid)
  if err != nil {
    return err
  }
  defer windows.CloseHandle(windows.Handle(hThread))

  if _, err := windows.ResumeThread(hThread); err != nil {
    return err
  }
  
  windows.WaitForSingleObject(windows.Handle(hThread), windows.INFINITE)
  
  return nil
}

func ResumeMainThread(pid int) error {

  /*
  It is worth mentioning that this function resume the first thread found of the specified process.
  Which _should_ be the main thread. But technically this is not _"correct"_.
  
  Why using `CreateToolhelp32Snapshot()` then you might ask?
  Long story short, this is because `os/exec` does not return the handle of the main thread from `CreateProcessW()`.
  And for once, I didn't feel like re-inventing the wheel, ie: doing my own os/exec and/or wrapper of `CreateProcessW()`, 
  just to be able to create a suspended process and resume it.
  */

  hSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, uint32(pid))
  if err != nil {
    return err
  }
  defer windows.CloseHandle(windows.Handle(hSnapshot))

  var entry windows.ThreadEntry32
  entry.Size = uint32(unsafe.Sizeof(entry))
  if err := windows.Thread32First(hSnapshot, &entry); err != nil {
    return err
  }

  for
  {
    if err := windows.Thread32Next(hSnapshot, &entry); err != nil {
      return err
    }
    
    if entry.OwnerProcessID == uint32(pid) && entry.ThreadID != 0 {
      return ResumeThread(entry.ThreadID)
    }
  }
}