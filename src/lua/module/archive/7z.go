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
  "path/filepath"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/lua/type/failure"
  "github.com/yuin/gopher-lua"
  "github.com/bodgit/sevenzip"
)

func Un7z(L *lua.LState) int {

  path := L.CheckString(1)
  if len(path) > 0 {
    path = fs.Resolve(expand.ExpandVariables(path))
    if filepath.Ext(path) != ".7z" {
      L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", "Not a .7z file !"))
      return 1
    }
  } else {
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", "Archive file path is empty!"))
    return 1
  }

  destDir := L.CheckString(2)
  if len(destDir) > 0{
    destDir = fs.Resolve(expand.ExpandVariables(destDir))
  } else {
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", "Destination dir is empty!"))
    return 1
  }

  excludeList := []string{}
  if L.GetTop() >= 3 {
    list := L.CheckTable(3)
    list.ForEach(func(_, value lua.LValue) {
      if str, ok := value.(lua.LString); ok {
        excludeList = append(excludeList, filepath.FromSlash(string(str)))
      }
    })
  }

  r, err := sevenzip.OpenReader(path)
  if err != nil {
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 1
  }
  defer r.Close()

  for _, file := range r.File {
    fpath := filepath.Join(destDir, file.Name)

    skip := false
    for _, exclude := range excludeList {
      if strings.HasPrefix(filepath.FromSlash(file.Name), exclude) {
        skip = true
        break
      }
    }
    if skip { continue }
    
    if file.FileInfo().IsDir() {
      os.MkdirAll(fpath, os.ModePerm)
      continue
    }

    if !strings.HasPrefix(fpath, filepath.Clean(destDir) + string(os.PathSeparator)) {
      continue
    }

    if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
      L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
      return 1
    }

    outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
    if err != nil {
      L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
      return 1
    }
    defer outFile.Close()

    rc, err := file.Open()
    if err != nil {
      L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
      return 1
    }
    
    _, err = io.Copy(outFile, rc)
    rc.Close()
    if err != nil {
      L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
      return 1
    } 
  }

  return 0
}