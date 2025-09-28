local user = require("user")
local file = require("file")
local INI = require("config/ini")

local path = "%APPDATA%/GSE Saves/settings"

local content = file.Read(path .. "/configs.user.ini")
local steam = INI.Parse(content, {
  number = false
})
local update = false

if not steam["user::general"] then
  steam["user::general"] = {}
  update = true
end
    
-- User name
    
if not steam["user::general"]["account_name"] or Array.includes({"", "Noob"}, steam["user::general"]["account_name"]) then
  steam["user::general"]["account_name"] = user.name
  update = true
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
    
if not steam["user::general"]["language"] or steam["user::general"]["language"] == "" then
  steam["user::general"]["language"] = user.language
  if not Array.includes(steam_languages, steam["user::general"]["language"]) then
    steam["user::general"]["language"] = "english"
  elseif steam["user::general"]["language"] == "spanish" and user.locale.region ~= "ES" then
    steam["user::general"]["language"] = "latam"
  elseif steam["user::general"]["language"] == "portuguese" and user.locale.region == "BR" then
    steam["user::general"]["language"] = "brazilian"
  elseif steam["user::general"]["language"] == "chinese" then
    if user.locale.region == "CN" or user.locale.region == "SG" then
      steam["user::general"]["language"] = "schinese"
    else
      steam["user::general"]["language"] = "tchinese"
    end
  elseif steam["user::general"]["language"] == "korean" then
    steam["user::general"]["language"] = "koreana"
  end
  update = true
end

-- User Region
    
if not steam["user::general"]["ip_country"] or steam["user::general"]["ip_country"] == "" then
  steam["user::general"]["ip_country"] = user.locale.region
  update = true
end

-- Save

if update then
  local data = INI.Stringify(steam)
  if data ~= "" then
    file.Write(path .. "/configs.user.ini", data)
  end
end

-- Steam Loader

-- NB: You also have to set env var with `env:{key:value,...}` in launcher.json
-- Example: 
-- "env": {
--    "SteamAppId": "480",
--    "SteamGameId": "480",
--    "SteamClientLaunch": "1",
--    "SteamEnv": "1",
--    "SteamPath": "%CURRENTDIR%\\Launcher.exe"
-- }
-- You need to use option `wait: true` in launcher.json if you want to restore modified values on game exit
-- To force inject steamclient/GameOverlayRenderer dll(s) use the `addons` option in launcher.json

local process = require("process")
local steamclient = require("steamclient")
local SteamID = require("SteamID")

if steamclient.HasGenuineDLL() then
  local client = {}
  local sid64 = steam["user::general"]["account_steamid"]
  if sid64 ~= "" then
    client.user = SteamID(sid64).accountid
  end

  local backup = steamclient.Backup()
  steamclient.Load(client)
  process.On("will-quit", function() 
    steamclient.Restore(backup)
  end)
 end