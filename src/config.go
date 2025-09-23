/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

type Addon struct {
  Path            string               `json:"path"`
  Required        bool                 `json:"required"`
}

type File struct {
  Path            string               `json:"path"`
  SRI             string               `json:"sri"` 
  Size            int64                `json:"size"`
  Signed          bool                 `json:"signed"`
}

type Splash struct {
  Show            *bool               `json:"show"`
  Images          []string            `json:"image"`
  Timeout         uint                `json:"timeout"`
  Wait            string              `json:"wait"`
}

type Link struct {
  Origin          string              `json:"path"`
  Destination     string              `json:"dest"`
}

type Attrib struct {
  Path            string               `json:"path"`
  Hidden          bool                 `json:"hidden"` 
  ReadOnly        bool                 `json:"readonly"`
}

type CompatFlags struct {
  Version         string              `json:"version"`
  Fullscreen      *bool               `json:"fullscreen"`
  Admin           *bool               `json:"admin"`
  Invoker         *bool               `json:"invoker"`
  Aware           *bool               `json:"aware"`
}

type Patch struct {
  LAA             *bool               `json:"laa"`
}

type WinePrefix struct {
  WinVer          string              `json:"winver"`
  DllOverrides    map[string]string   `json:"overrides"`
  DPI             uint32              `json:"dpi"`
}

type Script struct {
  Path            string              `json:"path"`
  Fs              *bool               `json:"fs"`
  Net             *bool               `json:"net"`
  Reg             *bool               `json:"reg"`
  Exec            *bool               `json:"exec"`
  Import          *bool               `json:"import"`
}

type Shortcut struct {
  Name            string              `json:"name"`
  Desktop         *bool               `json:"desktop"`
  StartMenu       *bool               `json:"startmenu"`
}

type Config struct {
  Bin             string              `json:"bin"`
  Cwd             string              `json:"cwd"`
  Args            string              `json:"args"`
  Priority        string              `json:"priority"`
  Env             map[string]string   `json:"env"`
  Hide            *bool               `json:"hide"`
  Shell           *bool               `json:"shell"`
  Wait            *bool               `json:"wait"`
  Script          Script              `json:"script"`
  Addons          []Addon             `json:"addons"`
  Integrity       []File              `json:"integrity"`
  Splash          Splash              `json:"splash"`
  Symlink         []Link              `json:"symlink"`
  Compatibility   CompatFlags         `json:"compatibility"`
  Patch           Patch               `json:"patch"`
  Prefix          WinePrefix          `json:"prefix"`
  Attrib          []Attrib            `json:"attrib"`
  Menu            map[string]string   `json:"menu"`
  Shortcut        Shortcut            `json:"shortcut"`
}

func mergeConfig(config *Config, override *Config) {
  
  //string
  if len(override.Bin) > 0 { config.Bin = override.Bin }
  if len(override.Cwd) > 0 { config.Cwd = override.Cwd }
  if len(override.Args) > 0 { config.Args = override.Args }
  if len(override.Priority) > 0 { config.Priority = override.Priority }
  
  //bool
  if override.Hide != nil { config.Hide = override.Hide }
  if override.Shell != nil { config.Shell = override.Shell }
  if override.Wait != nil { config.Wait = override.Wait }
  
  //map
  if len(override.Env) > 0 {
    for k, v := range override.Env { config.Env[k] = v }
  }
  
  //[]struct
  if override.Addons != nil && len(override.Addons) > 0 { config.Addons = override.Addons }
  if override.Integrity != nil && len(override.Integrity) > 0 { config.Integrity = override.Integrity }
  if override.Symlink != nil && len(override.Symlink) > 0 { config.Symlink = override.Symlink }
  if override.Attrib != nil && len(override.Attrib) > 0 { config.Attrib = override.Attrib }
  
  //Nested
  if len(override.Script.Path) > 0 { config.Script.Path = override.Script.Path }
  if override.Script.Fs != nil { config.Script.Fs = override.Script.Fs }
  if override.Script.Net != nil { config.Script.Net = override.Script.Net }
  if override.Script.Reg != nil { config.Script.Reg = override.Script.Reg }
  if override.Script.Exec != nil { config.Script.Exec = override.Script.Exec }
  if override.Script.Import != nil { config.Script.Import = override.Script.Import }
  
  if len(override.Shortcut.Name) > 0 { config.Shortcut.Name = override.Shortcut.Name }
  if override.Shortcut.Desktop != nil { config.Shortcut.Desktop = override.Shortcut.Desktop }
  if override.Shortcut.StartMenu != nil { config.Shortcut.StartMenu = override.Shortcut.StartMenu }

  if override.Splash.Show != nil { config.Splash.Show = override.Splash.Show }
  if override.Splash.Images != nil && len(override.Splash.Images) > 0 { config.Splash.Images = override.Splash.Images }
  if override.Splash.Timeout > 0 { config.Splash.Timeout = override.Splash.Timeout }
  if len(override.Splash.Wait) > 0 { config.Splash.Wait = override.Splash.Wait }
  
  if len(override.Compatibility.Version) > 0 { config.Compatibility.Version = override.Compatibility.Version }
  if override.Compatibility.Fullscreen != nil  { config.Compatibility.Fullscreen = override.Compatibility.Fullscreen }
  if override.Compatibility.Admin != nil { config.Compatibility.Admin = override.Compatibility.Admin }
  if override.Compatibility.Invoker != nil { config.Compatibility.Invoker = override.Compatibility.Invoker }
  if override.Compatibility.Aware != nil { config.Compatibility.Aware = override.Compatibility.Aware }
  
  if override.Patch.LAA != nil { config.Patch.LAA = override.Patch.LAA }
  
  if len(override.Prefix.WinVer) > 0 { config.Prefix.WinVer = override.Prefix.WinVer }
  if override.Prefix.DPI > 0 { config.Prefix.DPI = override.Prefix.DPI }
  if len(override.Prefix.DllOverrides) > 0 {
    for k, v := range override.Prefix.DllOverrides {
      config.Prefix.DllOverrides[k] = v
    }
  }
}