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

This is a module to help with file and path manipulation.

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
- `Dirname(path: string) string`
- `Extname(path: string) string`
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

#### `Dirname(path: string) string`

Returns the directory name of a path.

#### `Extname(path: string) string`

Returns the extension of a path.

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

This is a module to get info about the current process (Mini-Launcher) and the target process to start.

```lua
local process = require("process")
```

- `platform: string` : operating system target (GOOS)
- `arch: string` : architecture target (GOARCH)
- `pid: number` : process id
- `wine: bool`: whether process is running under wine/proton or not
- `path: string`: process absolute pathname
- `bin: string`: process file name
- `dir: string`: process parent dir
- `cwd: string`: process current working dir
- `args: string[]`: process arguments
- `env: { key: string, ... }` : process env. var. as key:value pairs
- `target.path: string`: target process absolute pathname
- `target.bin: string`: target process file name
- `target.dir: string`: target process parent dir
- `target.cwd: string`: target process current working dir
- `target.argv: string[]`: target process verbatim arguments
- `target.env: { key: string, ... }` : target process env. var. as key:value pairs
- `On(event: string, callback: function)` : register callback function to be run for specified event

**Events**

- `will-quit()` : Fired when process is about to terminate.
- `did-start(event: { pid: number })` : Fired when the target executable start sequence is over (spawning, addons, splash screen).

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
    "SteamPath": "%PROCESS%"
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
  You can specify the game's appid (defaults to the `SteamAppId` env var), steamclient dlls path and user account id.<br/>
  If omitted they are set automatically by looking for `steam_appid.txt`, `steamclient(64).dll` within the launcher's current working directory and parent directory.
  
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
    "SteamPath": "%PROCESS%"
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