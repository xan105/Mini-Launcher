About
=====

Mini-Launcher is an application launcher with the following features:

  - DLL Injection
  - Lua Scripting
  - Optional Splash screen
  - File integrity check
  - Expanding Variable
  - Verbatim Arguments
  - Setting Environnement Variables

🐧 This software has an emphasis on being compatible with Linux/Proton.
  
💻 This software is for my own personal use but feel free to use it. 
  
Command Line
============

### `-config string` (launcher.json)

File path to the json configuration file to use. Defaults to `launcher.json`.<br />
Path can be absolute or relative (to the current working dir)

### `-help`

Display a message box with all the command line arguments and a short description.

Config file
===========

```ts
{
  bin: string,
  cwd?: string,
  args?: string,
  env?: object,
  hide?: bool,
  script?: string,
  addons?: []{
    path: string, 
    required?: bool
  },
  integrity?: []{
    sri: string, 
    file?: string, 
    size?: number
  },
  splash?: {
    show: bool,
    image: []string,
    timeout?: number
  }
}
```

### `bin: string`

File path to the executable to launch.<br />
Path can be absolute or relative (to the current working dir).

`%VAR%` are expanded if any (see Expanding Variable for more details)

### `cwd?: string` (parent directory)

An option to override the current working dir of the executable to be launched.<br />
This is equivalent to the "Start In" option of a Windows shortcut.

By default the parent directory of the executable is used.<br />
Example: `G:\METAPHOR\METAPHOR.exe` => `G:\METAPHOR\`

### `args?: string` (none)

Optional argument(s) to pass to the executable.<br />
Argument(s) are passed "verbatim" ie: no quoting or escaping is done.

`%VAR%` are expanded if any (see Expanding Variable for more details)

### `env?: object` (none)

Add additional environment key-value pairs to the executable process.

Example:

```json
{
  "env": {
    "GAMEPAD_LED": "BATTERYLVL"
  }
}
```

`%VAR%` in value are expanded if any (see Expanding Variable for more details)

### `hide?: bool` (false)

When enabled, the executable will run without displaying a window, making it invisible to the user.<br /> 
This is useful for background tasks or command-line utilities that do not require user interaction.

### `script?: string` (none)

File path to a Lua script to be run just before the executable (see Lua Scripting for more details).<br />
Path can be absolute or relative (to the current working dir).

For now this is mainly to handle CD Key generation in old games.

### `addons?: []{ path: string, required?: boolean }` (none)

List of addons to inject to the executable process.<br />
When `required` is set to `true` and if the injection failed, alert the user and kill the process.

Path can be absolute or relative (to the current working dir).

Example:
```json
{
  "addons": [
    { "path": "Launcher/opencnconline.dll", "required": true }
  ]
}
```

### `integrity?: []{sri: string, file?: string, size?: number}`

Check file(s) integrity before starting the executable.

- `sri: string` 
  is a [Subresource Integrity](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity), algo supported are `sha256, sha384, sha512`.

- `file?: string` (executable path) 
  File path, can be absolute or relative (to the current working dir). If no path is specified then the sri targets the executable path.<br />
  `%VAR%` are expanded if any (see Expanding Variable for more details).

- `size?: number`
  optional file size (in bytes), to accelerate sum comparison.
  
### `splash?: { show: bool, image: []string, timeout?: number }`

Display a splash screen until the executable process change the cursor or display a window. The splash screen should be a BMP file. 

- `show: bool` (false)
  Wether to display the splash screen or not.

- `image: []string` 
  Splash screen filepath. When more than one, a splash screen is selected at random.<br />
  File path can be absolute or relative (to the current working dir).<br />
  `%VAR%` are expanded if any (see Expanding Variable for more details).
  
- `timeout?: number` (10 sec)
  Failsafe timeout in seconds.<br />
  There was no event dispatched under Linux/Proton on Wayland in my limited testing.

Expanding Variable
==================

List of variables that will get expanded:

- `%APPDATA%`
- `%LOCALAPPDATA%`
- `%PROGRAMDATA%`
- `%DESKTOP%`
- `%DOCUMENTS%`
- `%MUSIC%`
- `%PICTURES%`
- `%VIDEOS%`
- `%DOWNLOAD%`
- `%SAVEGAME%`
- `%HOMEDIR%`
- `%USERPROFILE%`
- `%PUBLIC%` 
- `%SYSTEMDIR%`
- `%TEMP%`
- `%TMP%`
- `%CURRENTDIR%`: Current working dir of the mini-launcher
- `%BINDIR%`: Dir where the mini-launcher is located at

Lua Scripting
=============

Very simple scripting engine powered by [yuin/gopher-lua](https://github.com/yuin/gopher-lua).

- Lua 5.1
- Libraries:
  + Package
  + Basic
  + Table
  + String
  + Math
  
Some standard libraries are not enabled by design.<br />
The followings modules are exposed to the Lua VM, I might add more later on.

### 🌐 Globals

### `sleep(ms int)

Suspends the execution of the Lua engine until the time-out interval elapses (interval is in milliseconds).

### 📦 Regedit

This is a module to read and write from/to the registry.

```lua
local regedit = require("regedit")
```

- `QueryStringValue(root: string, path: string, key: string)  string`
- `WriteStringValue(root: string, path: string, key: string, value: string)`

✔️ `root` key accepted values are `"HKCR", "HKCU", "HKLM", "HKU" or "HKCC"`.

#### `QueryStringValue(root: string, path: string, key: string) string`

> REG_SZ & REG_EXPAND_SZ

Return string value of given path/key.

💡For the default key `@` use `key = ""`

#### `WriteStringValue(root: string, path: string, key: string, value: string)`

> REG_SZ & REG_EXPAND_SZ

Write string value in given path/key (subkeys are created if necessary).

💡For the default key `@` use `key = ""`

### 📦 Random

This is a module to generate random things.

```lua
local random = require("random")
```

- `AlphaNumString(length: int) string`

#### `AlphaNumString(length: int) string`

Generate a random alpha numeric string of specified length.

### 📦 File

This is a module to read and write text data from/to file.

```lua
local file = require("file")
```

- `Write(filename: string, data: string, format?: string = "utf8")`

#### `Write(filename: string, data: string, format?: string = "utf8")`

Overwrite text data with specified format encoding (default to utf8). File is created if needed.

Encoding format:

  - `utf8`
  - `utf8sig`
  - `utf16le`
  - `windows1252`

Build
=====

- Golang v1.24.0
- [go-winres](https://github.com/tc-hib/go-winres) installed in `%PATH%` env var for win32 manifest & cie

Run `build.cmd` on Windows<br/>
Run `build.sh` on Linux<br/>

Output files are located in `./build/${platform}/${config}`
