/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

type Addon struct {
  Path      string                `json:"path"`
  Required  bool                  `json:"required"`
}

type File struct {
  Path      string                `json:"path"`
  SRI       string                `json:"sri"` 
  Size      int64                 `json:"size"`
}

type Splash struct {
  Show        bool                `json:"show"`
  Images      []string            `json:"image"`
  Timeout     uint                `json:"timeout"`
}

type Link struct {
  Origin      string              `json:"path"`
  Destination string              `json:"dest"`
}

type Config struct {
  Bin         string              `json:"bin"`
  Cwd         string              `json:"cwd"`
  Args        string              `json:"args"`
  Env         map[string]string   `json:"env"`
  Hide        bool                `json:"hide"`
  Script      string              `json:"script"`
  Addons      []Addon             `json:"addons"`
  Integrity   []File              `json:"integrity"`
  Splash      Splash              `json:"splash"`
  Symlink     []Link              `json:"symlink"`
}