/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//Ported to GoLang from https://github.com/xan105/node-ini (MIT)

package ini

import (
  "strconv"
  "strings"
  "runtime"
)

type StringifyOptions struct {
  Whitespace   bool
  BlankLine    bool
  Quote        bool
  Eol          string
}

func Stringify(data map[string]interface{}, options *StringifyOptions) string {
    var result []string
    
    if options == nil {
      options = &StringifyOptions{
        Whitespace: true,
        BlankLine: false,
        Quote: false,
      }
    }
    
    delimiter := "="
    if options.Whitespace {
      delimiter = " = "
    }
    
    eol := "\n"
    if options.Eol == "\n" || options.Eol == "\r\n" {
      eol = options.Eol
    } else if runtime.GOOS == "windows" {
      eol = "\r\n" 
    }
    
    for key, v := range data {
        switch value := v.(type) {
        case string:
            if options.Quote {
              result = append(result, key + delimiter + "\"" + value + "\"")
            } else {
              result = append(result, key + delimiter + value)
            }
        case bool:
            result = append(result, key + delimiter + strconv.FormatBool(value))
        case float64:
            result = append(result, key + delimiter + strconv.FormatFloat(value, 'f', -1, 64))
        case map[string]interface{}:
            section := "[" + key + "]"
            result = append(result, section)
            result = append(result, Stringify(value, options))
            if options.BlankLine { result = append(result, "") }
      }
    }

    return strings.Join(result, eol)
}