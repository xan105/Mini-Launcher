local regedit = require("regedit")
local process = require("process")

local steam = {}

function steam.backup()
  local backup = {}
  backup["ActiveUser"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser")
  backup["pid"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid")
  backup["SteamClientDll"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll")
  backup["SteamClientDll64"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64")
  backup["Universe"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe")
  backup["RunningAppID"] = regedit.QueryIntegerValue("HKCU", "Software/Valve/Steam", "RunningAppID")
  backup["SteamExe"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam", "SteamExe")
  backup["SteamPath"] = regedit.QueryStringValue("HKCU", "Software/Valve/Steam", "SteamPath")
  return backup
end

function steam.restore(backup)
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", backup["ActiveUser"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", backup["pid"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", backup["SteamClientDll"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", backup["SteamClientDll64"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", backup["Universe"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", backup["RunningAppID"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", backup["SteamExe"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", backup["SteamPath"])
end

function steam.load(client)
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", "1999874061")
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", tostring(process.pid))
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", client.dll)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", client.dll64)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", "Public")
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", client.appid)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", process.ExecPath())
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", process.Cwd())  
end

return steam