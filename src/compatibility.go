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
  "launcher/internal/wine"
)

func setCompatFlags(binary string, flags CompatFlags) {

  //Win10: "~ [Fullscreen Optimization] + [Privilege Level] + [Flags...] + [Compatibility Mode]"

  const path = "Software/Microsoft/Windows NT/CurrentVersion/AppCompatFlags/Layers"
  template := []string{}
    
  if flags.Fullscreen != nil && *flags.Fullscreen {
    template = append(template, "DISABLEDXMAXIMIZEDWINDOWEDMODE")
  }

  if flags.Admin != nil && *flags.Admin {
    template = append(template, "RUNASADMIN")
  } else if flags.Invoker != nil && *flags.Invoker {
    template = append(template, "RUNASINVOKER")
  }
    
  if flags.Aware != nil && *flags.Aware {
    template = append(template, "HIGHDPIAWARE")
  }
    
  if len(flags.Version) > 0 {
    versions := []string{
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
    if slices.Contains(versions, flags.Version) {
      template = append(template, flags.Version)
    }
  }

  if len(template) > 0 {
    slices.Insert(template, 0, "~")
    regedit.WriteStringValue("HKCU", path, binary, strings.Join(template, " "))
  } else {
    regedit.DeleteValue("HKCU", path, binary)
  }
}

func updatePrefixSettings(prefix WinePrefix) {

  if !wine.IsWineOrProton() { return }

  if len(prefix.WinVer) > 0 {
    versions := []string{
      "win11",
      "win10",
      "win81", 
      "win8", 
      "win7",  
      "vista", 
      "winxp",
    }
    if slices.Contains(versions, prefix.WinVer) {
      regedit.WriteStringValue("HKCU", "HKCU/Software/Wine", "Version", prefix.WinVer)
    }
  }
  
  if prefix.DPI >= 96 && prefix.DPI <= 480 {
    regedit.WriteDwordValue("HKCU", "Control Panel/Desktop", "LogPixels", prefix.DPI)
  }
  
  if len(prefix.DllOverrides) > 0 {
    overrides := []string{
      "native,builtin",
      "builtin,native",
      "native",
      "builtin",
    }
    for dll, override := range prefix.DllOverrides {
      if len(dll) > 0 && len(override) > 0 {
        if slices.Contains(overrides, override) {
          regedit.WriteStringValue("HKCU", "HKCU/Software/Wine/DllOverrides", strings.ToLower(dll), override)
        }
      }
    }
  }
}