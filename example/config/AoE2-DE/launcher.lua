local file = require("file")
local INI = require("config/ini")
local user = require("user")
local regedit = require("regedit")
local process = require("process")

-- Settings path

local path = "%APPDATA%/GSE Saves/settings"
local local_user_file = file.Read("steam_settings/configs.user.ini")

if local_user_file then
  user_cfg = INI.Parse(local_user_file)
  if user_cfg and user_cfg["user::saves"] then
    local_save_path = user_cfg["user::saves"]["local_save_path"]
    if not local_save_path or local_save_path == "" then
      saves_folder_name = user_cfg["user::saves"]["saves_folder_name"]
      if saves_folder_name and saves_folder_name ~= "" then
        path = "%APPDATA%/" .. saves_folder_name .. "/settings"
      end
    end
  end
end

local global_user_cfg = {}
local global_user_file = file.Read(path .. "/configs.user.ini")
if global_user_file then
  global_user_cfg = INI.Parse(global_user_file)
end

if not global_user_cfg["user::general"] then
  global_user_cfg["user::general"] = {}
end
    
-- User name
    
local account_name = global_user_cfg["user::general"]["account_name"]
if not account_name or account_name == "" or account_name == "Noob" then
  account_name = user.name
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
    
local language = global_user_cfg["user::general"]["language"]
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
end

-- User Region
    
local region = global_user_cfg["user::general"]["ip_country"]
if not region or region == "" then
  region = user.locale.region
end

-- Save
    
global_user_cfg["user::general"]["account_name"] = account_name
global_user_cfg["user::general"]["language"] = language
global_user_cfg["user::general"]["ip_country"] = region

file.Write(path .. "/configs.user.ini", INI.Stringify(global_user_cfg))

-- Steam Loader

local backup = {}
backup["ActiveUser"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser")
backup["pid"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid")
backup["SteamClientDll"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll")
backup["SteamClientDll64"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64")
backup["Universe"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe")
backup["RunningAppID"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam", "RunningAppID")
backup["SteamExe"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam", "SteamExe")
backup["SteamPath"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam", "SteamPath")

regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", "1999874061")
regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", tostring(process.pid))
regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", "%CURRENTDIR%\\steamclient.dll")
regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", "%CURRENTDIR%\\steamclient64.dll")
regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", "Public")
regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", "813780")
regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", "%CURRENTDIR%\\AoE2DE_s.exe")
regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", "%CURRENTDIR%")

process.On("will-quit", function() 
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", backup["ActiveUser"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", backup["pid"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", backup["SteamClientDll"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", backup["SteamClientDll64"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", backup["Universe"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", backup["RunningAppID"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", backup["SteamExe"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", backup["SteamPath"]) 
end)