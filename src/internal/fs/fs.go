/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package fs

import(
  "os"
  "io"
  "bufio"
  "errors"
  "encoding/json"
  "hash"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/base64"
  "golang.org/x/text/encoding"
  "golang.org/x/text/encoding/charmap"
  "golang.org/x/text/encoding/unicode"
  "golang.org/x/text/transform"
)

func ReadJSON[T any](filepath string) (config T, err error) {

  file, err := os.Open(filepath)
  if err != nil { return }
  defer file.Close()
  
  bytes, err := io.ReadAll(file)
  if err != nil { return }

  err = json.Unmarshal(bytes, &config)
  if err != nil { return }

  return
}

func FileExist(path string) bool {
  target, err := os.Stat(path)
  if err == nil {
    return !target.IsDir()
  }
  if errors.Is(err, os.ErrNotExist) {
    return false
  }
  return false
}

func CheckSum(filePath string, algo string) (result string, err error) {
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

func WriteFile(filename string, data string, format string) error {
  file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
  if err != nil {
    return err
  }
  defer file.Close()

  var enc encoding.Encoding
  switch format {
    case "utf8":
      enc = encoding.Nop //UTF-8 is default
    case "utf8sig":
      enc = unicode.UTF8BOM
    case "utf16le":
      enc = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
    case "windows1252":
      enc = charmap.Windows1252
    default:
      return errors.New("Unsupported encoding: " + format)
  }

  encoder := enc.NewEncoder()
  writer := bufio.NewWriter(transform.NewWriter(file, encoder))
  _, err = writer.WriteString(data)
  if err != nil {
    return err
  }

  return writer.Flush()
}