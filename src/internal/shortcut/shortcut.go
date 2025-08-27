/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package shortcut

import (
  "runtime"
  "github.com/go-ole/go-ole"
  "github.com/go-ole/go-ole/oleutil" 
)

type Shortcut struct {
  Path              string
  TargetPath        string
  Arguments         string
  WorkingDirectory  string
  Description       string
  IconLocation      string
}

func CreateShortcut(s Shortcut) (error){
  runtime.LockOSThread()
  defer runtime.UnlockOSThread()
  
  err := ole.CoInitialize(0)
  if err != nil {
    return err
  }
  defer ole.CoUninitialize()
  
  com, err := oleutil.CreateObject("WScript.Shell")
  if err != nil {
    return err
  }
  defer com.Release()
  
  wshell, err := com.QueryInterface(ole.IID_IDispatch)
  if err != nil {
    return err
  }
  defer wshell.Release()
  
  shortcut, err := oleutil.CallMethod(wshell, "CreateShortcut", s.Path)
  if err != nil {
    return err
  }
  sc := shortcut.ToIDispatch()
  defer sc.Release()
  
  if _, err := oleutil.PutProperty(sc, "TargetPath", s.TargetPath); err != nil {
    return err
  }
  if _, err := oleutil.PutProperty(sc, "Arguments", s.Arguments); err != nil {
    return err
  }
  if _, err := oleutil.PutProperty(sc, "WorkingDirectory", s.WorkingDirectory); err != nil {
    return err
  }
  if _, err := oleutil.PutProperty(sc, "Description", s.Description); err != nil {
    return err
  }
  if _, err := oleutil.PutProperty(sc, "IconLocation", s.IconLocation); err != nil {
    return err
  }
  
  _, err = oleutil.CallMethod(sc, "Save")
  if err != nil {
    return err;
  }
  
  return nil
}