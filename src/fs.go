/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "io"
  "encoding/json"
  "errors"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/base64"
  "hash"
)

func readJSON(filepath string) (config Config, err error) {
  file, err := os.Open(filepath)
  if err != nil { return }
  defer file.Close()
  
  bytes, err := io.ReadAll(file)
  if err != nil { return }

  err = json.Unmarshal(bytes, &config)
  if err != nil { return }

  return
}

func fileExist(path string) bool {
  target, err := os.Stat(path)
  if err == nil {
    return !target.IsDir()
  }
  if errors.Is(err, os.ErrNotExist) {
    return false
  }
  return false
}

func checkSum(filePath string, algo string) (result string, err error) {
    file, err := os.Open(filePath)
    if err != nil { return }
    defer file.Close()

    var h hash.Hash
    switch algo {
      case "sha256":
        h = sha256.New()
      case "sha384":
        h = sha512.New384()
      case "sha512":
        h = sha512.New()
      default:
        return "", errors.New("Unsupported hash algorithm: \"" + algo + "\"")
    }

    if _, err = io.Copy(h, file); err != nil { 
      return 
    }

    result = base64.StdEncoding.EncodeToString(h.Sum(nil))
    return
}