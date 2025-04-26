/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "launcher/internal/fs"
  "launcher/internal/expand"
)

func applyFileAttributes(attributes []Attrib) {
  if attributes != nil && len(attributes) > 0 {
    for _, attribute := range attributes {
      if len(attribute.Path) > 0 {
        file := fs.Resolve(expand.ExpandVariables(attribute.Path))
        if ok, _ := fs.FileExist(file); ok {   
          err := fs.SetFileAttributes(file, attribute.ReadOnly, attribute.Hidden)
          if err != nil {
            panic("File Attribute", err.Error())
          }
        }
      }
    }
  }
}