-- Copyright (c) Anthony Beaumont
-- This source code is licensed under the MIT License
-- found in the LICENSE file in the root directory of this source tree.

local file = require("file")
local regedit = require("regedit")
local process = require("process")

local steamclient = {}

function steamclient.HasGenuineDLL(root)
  root = root or process.cwd
  assert(type(root) == "string" and root ~= "", "Expected a non-empty string!")
  
  local dlls = {"steam_api64.dll", "steam_api.dll"}
  for _, dll in ipairs(dlls) do
    local path = file.Glob(root, dll, { recursive = true })
    if path[1] and path[1] ~= "" then
      local info = file.Info(path[1])
      return info.signed
    end
  end
  return false
end

function steamclient.Backup()
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

function steamclient.Restore(backup)
  assert(type(backup) == "table" and next(backup) ~= nil, "Expected non-empty table!")
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", backup["ActiveUser"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", backup["pid"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", backup["SteamClientDll"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", backup["SteamClientDll64"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", backup["Universe"])
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", backup["RunningAppID"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", backup["SteamExe"])
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", backup["SteamPath"])
end

function findAppID(locations)
  for _, location in ipairs(locations) do
    local paths = file.Glob(location, "steam_appid.txt", { recursive = true })
    for _, path in ipairs(paths) do
      local appid = file.Read(path)
      if appid ~= "" then
        return appid
      end
    end
  end
end

function findDLL(locations, is64)
  local name = is64 and "steamclient64.dll" or "steamclient.dll"
  for _, location in ipairs(locations) do
    local paths = file.Glob(location, name, { recursive = true, absolute = true })
    for _, path in ipairs(paths) do
      local info = file.Info(path)
      if not info.signed then
        return path
      end
    end
  end
end

function steamclient.Load(client)
  client = client or {}
  assert(type(client) == "table", "Expected table!")
  
  client.appid = client.appid or 
                 process.target.env["SteamAppId"] or 
                 process.target.env["SteamGameId"] or
                 process.target.env["SteamOverlayGameId"]
    
  local locations = { process.cwd }
  if process.dir ~= process.cwd then
    locations[#locations + 1] = process.dir
  end

  if not client.appid or client.appid == "" then
    client.appid = findAppID(locations)
  end

  if not client.dll or client.dll == "" then
    client.dll = findDLL(locations, false)
  end

  if not client.dll64 or client.dll64 == "" then
    client.dll64 = findDLL(locations, true)
  end
  
  if type(client.user) ~= "number" or client.user % 1 ~= 0 or client.user == 0 then
    client.user = 1999874061
  end

  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "ActiveUser", tostring(client.user))
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam/ActiveProcess", "pid", tostring(process.pid))
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll", client.dll)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "SteamClientDll64", client.dll64)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam/ActiveProcess", "Universe", "Public")
  regedit.WriteDwordValue("HKCU", "Software/Valve/Steam", "RunningAppID", client.appid)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamExe", process.path)
  regedit.WriteStringValue("HKCU", "Software/Valve/Steam", "SteamPath", process.dir)  
end

return steamclient