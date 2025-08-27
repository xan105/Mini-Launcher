local user = require("user")
local file = require("file")
local INI = require("config/ini")

local filePath = "steam_emu.ini"

local cfg, err = INI.Parse(file.Read(filePath))
if not err then
  local update = false
  if cfg["Settings"] then
    if cfg["Settings"]["UserName"] == "" or 
       cfg["Settings"]["UserName"] == "CODEX" or 
       cfg["Settings"]["UserName"] == "RUNE" then
     cfg["Settings"]["UserName"] = user.name
     update = true
     end
  end
  
  if cfg["Settings"] and cfg["Settings"]["Language"] == "" then
     cfg["Settings"]["Language"] = user.language
     update = true
  end
  
  if update then
    file.Write(filePath, INI.Stringify(cfg, { 
      whitespace = false, 
      blankLine = true 
    }))
  end
end