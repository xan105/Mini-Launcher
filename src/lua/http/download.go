/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package http

import (
  "io"
  "os"
  "mime"
  "strings"
  "net/http"
  "path/filepath"
  "github.com/yuin/gopher-lua"
  "launcher/internal/fs"
  "launcher/internal/expand"
)

func Download(L *lua.LState) int {
  url := L.CheckString(1)
  
  destDir := L.CheckString(2)
  if len(destDir) > 0{
    destDir = fs.Resolve(expand.ExpandVariables(destDir))
  } else {
    L.Push(lua.LNil)
    L.Push(lua.LString("Destination dir is empty!"))
    return 1
  }
  
  resp, err := http.Get(url)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }
  defer resp.Body.Close()

  //Determine the filename
  parts := strings.Split(url, "/")
  last := len(parts)-1
  filename := parts[last]
  if cd := resp.Header.Get("Content-Disposition"); cd != "" {
    if _, params, err := mime.ParseMediaType(cd); err == nil {
      if name, ok := params["filename"]; ok {
        filename = name
      }
    }
  }

  //Create
  if err := os.MkdirAll(destDir, 0755); err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }
  filePath := filepath.Join(destDir, filename)
  out, err := os.Create(filePath)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }
  defer out.Close()

  //Write the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(lua.LString(err.Error()))
    return 2
  }

  L.Push(lua.LString(filePath))
  return 1
}