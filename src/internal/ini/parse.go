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

func translate(value string, options *ParserOptions) any {
  
  if options.Unquote && len(value) > 2 {
    first := value[0]
    last  := value[len(value)-1]
    if (first == '"' || first == '\'') && first == last {
      value = value[1 : len(value)-1]
    }
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

func Parse(data string, options *ParserOptions) map[string]any {

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
  commentRegex := regexp.MustCompile(`^\s*[;#]`)
  blankLineRegex := regexp.MustCompile(`^\s*$`)

  scanner := bufio.NewScanner(strings.NewReader(data))
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" || commentRegex.MatchString(line) || blankLineRegex.MatchString(line)  { continue }

    if strings.HasPrefix(line, "[") {
      match := sectionRegex.FindStringSubmatch(line)
      if match != nil && len(match) > 1 {
        name := strings.TrimSpace(match[1])
        ignore = slices.Contains(options.Filter, name)
        if !ignore {
          if _, exists := result[name]; !exists {
            result[name] = make(map[string]any)
          }
          section = result[name].(map[string]any)
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