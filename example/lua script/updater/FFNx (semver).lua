local file = require("file")
local http = require("http")
local JSON = require("config/json")
local time = require("time")
local archive = require("archive")

local DIR = "%CURRENTDIR%/"
local FILEPATH = DIR .. "AF3DN.P"
local REPO = "julianxhokaxhiu/FFNx"
local ASSET = "^FFNx%-Steam%-v[%d%.]+%.zip$"
local URL = "https://api.github.com/repos/" .. REPO .. "/releases/latest"

function parse_version(str)
  local major, minor, build, revision = str:match("^(%d+)%.(%d+)%.(%d+)%.?(%d*)$")
  return {
    major    = tonumber(major) or 0,
    minor    = tonumber(minor) or 0,
    build    = tonumber(build) or 0,
    revision = tonumber(revision) or 0
  }
end

function compare_versions(a, b)
  if a.major > b.major then
    return true
  elseif a.major < b.major then
    return false
  elseif a.minor > b.minor then
    return true
  elseif a.minor < b.minor then
    return false
  elseif a.build > b.build then
    return true
  elseif a.build < b.build then
    return false
  elseif a.revision > b.revision then
    return true
  else
    return false
  end
end

local info, err = file.Info(FILEPATH)
if err then
  error(tostring(err))
end

local res, err = http.Fetch(URL, {
  method = "GET",
  headers = {
    ["Accept"] = "application/vnd.github.v3+json",
    ["User-Agent"] = "Chrome/"
  }
})
if err then 
  error(tostring(err))
end

local github, err = JSON.Parse(res.body)
if err then
  error(tostring(err))
end

local remote = parse_version(github["tag_name"])
if compare_versions(remote, info.version) then

  local target = Array.find(github.assets, function(asset) return asset.name:match(ASSET) end)
  if not target then
    error("No asset found in GitHub response")
  end

  local TMP = "%TEMP%/FFNX"
  
  local path, err = http.Download(target.browser_download_url, TMP)
  if err then
    error(tostring(err))
  end

  local err = archive.Unzip(path, DIR, { "steam_api.dll" })
  file.Remove(TMP)
  if err then
    error(tostring(err))
  end

end