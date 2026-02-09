About
=====

<p align="center">
  <img src="https://github.com/xan105/Mini-Launcher/raw/main/banner.png"/>
  <br/><em>To the moon with Go + Lua ✨</em>
</p>

Mini-Launcher is a CLI application launcher with various gaming-related features.

I built this tool as part of my personal game preservation efforts.

Here is a highlight of its features with some usage example:

  - DLL Injection / ASI plugin loading
      for mods and other compatibility tools when you can't use DLL sideloading.
  - Lua Scripting
      CD Key generation, dealing with mods and config files, [bypassing DRM-like checks](https://xan105.com/blog/scripting-a-steam-loader-using-gopherlua), etc.
  - Splash screen (optional)
  - Large Address Aware patching
      access up to 4GB of virtual memory for old x86 games.
  - Setting environment variables and arguments with expanding variables.
  - Setting compatibility flags or Wine/Proton prefix options.
  - Symlink automation 
      Mimic cloud savegames by relocating savegames to a _"synced folder"_ like OneDrive, Dropbox, etc.

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

JSON configuration file. By default `launcher.json` (use `--config` to change it).

_Example:_

```json
{
  "bin": "Data/ra3_1.12.game",
  "args": "-config \"%CURRENTDIR%\\RA3_%LANGUAGE%_1.12.SkuDef\"",
  "script": {
    "path": "keygen.lua",
    "reg": true
  },
  "addons": [
    { "path": "opencnconline.dll", "required": true }
  ],
  "patch": {
    "laa": true
  }
}
```

<details>
<summary>JSON schema:</summary>
<br />
```ts
{
  bin: string,
  cwd?: string,
  args?: string,
  env?: object,
  priority?: string,
  hide?: bool,
  shell?: bool,
  wait?: bool,
  script?: {
    path: string,
    fs?: bool,
    net?: bool,
    reg?: bool,
    exec?: bool,
    import?: bool
  },
  addons?: []{
    path: string, 
    required?: bool
  },
  suspended?: bool,
  integrity?: []{
    sri: string, 
    path?: string, 
    size?: number,
    signed?: bool
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
  shortcut?: {
    name: string,
    desktop?: bool,
    startmenu?: bool
  },
  compatibility?: {
    version?: string,
    fullscreen?: bool,
    admin?: bool,
    invoker?: bool,
    aware?: bool
  },
  patch?: {
    laa?: bool
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
</details>

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

### `priority?: string` (none)

Set the executable process priority class:

- `IDLE`:          Threads run only when the system is idle.
- `BELOW_NORMAL`:  _Above `IDLE` but below `NORMAL`._
- `NORMAL`:        This is the default.
- `ABOVE_NORMAL`:  _Above `NORMAL` but below `HIGH`._ 
- `HIGH`:          Performs time-critical tasks that must be executed immediately. 
- `REALTIME`:      Highest possible priority, higher than operating system processes performing important tasks.

> [!WARNING] 
> Use extreme care when using **HIGH**, because a high-priority class application can use nearly all available CPU time.
> Even more so with **REALTIME** which can cause issues like disk cache not flushing or unresponsive mouse.

### `hide?: bool` (false)

When enabled, the executable will run without displaying a window, making it invisible to the user.

### `shell?: bool` (false)

When enabled runs inside of a shell (%COMSPEC% ie `cmd.exe`).<br />
💡 Use the `hide` option above to hide the shell.

### `wait?: bool` (false)

When enabled, will wait for the executable to terminate before exiting.

> [!TIP]
> This option can be used in conjunction with the Lua API (see below) to run some code on executable exit! 

```lua
local process = require("process")
process.On("will-quit", function()
  print("Bye bye!")
  //do some cleanup for example
  //like idk deleting logs generated by mods 🙃
end)
```

### `script?: { path: string, fs, net, reg, exec, import?: bool }` (none)

Lua script to be run just before the executable.

Originally this feature was for handling CD Key generation in old games.
The Lua API has since become quite extensive, so please kindly see the **Lua Scripting** section below for more details.

It is quite handy for dealing with old games that requires a bunch of mods to function properly.

See the `./example` directory for some examples.

- `path: string`

  Script file path can be absolute or relative (to the current working dir).<br />
  `%VAR%` are expanded if any (see Expanding Variable for more details).

- `fs, net, reg, exec, import ?: bool` (false)

  These flags act as a simple permissions system: 
  + `fs`: Filesystem operation
  + `net`: Network request
  + `reg`: Windows registry
  + `exec`: Shell command
  + `import`: Load external Lua code (via require, loadfile, loadstring, ...)

  You must explicitly grant access to these resources.

### `addons?: []{ path: string, required?: boolean }` (none)

> DLL Injection / ASI plugin loading

> [!WARNING]
> DLL Injection is generally associated with malware or cheats.
> Using this feature has a high chance of triggering an antivirus false positive.

List of addons to inject to the executable process.<br />
Supported file extension:
  - `.dll` (Dynamic link library)
  - `.asi` (Commonly used in game modding; technically it's a `.dll` file with a different extension).

When `required` is set to `true` and the injection failed: alert the user and kill the process.

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

### `suspended?: bool (false)`

When enabled, the main thread of the executable process is created in a suspended state and does not run until resumed.

The main thread will be automatically resumed after addons are injected to the executable process (see `addons` option above).

> ![TIP]
> Some addons may require the executable to be started in such state in order to function properly.

### `integrity?: []{sri: string, path?: string, size?: number, signed?: bool}` (none)

Check file(s) integrity before starting the executable.

- `sri: string` 
  is a [Subresource Integrity](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity), algo supported are `sha256, sha384, sha512`.

- `path?: string` (executable path) 
  File path, can be absolute or relative (to the current working dir). If no path is specified then the sri targets the executable path.<br />
  `%VAR%` are expanded if any (see Expanding Variable for more details).

- `size?: number`
  optional file size (in bytes), to accelerate sum comparison.
  
- `signed?: bool` (false)
  When true, will check if the file is signed and trusted.
  
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

> [!TIP]
> On Windows the launcher will restart automatically with elevated privileges triggering an UAC prompt for the end user.

### `shortcut?: { name: string, desktop?: bool, startmenu?: bool}` (none)

Creates desktop / start menu shortcut(s) on launch.

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

### `patch?: object` (none)

Patch(es) to apply to the executable.

Does nothing by default; You must explicitly set a patch to `true` or `false`.

- `laa?: bool`:<br/>
  **Large Adress Aware**<br/>
  Allow x86 binary to access up to 4GB of virtual memory (instead of only 2) on x64 Windows.<br/>
  This option set (`true`) or unset (`false`) the LAA flag bit in the PE Header.<br/>
  
NB: Patching is done prior to any integrity check because it directly modifies the executable.

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

💡NB: This is done before running any Lua Script.
If you are dynamically creating/modifying file(s), consider the Lua API: `file.SetAttributes()` (see below).

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

Very simple scripting engine powered by [yuin/gopher-lua](https://github.com/yuin/gopher-lua) (Lua 5.1).<br />
See the `./example` directory for some examples.

Standard libs available are:

  - Package
  - Basic
  - Table
  - String
  - Math
  - Coroutine

Some standard libraries are not enabled by design.<br />
`goto` and `::label::` statements from Lua 5.2 are supported.<br />

The followings modules are exposed to the Lua VM, I might add more later on.

> [!NOTE]
By default the VM is mostly sandboxed: you can only _require_ from the available modules.
If you want to _require_ an external lua file you must set the permission `import: true` in the config file.

## 🌐 Globals

### `sleep(ms: int)`

Suspends the execution of the Lua engine until the time-out interval elapses (interval is in milliseconds).

### `console: SetFuncs`

  + `log(any, ...)`
  + `warn(any, ...)`
  + `error(any, ...)`
  
Convenience methods to print value or array with timestamp and log level. Values are colored depending on their type.

💡 `print()` is an alias to `console.log()`

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

### `Failure(code?: string, message?: string) Failure{ code: string, message: string }`

Failure is a custom type (_userdata_) that represents an "error object" with an associated error code and message.
This provides a structured way to handle error.

- `code?: string` ("ERR_UNKNOWN")
  Error code.
  
- `message?: string` ("An unknown error occurred")
  Error message.
  
💡 `Failure` has a `__tostring` metamethod. If not invoked automatically, you can explicitly call it using `tostring(Failure)`

Example:

```lua
local err = Failure("ERR_NOT_FOUND", "The requested item was not found")
print(err.code)    -- "ERR_NOT_FOUND"
print(err.message) -- "The requested item was not found"
print(err)         -- "[ERR_NOT_FOUND]: The requested item was not found"

local value, err = Foo()
if err and err.code == "ERR_UNKNOWN" then
  error(err.message) -- Raise an error "An unknown error occurred"
  -- or
  error(tostring(err))
end
```

## 📦 Modules

### `📦 Regedit`

This is a module to read and write from/to the registry.

> Requires the `reg` permission.

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
- `Create(root: string, path: string)`
- `Delete(root: string, path: string)`
- `WriteStringValue(root: string, path: string, key: string, value: string)` //REG_SZ
- `WriteExpandStringValue(root: string, path: string, key: string, value: string)` //REG_EXPAND_SZ
- `WriteMultiStringValue(root: string, path: string, key: string, value: []string)` //REG_MULTI_SZ
- `WriteBinaryValue(root: string, path: string, key: string, value: string)` //REG_BINARY
- `WriteDwordValue(root: string, path: string, key: string, value: string)` //REG_DWORD 
- `WriteQwordValue(root: string, path: string, key: string, value: string)` //REG_QWORD
- `DeleteValue(root: string, path: string, key: string)`

✔️ `root` key accepted values are `"HKCR", "HKCU", "HKLM", "HKU" or "HKCC"`.<br />
💡For the default key `@` use `key = ""`

`%VAR%` in `WriteStringValue(..., value)` are expanded if any (see Expanding Variable for more details).

NB: `REG_DWORD` & `REG_QWORD` are represented as string due to floating-point precision limits, if you need to perform arithmetic on them in Lua use `tonumber()`.

### `📦 Random`

This is a module to generate random things.

```lua
local random = require("random")
```

- `AlphaNumString(length: number) string`
- `UserPID() number`
- `SteamID() number`

#### `AlphaNumString(length: number) string`

Generate a random alpha numeric string of specified length.

#### `UserPID() number`

Picks a random PID from the user-owned processes.

#### `SteamID() number`

Generate a random SteamID64.

### `📦 File`

This is a module to read and write text data from/to file.

> Requires the `fs` permission.

```lua
local file = require("file")
```

- `Write(filename: string, data: string, format?: string = "utf8") Failure`
- `Read(filename: string, format?: string = "utf8") string, Failure`
- `Remove(path: string) Failure`
- `Info(filename: string) table, Failure`
- `Glob(root: string, pattern: string, options?: { recursive?: bool = false, absolute?: bool = false }) []string, Failure`
- `Basename(path: string, suffix?: bool = true) string`
- `SetAttributes(filename: string, flags?: { readonly?: bool = false, hidden?: bool = false }) Failure`

Encoding format:

  - `utf8`
  - `utf8sig`
  - `utf16le`
  - `windows1252`

`%VAR%` in `filename` / `root` are expanded if any (see Expanding Variable for more details).

#### `Write(filename: string, data: string, format?: string = "utf8") Failure`

Overwrite text data with specified format encoding (default to utf8).<br /> 
Create target parent dir if doesn't exist.<br />
File is created if doesn't exist.

#### `Read(filename: string, format?: string = "utf8") string, Failure`

Read text data as specified format encoding (default to utf8).

#### `Remove(path: string) Failure`

Delete file or directory and any children it contains at the given path.

#### `Info(filename: string) table, Failure`

Retrieves information for the specified path.<br/>
Time information are represented as Unix epoch time (seconds).<br/>
If the target is a file this will also include the file version information (if any) and whether the file is signed and trusted or not.

```ts
{
  size: number, 
  time: { 
    modification: number, 
    creation?: number, 
    access?: number
  }, 
  version?: { 
    major: number, 
    minor: number, 
    build: number, 
    revision: number 
  },
  signed?: bool
}
```

#### `Glob(root: string, pattern: string, options?: { recursive?: bool = false, absolute?: bool = false }) []string, Failure`

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

#### `SetAttributes(filename: string, flags?: { readonly?: bool = false, hidden?: bool = false }) Failure`

Set file attributes: read only and/or hidden.

### `📦 Config`

This is a module to parse/stringify config files.

```lua
local JSON = require("config/json")
local TOML = require("config/toml")
local INI  = require("config/ini")
local YAML = require("config/yaml")
local XML  = require("config/xml")
```

- `JSON`
  + `Parse(data: string) table | nil, Failure`
  + `Stringify(data: table, pretty?: bool = true) string | nil, Failure`
- `TOML`
  + `Parse(data: string) table | nil, Failure`
  + `Stringify(data: table) string | nil, Failure`
- `INI`
  + `Parse(data: string, options?: table) table`
  + `Stringify(data: table, options?: table) string`
- `YAML`
  + `Parse(data: string) table | nil, Failure`
  + `Stringify(data: table) string | nil, Failure`
- `XML`
  + `Parse(data: string) table | nil, Failure`
  + `Stringify(data: table, pretty?: bool = true) string | nil, Failure`
  
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

### `📦 Http`

This is a module to do http request.

> Requires the `net` permission.

```lua
local http = require("http")
```

- `Fetch(url: string, options?: {method?: string, headers?: table, body?: string }) {status: number , body: string, headers: table} | nil, Failure`
- `Download(url: string, destDir: string) string, Failure`

#### `Fetch(url: string, options?: {method?: string, headers?: table, body?: string }) {status: number , body: string, headers: table} | nil, Failure`

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

#### `Download(url: string, destDir: string) string, Failure`

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

### `📦 Archive`

This is a module to decompress archive file.

> Requires the `fs` permission.

```lua
local archive = require("archive")
```

- `Unzip(filePath: string, destDir: string, excludeList?: []string) Failure`
- `Un7z(filePath: string, destDir: string, excludeList?: []string) Failure`

Extract `.zip` / `7.z` archive to `destDir`. Overwriting existing files.

`%VAR%` are expanded if any (see Expanding Variable for more details).

### `📦 User`

This is a module to get info about the current user.

```lua
local user = require("user")
```

- `name: string` : User name
- `admin: bool` : has elevated rights ?
- `language: string`: User's language in English (ex: `english`, `french`, `german`)
- `locale`: User's language as ISO 639
  + `code: string`: language code (ex: `en`, `fr`, `de`)
  + `region: string`: language region (ex: `US`, `BE`, `DE`)
  
### `📦 Video`

This is a module to get info about the current display mode.

```lua
local video = require("video")
```

- `Current() { width?: number (px), height?: number (px), hz?: number, scale?: number (%)}, Failure`

### `📦 Process`

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

### `📦 Shell`

This is a module to execute shell command. 

> Requires the `exec` permission.

```lua
local shell = require("shell")
```

- `Run(command: string) {stdout?: string, stderr?: string}, Failure`

#### `Run(command: string) {stdout?: string, stderr?: string}, Failure`

Spawns a shell then execute the command within that shell (ComSpec).

### `📦 Time`

This is a module to handle time conversion. 

```lua
local time = require("time")
```

- `Current() number`: Current Unix time
- `HumanizeDuration(seconds: number) string`
- `ToUnix(datetime: string, format?: string = "ISO8601") number, Failure`
- `ToIso8601(datetime: number) string`

NB: `ToUnix()` supported formats are:
- "ISO8601"
- "YYYY-MM-DD"
- "YYYY/MM/DD"
- "YYYY_MM_DD"
- "DD-MM-YYYY"
- "DD/MM/YYYY"
- "MM-DD-YYYY"
- "MM/DD/YYYY"
- "YYYY-MM-DD HH:MM:SS"
- "YYYY/MM/DD HH:MM:SS"

### `📦 SteamID`

This is a module to help working with Steam-related user identification.

```lua
local SteamID = require("SteamID")
```

#### `SteamID(userid: string) SteamID{...}`

`SteamID` is a custom type (_userdata_) that represents a _Steam ID_ with its associated _universe_, _type_, _instance_, and _account ID_.

It is created from a _"Steam2 ID"_ (`STEAM_X:Y:Z`), a _"Steam3 ID"_ (`[U:1:Z]`) or a _"Steam64 ID"_ string.

This provides a structured and easy way to handle conversion:

- `universe: number`

- `type: number`

- `instance: number`

- `accountid: number`

- `:asSteam2() string`

  Returns a "Steam2 ID": `STEAM_X:Y:Z`
  
  eg: `STEAM_1:0:354782281`

- `:asSteam3() string`

  Returns a "Steam3 ID": `[U:1:Z]`
  
  eg: `[U:1:709564562]`

- `:asSteam64() string`
  
  Returns a "Steam64 ID"
  
  eg: `76561198669830290`


Example:

```lua
local SteamID = require("SteamID")
local id = SteamID("76561198669830290")

print(id.accountid)   -- 709564562
print(id:asSteam2())  -- STEAM_1:0:354782281
print(id:asSteam3())  -- [U:1:709564562]
print(id:asSteam64()) -- 76561198669830290
```

### `📦 Steam Client`

This module provides utilities to help launching games that require the Steam client.

These utilities can be used to create what is often referred to as a _"Steam loader"_.

💡If you have no idea what I'm talking about, I invite you to read [my blog post about it](https://xan105.com/blog/scripting-a-steam-loader-using-gopherlua).

> [!IMPORTANT]
> You also have to set Steam-related env. var. with `env:{key:value,...}` in the config file.

```json
{
  "env": {
    "SteamAppId": "480",
    "SteamGameId": "480",
    "SteamClientLaunch": "1",
    "SteamEnv": "1",
    "SteamPath": "%CURRENTDIR%\\Launcher.exe"
  }
}
```

> Requires the `reg` permission.

```lua
local steamclient = require("steamclient")
```

- `HasGenuineDLL(root?: string) bool`

  Recursively search, within the specified root directory, for the presence of genuine (signed) `steam_api(64).dll`.
  If omitted then the launcher's current working directory is used.
  
- `Backup() table`

  Backup the Steam-related registry values.
  
- `Restore(backup: table)`

  Restore previously backed up Steam-related registry values.
  
> [!TIP]
> Use the event _"will-quit"_ from the `process` module to restore the values later on.
> 
> You can also set the option `wait: true` in the config file so the event triggers when the game exits rather than when the launcher terminates.
  
- `Load(client?: { appid?: string, dll:? string, dll64?: string, user?: number })`

  Write the Steam-related values to the registry.<br/>
  You can specify the game's appid, steamclient dlls path and user account id.<br/>
  If omitted they are set automatically by looking for `steam_appid.txt`, `steamclient(64).dll` within the launcher's current working directory.
  
> [!TIP]
> To force inject steamclient/GameOverlayRenderer dll(s) use the `addons` option.
 

**Full example:**

```lua
local process = require("process")
local steamclient = require("steamclient")

if steamclient.HasGenuineDLL() then
  local backup = steamclient.Backup()
  steamclient.Load()
  process.On("will-quit", function() 
    steamclient.Restore(backup)
  end)
 end
```

Config file

```json
{
  "env": {
    "SteamAppId": "480",
    "SteamGameId": "480",
    "SteamClientLaunch": "1",
    "SteamEnv": "1",
    "SteamPath": "%CURRENTDIR%\\Launcher.exe"
  },
  "wait": true,
  "addons": [
    { "path": "steamclient64.dll", "required": true},
    { "path": "GameOverlayRenderer64.dll", "required": true}
  ]
}
```

> [!TIP]
> 🐧 Linux/Proton: you may need to set the env. var. `PROTON_DISABLE_LSTEAMCLIENT=1` _(Linux environment)_ to disable Proton Steam client bridge shenanigans, otherwise it may conflict with the `steamclient(64).dll`.

### `📦 Types`

This is module for type checking at runtime.

```lua
local types = require("types")
```

- `is(typestring: string, value: unknown) bool`
- `as(typestring: string, value: unknown) unknown|nil`
  
  Return the given value when the condition is true otherwise nil.

- `should(typestring: string, value: unknown) unknown`
  
  Return the given value when the condition is true otherwise raise an Error (similar to Lua's `assert()`).

<details><summary>List of supported "type string":</summary>

- string
- str
- number
- nbr
- int
- integer
- uint
- boolean
- bool
- table
- array
- arr
- function
- func
- fn
- userdata
- thread

You can add the suffix `[]` for an array, and add a number for fixed length array.

Example:

```lua
local types = require("types")
print(types.is("string[]", {"hello", "foo", "bar"})) -- true
```

</details>

Build
=====

- Golang v1.25.x
- [go-winres](https://github.com/tc-hib/go-winres) installed in `%PATH%` env var for win32 manifest & cie

Run `build.cmd` on Windows<br/>
Run `build.sh` on Linux<br/>

Output files are located in `./build/${platform}/${config}`

## Github Actions

- `./.github/workflows/go-build-on-linux.yaml` 
- `./.github/workflows/go-build-on-windows.yaml`

Build all targets (debug and release) and create a release.