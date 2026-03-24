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

By default, the launcher uses `launcher.json` (you can override it with `--config`).

Most configuration options are straightforward and self-explanatory, so you should be up and running quickly.

_Example with Red Alert 3:_

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

The configuration is quite flexible, covering everything from process launch parameters to DLL injection and compatibility settings.

There are a lot of configuration options available but depending on your use case you may only need a small subset of them.

ℹ️ For a detailed explanation of each option, see [CONFIG_FILE.md](/CONFIG_FILE.md).

You can also check the `./example` directory for some examples.

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

Originally this feature was for handling CD Key generation in old games. It has since become quite extensive 😅.

ℹ️ The Lua API provided to the script is documented in [LUA_API.md](/LUA_API.md)

There are also a few examples in `./example` directory.

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
