About
=====

Mini-Launcher is an application launcher with the following features:

  - DLL Injection
  - Lua Scripting
  - Splash Screen (optional)
  - File Integrity
  - Expanding Variables
  - Verbatim Arguments
  - Setting Environnement Variables
  - Setting PCA Flags

🐧 This software has an emphasis on being compatible with Linux/Proton.
  
💻 This software is for my own personal use but feel free to use it. 
  
Command Line
============

### `--config string` (launcher.json)

File path to the json configuration file to use. Defaults to `launcher.json`.<br />
Path can be absolute or relative (to the current working dir)

### `--dry-run` (false)

Program will exit before starting the executable.

💡 This flag can come in handy when testing Lua Script.

### `--help` (false)

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
  shell?: bool,
  wait?: bool,
  script?: string,
  addons?: []{
    path: string, 
    required?: bool
  },
  integrity?: []{
    sri: string, 
    path?: string, 
    size?: number
  },
  splash?: {
    show: bool,
    image: []string,
    timeout?: number,
    wait?: string
  },
  symlink?: []{
    path: string,
    dest: string
  },
  compatibility?: {
    version?: string,
    fullscreen?: bool,
    admin?: bool,
    invoker?: bool,
    aware?: bool
  },
  prefix?: {
    winver?: string,
    dpi?: number,
    overrides?: object
  },
  attrib?: []{
    path: string,
    hidden?: boolean,
    readonly?: boolean
  },
  menu?: object
}
```

### `bin: string`

File path to the executable to launch.<br />
Path can be absolute or relative (to the current working dir).

`%VAR%` are expanded if any (see Expanding Variable for more details).

### `cwd?: string` (parent directory)

An option to override the current working dir of the executable to be launched.<br />
This is equivalent to the "Start In" option of a Windows shortcut.

By default the parent directory of the executable is used.<br />
Example: `G:\METAPHOR\METAPHOR.exe` => `G:\METAPHOR\`

### `args?: string` (none)

Optional argument(s) to pass to the executable.<br />
Argument(s) are passed "verbatim" ie: no quoting or escaping is done.

`%VAR%` are expanded if any (see Expanding Variable for more details).

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

`%VAR%` in value are expanded if any (see Expanding Variable for more details).

### `hide?: bool` (false)

When enabled, the executable will run without displaying a window, making it invisible to the user.

### `shell?: bool` (false)

When enabled runs inside of a shell (%COMSPEC% ie `cmd.exe`).<br />
💡Use the `hide` option above to hide the shell.

### `wait?: bool` (false)

When enabled, will wait for the executable to terminate before exiting.

### `script?: string` (none)

File path to a Lua script to be run just before the executable (see **Lua Scripting** below for more details).<br />
Path can be absolute or relative (to the current working dir).

`%VAR%` are expanded if any (see Expanding Variable for more details).

For now this is mainly to handle CD Key generation in old games.<br />
See the `./example` directory for some examples.

### `addons?: []{ path: string, required?: boolean }` (none)

List of addons to inject to the executable process.<br />
When `required` is set to `true` and if the injection failed, alert the user and kill the process.

Path can be absolute or relative (to the current working dir).<br />
`%VAR%` are expanded if any (see Expanding Variable for more details).

Example:
```json
{
  "addons": [
    { "path": "Launcher/opencnconline.dll", "required": true }
  ]
}
```

> [!IMPORTANT]
> This launcher does not support Wow64 injection so make sure the launcher, the executable and the addon are all the same arch (x86 or x64).

### `integrity?: []{sri: string, path?: string, size?: number}` (none)

Check file(s) integrity before starting the executable.

- `sri: string` 
  is a [Subresource Integrity](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity), algo supported are `sha256, sha384, sha512`.

- `path?: string` (executable path) 
  File path, can be absolute or relative (to the current working dir). If no path is specified then the sri targets the executable path.<br />
  `%VAR%` are expanded if any (see Expanding Variable for more details).

- `size?: number`
  optional file size (in bytes), to accelerate sum comparison.
  
### `splash?: { show: bool, image: []string, timeout?: number }` (none)

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
  
- `wait?: string` ("FOREGROUND")
  Which window event to wait for (case insensitive):
  
  + `FOREGROUND`: The foreground window has changed and is visible.
  + `WINDOW`: An hidden window is shown.
  + `CURSOR`: An hidden cursor is shown.
  
### `symlink?: []{path: string, dest: string}` (none)

Creates folder symlink before starting the executable.<br />
Path can be absolute or relative (to the current working dir).<br />
`%VAR%` are expanded if any (see Expanding Variable for more details).

💡 Useful for savegames:

```json
{
  "symlink": [
    { 
      "path": "%DOCUMENTS%/Telltale Games/The Walking Dead",
      "dest": "%SAVEGAME%/The Walking Dead"
    }
  ]
}
```

> [!CAUTION]
> This requires elevated privileges ("Admin rights") or the `SeCreateSymbolicLinkPrivilege` privilege.

### `compatibility?: object` (none)

Set `Program Compatibility Assistant` (PCA) flags, this is equivalent to the `right click > Properties > Compatibility tab` on Windows.

PCA flag(s) are set in `HKCU/Software/Microsoft/Windows NT/CurrentVersion/AppCompatFlags/Layers`.

🐧 This has no effect on application behavior under Wine/Proton

- `version?: string`
  Run the executable in compatibility mode for:
    + `WIN95`
    + `WIN98`
    + `WIN2000`
    + `WINXP`
    + `WINXPSP1`
    + `WINXPSP2`
    + `WINXPSP3`
    + `VISTARTM`
    + `VISTASP1`
    + `VISTASP2`
    + `WIN7RTM`
    + `WIN8RTM`
    
- `fullscreen?: bool`
  Disable fullscreen optimizations.

- `admin?: bool`
  Run the executable as an Administrator.
  
- `invoker?: bool`
  Enforce running the executable as the invoker.
  
- `aware?: bool`
  Override high DPI scaling behavior (Application).

NB: `admin` and `invoker` are mutually exclusive. If both are set `admin` supersede `invoker`.

### `prefix?: {winver?: string, dpi?: number, overrides?: object}` (none)

🐧 Linux Wine/Proton only.

Update current prefix settings.

- `winver?: string`
  Windows version. Accepted values are: 
    + `win11`
    + `win10`
    + `win81` 
    + `win8`
    + `win7`
    + `vista` 
    + `winxp`
  
- `dpi?: number`
  Scale factor: 96 (100%), 120 (125%), 144 (150%), 192 (200%) to 480 (500%).

- `overrides?: object` 
  key-value pairs of dll overrrides. Accepted values:
    + `native,builtin`
    + `builtin,native`
    + `native`
    + `builtin`
    
Example: 

```json
{
  "prefix": {
    "winver": "win10",
    "dpi": 96,
    "overrides": {
      "winmm": "native,builtin"
    }
  }
}
```

### `attrib?: []{path: string, hidden?: boolean = false, readonly?: boolean = false}` (none)

Set file(s) attributes: read only and/or hidden.

Path can be absolute or relative (to the current working dir).<br />
`%VAR%` are expanded if any (see Expanding Variable for more details).

### `menu?: object` (none)

Show a very simple button menu where each key/value pair is a button label and its corresponding override config file.

Example:

```json
{
  "menu": {
    "Mass Effect": "me.json",
    "Mass Effect 2": "me2.json",
    "Mass Effect 3": "me3.json"
  }
}
```

Path can be absolute or relative (to the current working dir).
`%VAR%` in the path are expanded if any (see Expanding Variable for more details).

💡 You can point an entry to the default config file.
If empty, the default is assumed ie: the value of `--config string` which defaults to `launcher.json`

Example:

```json
{
  "binary": "bin/bg3_dx11.exe"
  "menu": {
    "Baldurs Gate 3": "",
    "Baldurs Gate 3 (Vulkan)": "vulkan.json"
  }
}
```

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
- `%HOMEDIR%`, `%USERPROFILE%`
- `%PUBLIC%` 
- `%SYSTEMDIR%`
- `%TEMP%`, `%TMP%`
- `%CURRENTDIR%`: Current working dir of the mini-launcher
- `%BINDIR%`: Dir where the mini-launcher is located at
- `%USERNAME%`
- `%LANGCODE%`: User's language as ISO 639 language code (ex: `en`, `fr`, `de`)
- `%LANGUAGE%`: User's language in English (ex: `english`, `french`, `german`)
- `%SCREENWIDTH%`, `%XRES%`: Current primary display horizontal resolution (DPI Aware)
- `%SCREENHEIGHT%`, `%YRES%`: Current primary display vertical resolution (DPI Aware)
- `%SCREENREFRESH%`: Current primary display refresh rate

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

### `sleep(ms: int)`

Suspends the execution of the Lua engine until the time-out interval elapses (interval is in milliseconds).

### `console: SetFuncs`

  + `log(any, ...)`
  + `warn(any, ...)`
  + `error(any, ...)`
  
Convenience methods to print value or array with timestamp and log level. Values are colored depending on their type.

### `Array: SetFuncs`

  + `find(table, func) any`
  + `some(table, func) bool`
  + `includes(table, any) bool`
  
Convenience methods to search ~~array~~ Lua table.

Example: 

```lua
local arr = {1, 2, 3, 4, 5}

Array.find(arr, function(x) return x > 3 end)
Array.includes(arr, 3)

local arr = {
  {foo = "bar", value = 1},
  {foo = "baz", value = 2}
}

Array.find(arr, function(x) return x.foo == "bar" end)
Array.some(arr, function(x) return x.foo == "baz" end)
```

### 📦 Regedit

This is a module to read and write from/to the registry.

```lua
local regedit = require("regedit")
```

- `KeyExists(root: string, path: string) bool`
- `ListAllSubkeys(root: string, path: string) []string`
- `ListAllValues(root: string, path: string) []string`
- `QueryValueType(root: string, path: string, key: string) string`
- `QueryStringValue(root: string, path: string, key: string) string` //REG_SZ & REG_EXPAND_SZ
- `QueryMultiStringValue(root: string, path: string, key: string) []string` //REG_MULTI_SZ
- `QueryBinaryValue(root: string, path: string, key: string) string` //REG_BINARY
- `QueryIntegerValue(root: string, path: string, key: string) string` //REG_DWORD & REG_QWORD
- `WriteKey(root: string, path: string)`
- `DeleteKey(root: string, path: string)`
- `WriteStringValue(root: string, path: string, key: string, value: string)` //REG_SZ
- `WriteExpandStringValue(root: string, path: string, key: string, value: string)` //REG_EXPAND_SZ
- `WriteMultiStringValue(root: string, path: string, key: string, value: []string)` //REG_MULTI_SZ
- `WriteBinaryValue(root: string, path: string, key: string, value: string)` //REG_BINARY
- `WriteDwordValue(root: string, path: string, key: string, value: string)` //REG_DWORD 
- `WriteQwordValue(root: string, path: string, key: string, value: string)` //REG_QWORD
- `DeleteKeyValue(root: string, path: string, key: string)`

✔️ `root` key accepted values are `"HKCR", "HKCU", "HKLM", "HKU" or "HKCC"`.<br />
💡For the default key `@` use `key = ""`

`%VAR%` in `WriteStringValue(..., value)` are expanded if any (see Expanding Variable for more details).

NB: `REG_DWORD` & `REG_QWORD` are represented as string due to floating-point precision limits, if you need to perform arithmetic on them in Lua use `tonumber()`.

### 📦 Random

This is a module to generate random things.

```lua
local random = require("random")
```

- `AlphaNumString(length: number) string`
- `UserPID() number`

#### `AlphaNumString(length: number) string`

Generate a random alpha numeric string of specified length.

#### `UserPID() number`

Picks a random PID from the user-owned processes.

### 📦 File

This is a module to read and write text data from/to file.

```lua
local file = require("file")
```

- `Write(filename: string, data: string, format?: string = "utf8") error`
- `Read(filename: string, format?: string = "utf8") string, error`
- `Remove(path: string) error`
- `Version(filename: string) {Major, Minor, Build, Revision: number}, error`
- `Glob(root: string, pattern: string, options?: { recursive?: false, absolute?: false }) []string, err`
- `Basename(path: string, suffix?: bool = true) string`

Encoding format:

  - `utf8`
  - `utf8sig`
  - `utf16le`
  - `windows1252`

`%VAR%` in `filename` / `root` are expanded if any (see Expanding Variable for more details).

#### `Write(filename: string, data: string, format?: string = "utf8") error`

Overwrite text data with specified format encoding (default to utf8).<br /> 
Create target parent dir if doesn't exist.<br />
File is created if doesn't exist.

#### `Read(filename: string, format?: string = "utf8") string, error`

Read text data as specified format encoding (default to utf8).

#### `Remove(path: string) error`

Delete file or directory and any children it contains at the given path.

#### `Version(filename: string) {Major, Minor, Build, Revision: number}, error`

Retrieves version information for the specified file. Binary only (exe/dll).

#### `Glob(root: string, pattern: string, options?: { recursive?: false, absolute?: false }) []string, err`

Returns the names of all files matching pattern. The pattern syntax is the same as in Go [path/filepath Match](https://pkg.go.dev/path/filepath#Match).
With the addition that, to return only directories the pattern should end with `/`.

#### `Basename(path: string, suffix?: bool = true) string`

Returns the last element of path. When `suffix` is `false` the file extension is removed.

Example:

```lua
file.Basename("/foo/bar/quux.html");
-- Returns: "quux.html"

file.Basename("/foo/bar/quux.html", false);
-- Returns: "quux" 
```

### 📦 Config

This is a module to parse/stringify config files.

```lua
local JSON = require("config/json")
local TOML = require("config/toml")
local INI  = require("config/ini")
local YAML = require("config/yaml")
local XML  = require("config/xml")
```

- `JSON`
  + `Parse(data: string) table, err`
  + `Stringify(data: table, pretty?: bool = true) string, err`
- `TOML`
  + `Parse(data: string) table, err`
  + `Stringify(data: table) string, err`
- `INI`
  + `Parse(data: string, options?: table) table`
  + `Stringify(data: table, options?: table) string`
- `YAML`
  + `Parse(data: string) table, err`
  + `Stringify(data: table) string, err`
- `XML`
  + `Parse(data: string) table, err`
  + `Stringify(data: table, pretty?: bool = true) string, err`
  
⚠️ Due to GoLang using hashmap the key order is not guaranteed !

#### `INI`

Parse options: 

- `filter?: []string (none)` Section filter
- `global?: bool (true)` Include global section
- `unquote?: bool (true)` Unquote string (starting/ending with `"` or `'`)
- `boolean?: bool (true)` String to boolean type conversion
- `number?: bool (true)` String to number type conversion (same rules as JavaScript's JSON.parse())

Stringify options:

- `whitespace?: bool (true)` add space between delimiter `=`
- `blankLine?: bool (false)` add empty line between sections
- `quote?: bool (false)` quote string with `"`
- `eol?: string (system)` Either `\n` or `\r\n`

### 📦 Http

This is a module to do http request.

```lua
local http = require("http")
```

- `Fetch(url: string, options?: {method?: string, headers?: table, body?: string }) {status, body, headers}, err`
- `Download(url: string, destDir: string) string, err`

#### `Fetch(url: string, options?: {method?: string, headers?: table, body?: string }) {status, body, headers}, err`

A [Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) like API.

Example:

```lua
local http = require("http")
local JSON = require("config/json")

local repo = "xan105/Mini-Launcher"
local url = "https://api.github.com/repos/" .. repo .. "/releases/latest"

local res, err = http.Fetch(url, {
  method = "GET",
  headers = {
    ["Accept"] = "application/vnd.github.v3+json",
    ["User-Agent"] = "Chrome/"
  }
})
if err then
  console.error(err)
end

local github, err = JSON.Parse(res.body)
if err then
  console.error(err)
end

local latestRelease = github["tag_name"]
```

#### `Download(url: string, destDir: string) string, err`

Download a file. Filename is determined by the `Content-Disposition` header.
Create target parent dir if doesn't exist.<br />
Overwrite existing file.<br />
Return the downloaded file path.

`%VAR%` in `destDir` are expanded if any (see Expanding Variable for more details).

Example:

```lua
local http = require("http")

local filepath, err = http.Download("http://.../foo.bar", "%DOWNLOAD%")
if err then
  console.error(err)
else
  console.log("downloaded: " .. filepath)
end
```

### 📦 Archive

This is a module to decompress archive file.

```lua
local archive = require("archive")
```

- `Unzip(filePath: string, destDir: string) err`

#### `Unzip(filePath: string, destDir: string) err`

Extract `.zip` archive to `destDir`. Overwriting existing files.

`%VAR%` are expanded if any (see Expanding Variable for more details).

### 📦 User

This is a module to get info about the current user.

```lua
local user = require("user")
```

- `name: string` : User name
- `language: string`: User's language in English (ex: `english`, `french`, `german`)
- `locale`: User's language as ISO 639
  + `code: string`: language code (ex: `en`, `fr`, `de`)
  + `region: string`: language region (ex: `US`, `BE`, `DE`)
  
### 📦 Video

This is a module to get info about the current display mode.

```lua
local video = require("video")
```

- `Current() { width: number (px), height: number (px), hz: number, scale: number (%)}, err`

### 📦 Process

This is a module to get info about the current Mini-Launcher process.

```lua
local process = require("process")
```

- `platform: string` : operating system target (GOOS)
- `arch: string` : architecture target (GOARCH)
- `pid: number` : process id
- `Cwd() string` : process current working dir
- `ExecPath() string` : process absolute pathname
- `On(event: string, callback: function)` : register callback function to be run for specified event

**Events**

- `will-quit` : Fired when process is about to terminate.

Build
=====

- Golang v1.24.0
- [go-winres](https://github.com/tc-hib/go-winres) installed in `%PATH%` env var for win32 manifest & cie

Run `build.cmd` on Windows<br/>
Run `build.sh` on Linux<br/>

Output files are located in `./build/${platform}/${config}`
