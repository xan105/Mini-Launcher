local regedit = require("regedit")
local process = require("process")

-- Steam Loader

local appid = "480"
local steamPath = "%CURRENTDIR%"
local steamExe = "%CURRENTDIR%\\Example.exe"
local clientDLL = "%CURRENTDIR%\\steamclient.dll"
local clientDLL64 = "%CURRENTDIR%\\steamclient64.dll"

-- NB: You also have to set env var with `env:{key:value,...}` in launcher.json
-- Example: 
-- "env": {
--    "SteamAppId": "480",
--    "SteamGameId": "480",
--    "SteamClientLaunch": "1",
--    "SteamEnv": "1",
--    "SteamPath": "%CURRENTDIR%\\Example.exe"
-- }

-- Backup

local backup = {}
backup["ActiveUser"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser")
backup["pid"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid")
backup["SteamClientDll"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll")
backup["SteamClientDll64"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64")
backup["Universe"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe")
backup["RunningAppID"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam", "RunningAppID")
backup["SteamExe"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam", "SteamExe")
backup["SteamPath"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam", "SteamPath")

-- Write new values

regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", "1999874061")
regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", tostring(process.pid))
regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", clientDLL)
regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", clientDLL64)
regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", "Public")
regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", appid)
regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", steamExe)
regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", steamPath)

-- Restore old values

process.On("will-quit", function() -- You may need to use option `wait: true` in launcher.json depending on how the game behave
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", backup["ActiveUser"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", backup["pid"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", backup["SteamClientDll"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", backup["SteamClientDll64"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", backup["Universe"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", backup["RunningAppID"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", backup["SteamExe"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", backup["SteamPath"]) 
end)