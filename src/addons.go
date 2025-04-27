/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "strings"
  "runtime"
  "path/filepath"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/pe"
  "launcher/internal/thread"
)

func loadAddons(binary string, process *os.Process, addons []Addon) {
  if addons != nil && len(addons) > 0 {
    targetArch, err := pe.GetArchFromMachineType(binary)
    for _, addon := range addons {
      if len(addon.Path) > 0 {
        dylib := fs.Resolve(expand.ExpandVariables(addon.Path))
        ext := strings.ToLower(filepath.Ext(dylib))
        if ext == ".dll" || ext == ".asi" {
          if ok, _ := fs.FileExist(dylib); ok {

            if targetArch != runtime.GOARCH {
              if addon.Required {
                process.Kill()
                if err != nil {
                  panic("Remote Thread", "\"" + filepath.Base(binary) + "\": " + err.Error())
                } else {
                  panic("Remote Thread", "\"" + filepath.Base(binary)  + "\" and the Launcher are of different architecture!")
                }
              } else {
                continue
              } 
            }
            
            if arch, err := pe.GetArchFromMachineType(dylib); arch != runtime.GOARCH {
              if addon.Required {
                process.Kill()
                if err != nil {
                  panic("Remote Thread", "\"" + filepath.Base(dylib) + "\": " + err.Error())
                } else {
                  panic("Remote Thread", "\"" + filepath.Base(dylib)  + "\" and the target process are of different architecture!")
                }
              } else {
                continue
              }
            }
   
            if err := thread.CreateRemoteThread(process.Pid, dylib); err != nil {
              if addon.Required {
                process.Kill()
                panic("Remote Thread", "\"" + filepath.Base(dylib) + "\": " + err.Error())
              } else {
                continue
              }
            }
            
          }
        }
      }
    }
  }
}