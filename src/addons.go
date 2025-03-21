/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "path/filepath"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/thread"
)

func loadAddons(process *os.Process, addons []Addon) {
  if addons != nil && len(addons) > 0 {
    for _, addon := range addons {
      if len(addon.Path) > 0 {
        dylib := fs.Resolve(expand.ExpandVariables(addon.Path))
        if filepath.Ext(dylib) == ".dll" {
          if ok, _ := fs.FileExist(dylib); ok {       
            err := thread.CreateRemoteThread(process.Pid, dylib)
            if err != nil {
              if addon.Required {
                process.Kill()
                panic("Remote Thread", err.Error())
              }
            }
          }
        }
      }
    }
  }
}