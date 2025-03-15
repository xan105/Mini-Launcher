local file = require("file")
local user = require("user")
local INI = require("config/ini")

local filePath = "steam_emu.ini"

local cfg, err = INI.Parse(file.Read(filePath))
if not err then
  if cfg["Settings"] and cfg["Settings"]["UserName"] == "" then
     cfg["Settings"]["UserName"] = user.name
  end
  if cfg["Settings"] and cfg["Settings"]["Language"] == "" then
     cfg["Settings"]["Language"] = user.language
  end

  file.Write(filePath, INI.Stringify(cfg, { 
    whitespace = false, 
    blankLine = true 
  }))
end