local file = require("file")

-- Force Playstation button prompt
local _, err = file.Read("%LOCALAPPDATA%/Sandfall/Saved/Config/Windows/game.ini")
if err then
  local cfg = {}
  cfg["CommonInputPlatformSettings_Windows CommonInputPlatformSettings"] = {
    DefaultGamepadName = "PS5",
    bCanChangeGamepadType = "False"
  }
  file.Write("%LOCALAPPDATA%/Sandfall/Saved/Config/Windows/game.ini", INI.Stringify(cfg))
  file.SetAttributes("%LOCALAPPDATA%/Sandfall/Saved/Config/Windows/game.ini", { readonly = true })
end