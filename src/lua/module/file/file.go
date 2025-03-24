/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package file

import (
  "strings"
  "path/filepath"
  "github.com/yuin/gopher-lua"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/version"
  "launcher/lua/type/failure"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Write": Write,
    "Read": Read,
    "Remove": Remove,
    "Version": Version,
    "Glob": Glob,
    "Basename": Basename,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Write(L *lua.LState) int {
  filename := L.ToString(1)  
  data     := L.ToString(2)
  format   := L.ToString(3)

  if len(format) == 0 {
    format = "utf8"
  } 
 
  err := fs.WriteFile(
    fs.Resolve(expand.ExpandVariables(filename)), 
    data, 
    format,
  )
  if err != nil {
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 1
  }
    
  return 0
}

func Read(L *lua.LState) int {
  filename := L.ToString(1)  
  format   := L.ToString(2)

  if len(format) == 0 {
    format = "utf8"
  } 
  
  data, err := fs.ReadFile(
    fs.Resolve(expand.ExpandVariables(filename)), 
    format,
  )
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 2
  }

  L.Push(lua.LString(data))
  return 1
}

func Remove(L *lua.LState) int {
  path := L.ToString(1)
  
  err := fs.Remove(fs.Resolve(expand.ExpandVariables(path)))
  if err != nil {
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 1
  }
  
  return 0
}

func Version(L *lua.LState) int {
  filename := L.ToString(1)  

  fileInfo, err := version.FromFile(
    fs.Resolve(expand.ExpandVariables(filename)),
  )
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 2
  }

  fileVersion := L.NewTable()
  L.SetField(fileVersion, "Major", lua.LNumber(fileInfo.Major))
  L.SetField(fileVersion, "Minor", lua.LNumber(fileInfo.Minor))
  L.SetField(fileVersion, "Build", lua.LNumber(fileInfo.Build))
  L.SetField(fileVersion, "Revision", lua.LNumber(fileInfo.Revision))

  L.Push(fileVersion)
  return 1
}

func Glob(L *lua.LState) int {
  root := L.CheckString(1)
  pattern := L.CheckString(2)
  recursive := false
  absolute := false
  if L.GetTop() >= 3 {
    table := L.CheckTable(3)
    table.ForEach(func(key lua.LValue, value lua.LValue) {
      switch key.String() {
      case "recursive":
        if b, ok := value.(lua.LBool); ok {
          recursive = bool(b)
        }
      case "absolute":
        if b, ok := value.(lua.LBool); ok {
          absolute = bool(b)
        }
      }
    })
  }
      
  matches, err := fs.Glob(fs.Resolve(expand.ExpandVariables(root)), pattern, recursive, absolute)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 2
  }

  table := L.NewTable()
  if matches != nil {
    for _, match := range matches {
      table.Append(lua.LString(match))
    }
  }

  L.Push(table)
  return 1
}

func Basename(L *lua.LState) int {
  path := L.CheckString(1)
  suffix := L.OptBool(2, true)
 
  filename := filepath.Base(path)
  if !suffix {
    filename = strings.TrimSuffix(filename, filepath.Ext(filename))
  }
  
  L.Push(lua.LString(filename))
  return 1
}