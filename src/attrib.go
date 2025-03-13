/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "golang.org/x/sys/windows"
  "launcher/internal/fs"
  "launcher/internal/expand"
)

func toggleFileAttributes(filePath string, readonly bool, hidden bool) error {
  attrs, err := windows.GetFileAttributes(windows.StringToUTF16Ptr(filePath))
  if err != nil { return err }

  if readonly {
    if attrs & windows.FILE_ATTRIBUTE_READONLY != 0 {
      attrs &^= windows.FILE_ATTRIBUTE_READONLY //remove
    } else {
      attrs |= windows.FILE_ATTRIBUTE_READONLY //add
    }
  }

  if hidden {
    if attrs & windows.FILE_ATTRIBUTE_HIDDEN != 0 {
      attrs &^= windows.FILE_ATTRIBUTE_HIDDEN //remove
    } else {
      attrs |= windows.FILE_ATTRIBUTE_HIDDEN //add
    }
  }

  err = windows.SetFileAttributes(windows.StringToUTF16Ptr(filePath), attrs)
  if err != nil { return err }
  
  return nil
}

func applyFileAttributes(attributes []Attrib) {
  if attributes != nil && len(attributes) > 0 {
    for _, attribute := range attributes {
      if len(attribute.Path) > 0 {
        file := fs.Resolve(expand.ExpandVariables(attribute.Path))
        if ok, _ := fs.FileExist(file); ok {   
          err := toggleFileAttributes(file, attribute.ReadOnly, attribute.Hidden)
          if err != nil {
            panic("File Attribute", err.Error())
          }
        }
      }
    }
  }
}