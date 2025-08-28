local user = require("user")
local file = require("file")
local INI = require("config/ini")

local filePath = "steam_emu.ini"

local steam_languages = {
  "arabic", "bulgarian", "chinese", "czech",
  "danish", "dutch", "english", "finnish", "french",
  "german", "greek", "hungarian", "italian", "japanese", 
  "korean", "norwegian", "polish", "portuguese",
  "romanian", "russian", "spanish", "swedish",
  "thai", "turkish", "ukrainian", "vietnamese"
}

local content, err = file.Read(filePath)
if err or content == "" then
  return
end

local steam = INI.Parse(content)
local update = false

if steam["Settings"] then
  if not steam["Settings"]["UserName"] or Array.includes({"", "CODEX", "RUNE"}, steam["Settings"]["UserName"]) then
    steam["Settings"]["UserName"] = user.name
    update = true
  end

  if not steam["Settings"]["Language"] or steam["Settings"]["Language"] == "" then
    steam["Settings"]["Language"] = user.language
    if not Array.includes(steam_languages, steam["Settings"]["Language"]) then
      steam["Settings"]["Language"] = "english"
    elseif steam["Settings"]["Language"] == "spanish" and user.locale.region ~= "ES" then
      steam["Settings"]["Language"] = "latam"
    elseif steam["Settings"]["Language"] == "portuguese" and user.locale.region == "BR" then
      steam["Settings"]["Language"] = "brazilian"
    elseif steam["Settings"]["Language"] == "chinese" then
      if user.locale.region == "CN" or user.locale.region == "SG" then
        steam["Settings"]["Language"] = "schinese"
      else
        steam["Settings"]["Language"] = "tchinese"
      end
    elseif steam["Settings"]["Language"] == "korean" then
      steam["Settings"]["Language"] = "koreana"
    end
    update = true
  end
end
  
if update then
  local data = INI.Stringify(steam, { 
    whitespace = false, 
    blankLine = true 
  })
  
  if data ~= "" then
    local err = file.Write(filePath, data)
    if err then
      error(err.message)
    end
  end
end