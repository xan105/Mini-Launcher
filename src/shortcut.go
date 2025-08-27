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
  "launcher/internal/shortcut"
)

func makeShortcut(binary string, s Shortcut) {
  if len(s.Name) > 0 {
    
    self, err := os.Executable()
    if err != nil { 
      panic("Creating Shortcut", err.Error()) 
    }
    
    if s.StartMenu != nil && *s.StartMenu {
      path := fs.Resolve(expand.ExpandVariables("%APPDATA%/Microsoft/Windows/Start Menu/Programs/" + s.Name + ".lnk"))
      if ok, err := fs.FileExist(path); !ok {
        if err != nil { 
          panic("Creating Shortcut ()", "Start Menu: " + err.Error())
        }

        mslink := shortcut.Shortcut{
          Path: path,
          TargetPath: self,
          WorkingDirectory: filepath.Dir(self),
          IconLocation: binary + ",0",
        }

        if err := shortcut.CreateShortcut(mslink); err != nil {
          panic("Creating Shortcut ()", "Start Menu: " + err.Error()) 
        }
      }
    }

    if s.Desktop != nil && *s.Desktop {
      path := fs.Resolve(expand.ExpandVariables("%DESKTOP%/" + s.Name + ".lnk"))
      if ok, err := fs.FileExist(path); !ok {
        if err != nil { 
          panic("Creating Shortcut ()", "Desktop: " + err.Error())
        }

        mslink := shortcut.Shortcut{
          Path: path,
          TargetPath: self,
          WorkingDirectory: filepath.Dir(self),
          IconLocation: binary + ",0",
        }

        if err := shortcut.CreateShortcut(mslink); err != nil {
          panic("Creating Shortcut ()", "Desktop: " + err.Error()) 
        }
      }
    }
    
  }
}