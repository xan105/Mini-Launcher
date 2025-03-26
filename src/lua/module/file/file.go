/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package file

import (
  "os"
  "slices"
  "strings"
  "syscall"
  "path/filepath"
  "github.com/yuin/gopher-lua"
  "launcher/internal/fs"
  "launcher/internal/expand"
  "launcher/internal/version"
  "launcher/internal/trust"
  "launcher/lua/type/failure"
)

func Loader(L *lua.LState) int {
  var exports = map[string]lua.LGFunction{
    "Write": Write,
    "Read": Read,
    "Remove": Remove,
    "Info": Info,
    "Glob": Glob,
    "Basename": Basename,
  }
    
  mod := L.SetFuncs(L.NewTable(), exports)
  L.Push(mod)
  return 1
}

func Write(L *lua.LState) int {
  filename := L.CheckString(1)  
  data     := L.CheckString(2)
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
  filename := L.CheckString(1)  
  format   := L.ToString(2)

  if len(format) == 0 {
    format = "utf8"
  } 
  
  data, err := fs.ReadFile(
    fs.Resolve(expand.ExpandVariables(filename)), 
    format,
  )
  if err != nil {
    L.Push(lua.LString(""))
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 2
  }

  L.Push(lua.LString(data))
  return 1
}

func Remove(L *lua.LState) int {
  path := L.CheckString(1)
  
  err := fs.Remove(fs.Resolve(expand.ExpandVariables(path)))
  if err != nil {
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 1
  }
  
  return 0
}

func Info(L *lua.LState) int {
  filename := L.CheckString(1)
  filePath := fs.Resolve(expand.ExpandVariables(filename))
  
  fileInfo, err := os.Stat(filePath)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_FILE_SYSTEM", err.Error()))
    return 2
  }

  info := L.NewTable()
  L.SetField(info, "size", lua.LNumber(fileInfo.Size()))
  
  time := L.NewTable()
  L.SetField(time, "modification", lua.LNumber(fileInfo.ModTime().Unix()))
  if sysInfo, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData); ok {
    L.SetField(time, "creation", lua.LNumber(sysInfo.CreationTime.Nanoseconds() / 1e9))
    L.SetField(time, "access", lua.LNumber(sysInfo.LastAccessTime.Nanoseconds() / 1e9))
  }
  L.SetField(info, "time", time)

  if fileInfo.IsDir() || 
     !slices.Contains([]string{".exe", ".dll"}, strings.ToLower(filepath.Ext(filePath))) {
    L.Push(info)
    return 1
  }

  fileVersionInfo, err := version.FromFile(filePath)
  if err != nil {
    L.Push(info)
    L.Push(failure.LValue(L, "ERR_WIN32_API", err.Error()))
    return 2
  }
  
  version := L.NewTable()
  L.SetField(version, "major", lua.LNumber(fileVersionInfo.Major))
  L.SetField(version, "minor", lua.LNumber(fileVersionInfo.Minor))
  L.SetField(version, "build", lua.LNumber(fileVersionInfo.Build))
  L.SetField(version, "revision", lua.LNumber(fileVersionInfo.Revision))
  L.SetField(info, "version", version)
  
  signed, _ := trust.VerifySignature(filePath)
  L.SetField(info, "signed", lua.LBool(signed))
  
  L.Push(info)
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