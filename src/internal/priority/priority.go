/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package priority

import (
  "strings"
  "golang.org/x/sys/windows"
)

func GetPriorityClass(priority string) uint32 {

  var PRIORITY_CLASS = map[string]uint32{
    "IDLE":          windows.IDLE_PRIORITY_CLASS,
    "BELOW_NORMAL":  windows.BELOW_NORMAL_PRIORITY_CLASS,
    "NORMAL":        windows.NORMAL_PRIORITY_CLASS,
    "ABOVE_NORMAL":  windows.ABOVE_NORMAL_PRIORITY_CLASS,
    "HIGH":          windows.HIGH_PRIORITY_CLASS,
    "REALTIME":      windows.REALTIME_PRIORITY_CLASS,
  }

  if code, found := PRIORITY_CLASS[strings.ToUpper(priority)]; found {
    return code
  }

  return PRIORITY_CLASS["NORMAL"]
}