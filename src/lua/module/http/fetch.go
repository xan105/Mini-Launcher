/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package http

import (
  "io"
  "bytes"
  "net/http"
  "github.com/yuin/gopher-lua"
  "launcher/lua/type/failure"
)

func Fetch(L *lua.LState) int {
  url := L.CheckString(1)

  //Default options
  method := "GET"
  headers := make(map[string]string)
  var body io.Reader

  //Optional 'options' table
  if L.GetTop() >= 2 {
    options := L.CheckTable(2)
    options.ForEach(func(key lua.LValue, value lua.LValue) {
      if keyName, ok := key.(lua.LString); ok {
        switch string(keyName) {
          case "method":
            if methodStr, ok := value.(lua.LString); ok {
              method = string(methodStr)
            }
          case "body":
            if payload, ok := value.(lua.LString); ok {
              body = bytes.NewBuffer([]byte(payload))
            }
          case "headers":
            if headersTable, ok := value.(*lua.LTable); ok {
              headersTable.ForEach(func(hKey lua.LValue, hValue lua.LValue) {
                if hKeyStr, ok := hKey.(lua.LString); ok {
                  if hValueStr, ok := hValue.(lua.LString); ok {
                    headers[string(hKeyStr)] = string(hValueStr)
                  }
                }
              })
            }
        }
      }
    })
  }

  //Create HTTP request
  req, err := http.NewRequest(method, url, body)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_NET_HTTP", err.Error()))
    return 2
  }

  //Set headers
  for key, value := range headers {
    req.Header.Set(key, value)
  }

  //Make the request
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_NET_HTTP", err.Error()))
    return 2
  }
  defer resp.Body.Close()

  //Read response body
  respBody, err := io.ReadAll(resp.Body)
  if err != nil {
    L.Push(lua.LNil)
    L.Push(failure.LValue(L, "ERR_NET_HTTP", err.Error()))
    return 2
  }

  //Create response table
  result := L.NewTable()
  L.SetField(result, "status", lua.LNumber(resp.StatusCode))
  L.SetField(result, "body", lua.LString(string(respBody)))
  headersTable := L.NewTable()

  //Extract headers
  for key, values := range resp.Header {
    if len(values) > 0 {
      L.SetField(headersTable, key, lua.LString(values[0]))
    }
  }
  L.SetField(result, "headers", headersTable)

  L.Push(result)
  return 1
}