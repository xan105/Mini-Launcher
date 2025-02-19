/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "slices"
  "strings"
  "launcher/internal/regedit"
)

func setAppCompatFlags(binary string, flags CompatFlags) {

  //Win10: "~ [Fullscreen Optimization] + [Privilege Level] + [Flags...] + [Compatibility Mode]"

  const path = "Software/Microsoft/Windows NT/CurrentVersion/AppCompatFlags/Layers"

  if len(flags.Version) > 0 || flags.Fullscreen || flags.Admin || flags.Aware {

    version := []string{
      "WIN95",
      "WIN98", 
      "WIN2000",
      "WINXP", 
      "WINXPSP1",
      "WINXPSP2",
      "WINXPSP3",
      "VISTARTM",
      "VISTASP1",
      "VISTASP2",
      "WIN7RTM",
      "WIN8RTM",
    }

    template := []string{ "~" }
    
    if flags.Fullscreen {
      template = append(template, "DISABLEDXMAXIMIZEDWINDOWEDMODE")
    }

    if flags.Admin {
      template = append(template, "RUNASADMIN")
    }
    
    if flags.Aware {
      template = append(template, "HIGHDPIAWARE")
    }
    
    if slices.Contains(version, flags.Version) {
      template = append(template, flags.Version)
    }

    regedit.WriteStringValue("HKCU", path, binary, strings.Join(template, " "))
    
  } else { regedit.DeleteKeyValue("HKCU", path, binary) }
}