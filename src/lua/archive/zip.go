/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package archive

import (
  "os"
  "io"
  "strings"
  "archive/zip"
  "path/filepath"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Unzip": Unzip,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Unzip(L *lua.LState) int {

  path := L.CheckString(1)
  if len(path) > 0{
    path = fs.Resolve(expand.ExpandVariables(path))
    if filepath.Ext(path) != ".zip" {
      L.Push(lua.LString("Not a .zip file !"))
      return 1
    }
  } else {
    L.Push(lua.LString("Archive file path is empty!"))
    return 1
  }

  destDir := L.CheckString(2)
  if len(destDir) > 0{
    destDir = fs.Resolve(expand.ExpandVariables(destDir))
  } else {
    L.Push(lua.LString("Destination dir is empty!"))
    return 1
  }

  r, err := zip.OpenReader(path)
  if err != nil {
    L.Push(lua.LString(err.Error()))
    return 1
  }
  defer r.Close()

  for _, file := range r.File {
    fpath := filepath.Join(destDir, file.Name)

    if file.FileInfo().IsDir() {
      os.MkdirAll(fpath, os.ModePerm)
      continue
    }

    if !strings.HasPrefix(fpath, filepath.Clean(destDir) + string(os.PathSeparator)) {
      continue
    }

    if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
      L.Push(lua.LString(err.Error()))
      return 1
    }

    outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
    if err != nil {
      L.Push(lua.LString(err.Error()))
      return 1
    }
    defer outFile.Close()

    rc, err := file.Open()
    if err != nil {
      L.Push(lua.LString(err.Error()))
      return 1
    }
    defer rc.Close()

    _, err = io.Copy(outFile, rc)
    if err != nil {
      L.Push(lua.LString(err.Error()))
      return 1
    }
  }

  return 0
}