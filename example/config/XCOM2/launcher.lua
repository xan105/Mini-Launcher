local file = require("file")
local user = require("user")

local path = "%APPDATA%/Goldberg SteamEmu Saves/settings/"

-- User name

local account_name = file.Read(path .. "account_name.txt")
if not account_name or account_name == "" or account_name == "Noob" then
  file.Write(path .. "account_name.txt", user.name)
end

-- User Language

local steam_languages = {
  "arabic", "bulgarian", "chinese", "czech",
  "danish", "dutch", "english", "finnish", "french",
  "german", "greek", "hungarian", "italian", "japanese", 
  "korean", "norwegian", "polish", "portuguese",
  "romanian", "russian", "spanish", "swedish",
  "thai", "turkish", "ukrainian", "vietnamese"
}

local language = file.Read(path .. "language.txt")
if not language or language == "" then
  language = user.language
  if not Array.includes(steam_languages, language) then
    language = "english"
  elseif language == "spanish" and user.locale.region ~= "ES" then
    language = "latam"
  elseif language == "portuguese" and user.locale.region == "BR" then
    language = "brazilian"
  elseif language == "chinese" then
    if user.locale.region == "CN" or user.locale.region == "SG" then
      language = "schinese"
    else
      language = "tchinese"
    end
  elseif language == "korean" then
    language = "koreana"
  end
  file.Write(path .. "language.txt", language)
end

-- Mods

mods = file.Glob("XCom2-WarOfTheChosen/XComGame/Mods/", "*.XComMod", { recursive = true })
list = { "[Engine.XComModOptions]" }
if mods then
  for i, value in ipairs(mods) do
      table.insert(list, "ActiveMods=\"" .. file.Basename(value, false) .. "\"")
  end
end
config = table.concat(list, "\r\n")
file.Write("XCom2-WarOfTheChosen/XComGame/Config/DefaultModOptions.ini", config)