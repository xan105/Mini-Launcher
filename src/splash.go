/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "time"
  "slices"
  "strings"
  "math/rand"
  "path/filepath"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/splash"
)

func displaySplash(pid int, screen Splash) {
  if screen.Show && screen.Images != nil && len(screen.Images) > 0 {
    image := screen.Images[rand.Intn(len(screen.Images))]
    if len(image) > 0 {
      image = fs.Resolve(expand.ExpandVariables(image))
      if filepath.Ext(image) == ".bmp" {
        if ok, _ := fs.FileExist(image); ok {
        
          var timeout uint = 10
          if screen.Timeout > 0 {
            timeout = screen.Timeout
          }

          events := []string{"FOREGROUND", "WINDOW", "CURSOR"}
          var wait string = events[0]
          if slices.Contains(events, strings.ToUpper(screen.Wait)) {
            wait = strings.ToUpper(screen.Wait)
          }
        
          exit := make(chan bool)
          go splash.CreateWindow(exit, pid, image, wait)

          select {
            case <-exit:
              return
            case <-time.After(time.Second * time.Duration(timeout)):
              return
          }
        }
      }
    }
  }
}