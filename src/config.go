/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

type Addon struct {
  Path           string               `json:"path"`
  Required       bool                 `json:"required"`
}

type File struct {
  Path           string               `json:"path"`
  SRI            string               `json:"sri"` 
  Size           int64                `json:"size"`
}

type Splash struct {
  Show            bool                `json:"show"`
  Images          []string            `json:"image"`
  Timeout         uint                `json:"timeout"`
  Wait            string              `json:"wait"`
}

type Link struct {
  Origin          string              `json:"path"`
  Destination     string              `json:"dest"`
}

type Attrib struct {
  Path           string               `json:"path"`
  Hidden         bool                 `json:"hidden"` 
  ReadOnly       bool                 `json:"readonly"`
}

type CompatFlags struct {
  Version         string              `json:"version"`
  Fullscreen      bool                `json:"fullscreen"`
  Admin           bool                `json:"admin"`
  Invoker         bool                `json:"invoker"`
  Aware           bool                `json:"aware"`
}

type WinePrefix struct {
  WinVer          string              `json:"winver"`
  DllOverrides    map[string]string   `json:"overrides"`
  DPI             uint32              `json:"dpi"`
}

type Config struct {
  Bin             string              `json:"bin"`
  Cwd             string              `json:"cwd"`
  Args            string              `json:"args"`
  Env             map[string]string   `json:"env"`
  Hide            bool                `json:"hide"`
  Shell           bool                `json:"shell"`
  Wait            bool                `json:"wait"`
  Script          string              `json:"script"`
  Addons          []Addon             `json:"addons"`
  Integrity       []File              `json:"integrity"`
  Splash          Splash              `json:"splash"`
  Symlink         []Link              `json:"symlink"`
  Compatibility   CompatFlags         `json:"compatibility"`
  Prefix          WinePrefix          `json:"prefix"`
  Attrib          []Attrib            `json:"attrib"`
}