local regedit = require("regedit")
local process = require("process")
local file = require("file")

local steamclient = {}

function steamclient.hasGenuineDLL()
  local dlls = {"steam_api64.dll", "steam_api.dll"}
  
  for _, dll in ipairs(dlls) do
    local path = file.Glob(process.Cwd(), dll, { recursive = true })
    if path[1] then
      local info = file.Info(path[1])
      return info.signed
    end
  end
  return false
end

function steamclient.backup()
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

function steamclient.restore(backup)
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", backup["ActiveUser"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", backup["pid"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", backup["SteamClientDll"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", backup["SteamClientDll64"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", backup["Universe"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", backup["RunningAppID"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", backup["SteamExe"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", backup["SteamPath"])
end

function steamclient.load(client)

  if not client.appid or client.appid == "" then

    local paths = file.Glob(process.Cwd(), "steam_appid.txt", {
      recursive = true
    })

    for _, path in ipairs(paths) do
      client.appid = file.Read(path)
      if client.appid ~= "" then
        break
      end
    end  
  end

  if not client.dll or client.dll == "" then
    local paths, err = file.Glob(process.Cwd(), "steamclient.dll", {
      recursive = true,
      absolute = true
    })

    for _, path in ipairs(paths) do
      local info = file.Info(path)
      if not info.signed then
        client.dll = path
        break
      end
    end
  end

  if not client.dll64 or client.dll64 == "" then
    local paths = file.Glob(process.Cwd(), "steamclient64.dll", {
      recursive = true,
      absolute = true
    })
    
    for _, path in ipairs(paths) do
      local info = file.Info(path)
      if not info.signed then
        client.dll64 = path
        break
      end
    end
  end

  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", "1999874061")
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", tostring(process.pid))
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", client.dll)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", client.dll64)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", "Public")
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", client.appid)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", process.ExecPath())
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", process.Cwd())  
end

return steamclient