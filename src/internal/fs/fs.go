/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package fs

import(
  "os"
  "io"
  "hash"
  "bufio"
  "errors"
  "runtime"
  "path/filepath"
  "encoding/json"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/base64"
  "golang.org/x/text/encoding"
  "golang.org/x/text/transform"
  "golang.org/x/text/encoding/charmap"
  "golang.org/x/text/encoding/unicode"
)

func Resolve(filePath string) string {
  path := filepath.FromSlash(filePath)
  if !filepath.IsAbs(path) {
    fullPath, err := filepath.Abs(path) //Uses GetFullPathNameW() on Windows
    if err != nil && runtime.GOOS == "windows" {
      cwd, _ := os.Getwd()
      path = filepath.Join(cwd, path)
    } else {
      path = fullPath
    }
  }
  return path
}

func ReadJSON[T any](filePath string) (config T, err error) {

  file, err := os.Open(filePath)
  if err != nil { return }
  defer file.Close()
  
  bytes, err := io.ReadAll(file)
  if err != nil { return }

  err = json.Unmarshal(bytes, &config)
  if err != nil { return }

  return
}

func FileExist(filePath string) bool {
  target, err := os.Stat(filePath)
  if err == nil {
    return !target.IsDir()
  }
  if errors.Is(err, os.ErrNotExist) {
    return false
  }
  return false
}

func CheckSum(filePath string, algo string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil { return "", err }
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
      return "", err
    }

    return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func WriteFile(filePath string, data string, format string) error {
  
  dir := filepath.Dir(filePath)
  if err := os.MkdirAll(dir, 0755); err != nil {
    return err
  }
  
  file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
      return errors.New("Unsupported encoding: : \"" + format + "\"")
  }

  encoder := enc.NewEncoder()
  writer := bufio.NewWriter(transform.NewWriter(file, encoder))
  _, err = writer.WriteString(data)
  if err != nil {
    return err
  }

  return writer.Flush()
}

func ReadFile(filePath string, format string) (string, error) {

  file, err := os.Open(filePath)
  if err != nil {
    return "", err
  }
  defer file.Close()

  var enc encoding.Encoding
  switch format {
  case "utf8":
    enc = encoding.Nop //UTF-8 is default
  case "utf8sig":
    enc = unicode.UTF8BOM
  case "utf16le":
    enc = unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
  case "windows1252":
    enc = charmap.Windows1252
  default:
    return "", errors.New("Unsupported encoding: \"" + format + "\"")
  }

  decoder := enc.NewDecoder()
  reader := bufio.NewReader(transform.NewReader(file, decoder))
  data, err := io.ReadAll(reader)
  if err != nil {
    return "", err
  }

  return string(data), nil
}

func CreateFolderSymlink(origin string, destination string) error {

  target, err := os.Lstat(origin)
  if err != nil {
    if !errors.Is(err, os.ErrNotExist) {
      return err
    }
  }
  
  if err == nil {
    if target.Mode()&os.ModeSymlink != 0 { //Already a symlink
      return nil                           //Nothing to do
    }
  
    if !target.IsDir() { //Target is a file
      return errors.New("Symlink target is a file, aborting !")
    }
    
    entries, err := os.ReadDir(origin)
    if err != nil {
      return err
    }
    if len(entries) != 0 { 
      return errors.New("Symlink target is a non-empty dir, aborting !")
    }
    
    //Empty so safe to delete
    err = os.Remove(origin)
    if err != nil {
      return err
    }
  }

  err = os.MkdirAll(destination, 0755)
  if err != nil {
    return err
  }

  err = os.Symlink(destination, origin)
  if err != nil {
    return err
  }

  return nil
}