/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "path/filepath"
  "launcher/internal/pe"
)

func applyPatches(binary string, patches Patch) {
  if !patches.Allow { return }
  
  if err := pe.PatchLargeAddress(binary, patches.LAA); err != nil {
    panic("Patch (Large Address Aware)", "\"" + filepath.Base(binary) + "\": " + err.Error())
  }
}