/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "launcher/internal/ui"
)
  
func displayMenuOverride(menu map[string]string, defaultPath string) string {
  if len(menu) > 0 {
    labels := []string{}
    for label, _ := range menu {
      if len(label) > 0 {
        labels = append(labels, label)
      }
    }
    if len(labels) > 0 {
      button := ui.Menu(labels)
      index := <- button
      if index >= 0 && index <= len(labels) {
        name := labels[index]
        path := menu[name]
        if len(path) == 0 { path = defaultPath }
        return path
      }
    }
  }
  return ""
}