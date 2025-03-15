/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "errors"
  "golang.org/x/sys/windows"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/regedit"
  "launcher/internal/elevated"
)

func makeLink(links []Link) {
  if links != nil && len(links) > 0 {
    for _, link := range links {
      if len(link.Origin) > 0 && len(link.Destination) > 0 {
        err := fs.CreateFolderSymlink(
          fs.Resolve(expand.ExpandVariables(link.Origin)),
          fs.Resolve(expand.ExpandVariables(link.Destination)),
        )
        if err != nil {
          if linkErr, ok := err.(*os.LinkError); ok {
            if linkErr.Op == "symlink" {
              if !regedit.KeyExists("HCKU", "HKCU/Software/Wine") { //We don't do that here
                if errors.Is(linkErr.Err, windows.Errno(windows.ERROR_ACCESS_DENIED)) || 
                  errors.Is(linkErr.Err, windows.Errno(windows.ERROR_PRIVILEGE_NOT_HELD)) { 
                  if !elevated.IsElevated(){
                      elevated.RestartElevated()
                  }
                }
              }
            }
          }
          panic("Symlink", err.Error())
        } 
      }
    }
  } 
}