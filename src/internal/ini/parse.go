/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

//Ported to GoLang from https://github.com/xan105/node-ini (MIT)

package ini

import (
  "bufio"
  "regexp"
  "slices"
  "strings"
  "strconv"
)

type ParserOptions struct {
  Filter    []string
  Global    bool
  Unquote   bool
  Boolean   bool
  Number    bool
}

func unquote(str string) string {
  if len(str) < 2 {
    return str
  }
  
  first := str[0]
  last := str[len(str)-1]

  if (first == '"' || first == '\'') && first == last {
    return str[1 : len(str)-1]
  }

  return str
}

func translate(value string, options *ParserOptions) interface{} {
  
  if options.Unquote {
    value = unquote(value)
  }
  
  if options.Boolean {
    str := strings.ToLower(value)
    if str == "true" || str == "false" {
      if boolean, err := strconv.ParseBool(str); err == nil {
        return boolean
      }
    }
  }
  
  if options.Number {
    // Regex for valid number (same-ish as JSON.parse() rules)
    regex := regexp.MustCompile(`^-?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?$`)
    if regex.MatchString(value) {
      if number, err := strconv.ParseFloat(value, 64); err == nil {
        return number
      }
    } 
  }

  return value
}

func Parse(data string, options *ParserOptions) map[string]interface{} {

  if options == nil {
    options = &ParserOptions{
      Filter: []string{},
      Global: true,
      Unquote: true,
      Boolean: true,
      Number: true,
    }
  }

  result := make(map[string]interface{})
  
  var section map[string]interface{}
  ignore := !options.Global
  
  sectionRegex := regexp.MustCompile(`^\[([^\]]*)\]`)
  
  scanner := bufio.NewScanner(strings.NewReader(data))
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") { continue }

    if strings.HasPrefix(line, "[") {
      match := sectionRegex.FindStringSubmatch(line)
      if match != nil && len(match) > 1 {
        name := strings.TrimSpace(match[1])
        ignore = slices.Contains(options.Filter, name)
        if !ignore {
          if _, exists := result[name]; !exists {
            result[name] = make(map[string]interface{})
          }
          section = result[name].(map[string]interface{})
        }
      } else { 
        ignore = true 
      }
      continue
    }
    
    pos := strings.Index(line, "=")
    if pos < 1 { continue }

    key := strings.TrimSpace(line[:pos])
    value := strings.TrimSpace(line[pos+1:])

    if !ignore { 
      if section == nil {
        result[key] = translate(value, options)
      } else {
        section[key] = translate(value, options)
      }
    }
  }
  
  return result
}