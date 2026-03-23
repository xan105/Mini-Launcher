About
=====

<p align="center">
  <img src="https://github.com/xan105/Mini-Launcher/raw/main/banner.png" width="384" height="256"/>
</p>

Mini-Launcher is an application launcher with various gaming-related features. <br />
I built this tool as part of my personal game preservation efforts.

Here is a highlight of its features:

  - DLL Injection / ASI plugin loading
  - Lua Scripting
  - Splash screen (optional)
  - Large Address Aware patching
  - Setting environment variables and arguments with expanding variables
  - Setting compatibility flags or Wine/Proton prefix options
  - Set process priority and affinity (pin to specific cores)

🗑 No bloat: a single binary and a config file.

🐧 This software has an emphasis on being compatible with Linux/Proton.
  
💻 This software is for my own personal use but feel free to use it.

Command Line
============

### `--config string` (launcher.json)

File path to the json configuration file to use. Defaults to `launcher.json`.<br />
Path can be absolute or relative (to the current working dir)

### `--dry-run`

Program will exit before starting the executable.

💡 This flag can come in handy when testing Lua Script.

### `--wait`

Program will wait for the executable to terminate before exiting.

### `--help`

Display a message box with all the command line arguments and a short description.

Config file
===========

By default `launcher.json` (use `--config` to change it).

<details>
<summary>JSON configuration file:</summary>
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
  affinity?: []number,
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

ℹ️ The configuration options are quite extensive, for more details on each option please see [CONFIG_FILE.md](/CONFIG_FILE.md).

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
  },
  "prefix": {
    "overrides": {
      "xinput1_3": "native,builtin"
    }
  }
}
```

See the `./example` directory for more examples.

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
- `%CURRENTDIR%`: Current working dir of the mini-launcher process
- `%BINDIR%`: Parent dir of the mini-launcher process
- `%PROCESS%`: Absolute path of the the mini-launcher process
- `%USERNAME%`
- `%HOSTNAME%`
- `%LANGCODE%`: User's language as ISO 639 language code (ex: `en`, `fr`, `de`)
- `%LANGUAGE%`: User's language in English (ex: `english`, `french`, `german`)
- `%SCREENWIDTH%`, `%XRES%`: Current primary display horizontal resolution (DPI Aware)
- `%SCREENHEIGHT%`, `%YRES%`: Current primary display vertical resolution (DPI Aware)
- `%SCREENREFRESH%`: Current primary display refresh rate

Lua Scripting
=============

Simple scripting engine powered by [yuin/gopher-lua](https://github.com/yuin/gopher-lua) (Lua 5.1).

Originally this feature was for handling CD Key generation in old games. The Lua API has since become quite extensive 😅.

See the `./example` directory for some examples.

## Entry point and permissions

As explained in [CONFIG_FILE.md](/CONFIG_FILE.md), the main script is specified via the `script` option in the config file:

```json
{
  "script": {
    "path": "launcher.lua",
    "fs": false,
    "net": false,
    "reg": false,
    "exec": false,
    "import": false
  }
}
```

There are permissions for different ressources: filesystem operation, network request, etc. You must explicitly grant access to these resources.

The script runs just before the executable. In some cases, you may want to run some code on exit or when the target process runs.
This can be achieved by leveraging the `wait` option or `--wait` cmdline with callback(s) on event(s) from the `process` module:

```lua
local process = require("process")

process.On("will-quit", function()
  print("Bye bye!")
  -- Do something
end)

process.On("did-start", function()
  print("Running!")
  -- Do something
end)
```

## API

Standard libs available are:

  - Package
  - Basic
  - Table
  - String
  - Math
  - Coroutine

Some standard libraries are not enabled by design.<br />
`goto` and `::label::` statements from Lua 5.2 are supported.<br />

By default the VM is mostly sandboxed: you can only _require_ from the available modules.
If you want to _require_ an external lua file you must set the permission `import: true` in the config file.

**API summary**

- `regedit` : read and write from/to the registry.
- `random` : generate random things.
- `file` : file and path manipulation.
- `config` : parse/stringify config files.
- `http` : http request.
- `archive` : decompress archive file.
- `user` : get info about the current user.
- `video` : get info about the current display mode.
- `process` : get info about the current process and the target process to start.
- `shell` : execute shell command. 
- `time` : time conversion.
- `steamid` : Steam-related user identification.
- `steamclient` : utilities to help launching games that require the Steam client (Steamloader).
- `types` : type checking at runtime.
- and some globals for convenience stuff.

ℹ️ For more details on the API please see [LUA_API.md](/LUA_API.md).

⚠️ I may introduce breaking changes between minor version despite my best efforts not to. As the tool mature so will the API.

Build
=====

- Golang v1.26.x
- [go-winres](https://github.com/tc-hib/go-winres) installed in `%PATH%` env var for win32 manifest & cie

Run `build.cmd` on Windows<br/>
Run `build.sh` on Linux<br/>

Output files are located in `./build/${platform}/${config}`

## Github Actions

- `./.github/workflows/go-build-on-linux.yaml` 
- `./.github/workflows/go-build-on-windows.yaml`

Build all targets (debug and release) and create a release.
