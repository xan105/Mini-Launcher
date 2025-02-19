/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "strings"
  "regexp"
  "errors"
  "launcher/internal/fs"
  "launcher/internal/expand"
)

func verifyIntegrity(binary string, files []File) {
  if files != nil && len(files) > 0 {
    for _, file := range files {
      
      target := binary
      if len(file.Path) > 0 {
        target = fs.Resolve(expand.ExpandVariables(file.Path))
      }
      
      stats, err := os.Stat(target)
      if err != nil { 
        if errors.Is(err, os.ErrNotExist) {
          panic("Integrity failure", "File does not exist: \"" + target + "\"") 
        }
        panic("Integrity failure", err.Error())  
      }
      if file.Size > 0 {
        if stats.Size() != file.Size { 
          panic("Integrity failure", "Size mismatch: \"" + target + "\"") 
        }
      }
      
      re := regexp.MustCompile(`(?i)^(sha(?:256|384|512)-[A-Za-z0-9+/=]{43,}={0,2})$`)
      if !re.MatchString(file.SRI) {
        panic("Integrity failure", "Unexpected SRI format: \"" + file.SRI + "\"")
      }
      
      SRI := strings.SplitN(file.SRI, "-", 2)
      if len(SRI) != 2 {
        panic("Integrity failure", "Failed to parse SRI: \"" + file.SRI + "\"")
      }

      algo, expected := SRI[0], SRI[1]
      sum, err := fs.CheckSum(target, algo)
      if err != nil { panic("Integrity failure", err.Error()) }
      if sum != expected { 
        panic("Integrity failure", 
              "Hash mismatch: \"" + target + "\"\n" +
              "SRI: " + algo + "-" + sum)
      }
    }
  }
}