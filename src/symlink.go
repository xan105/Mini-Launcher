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

func makeLink(links []Link) {
  if links != nil && len(links) > 0 {
    for _, link := range links {
      if len(link.Origin) > 0 && len(link.Destination) > 0 {
        err := fs.CreateFolderSymlink(
          fs.Resolve(expand.ExpandVariables(link.Origin)),
          fs.Resolve(expand.ExpandVariables(link.Destination)),
        )
        if err != nil { 
          panic("Symlink", err.Error())
        } 
      }
    }
  } 
}