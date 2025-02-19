/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/hook"
)

func inject(process *os.Process, addons []Addon) {
  if addons != nil && len(addons) > 0 {
    for _, addon := range addons {
      if len(addon.Path) > 0 {
        dylib := fs.Resolve(expand.ExpandVariables(addon.Path))
        if fs.FileExist(dylib){
          err := hook.CreateRemoteThread(process.Pid, dylib)
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