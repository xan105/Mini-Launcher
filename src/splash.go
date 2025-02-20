/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "time"
  "math/rand"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/splash"
)

func displaySplash(pid int, screen Splash) {
  if screen.Show && screen.Images != nil && len(screen.Images) > 0 {
    image := screen.Images[rand.Intn(len(screen.Images))]
    if len(image) > 0 {
      image = fs.Resolve(expand.ExpandVariables(image))
      if fs.FileExist(image){
        var timeout uint = 10
        if screen.Timeout > 0 {
          timeout = screen.Timeout
        }
      
        exit := make(chan bool)
        go splash.CreateWindow(exit, pid, image)

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