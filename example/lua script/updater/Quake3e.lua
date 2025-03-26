local file = require("file")
local http = require("http")
local JSON = require("config/json")
local time = require("time")
local archive = require("archive")

local DIR = "bin/"
local FILEPATH = DIR .. "quake3e-vulkan.x64.exe"
local REPO = "ec-/Quake3e"
local ASSET = "quake3e-windows-msvc-x86_64.zip"
local URL = "https://api.github.com/repos/" .. REPO .. "/releases/latest"

local info = file.Info(FILEPATH)
if not info then
  return 
end

local res, err = http.Fetch(URL, {
  method = "GET",
  headers = {
    ["Accept"] = "application/vnd.github.v3+json",
    ["User-Agent"] = "Chrome/"
  }
})
if err then 
  error(err)
end

local github, err = JSON.Parse(res.body)
if err then
  error(err)
end

local target = Array.find(github.assets, function(asset) return asset.name == ASSET end)
if not target then
  error("No asset found in GitHub response")
end

if(time.ToUnix(target.updated_at) > info.time.modification) then

  local TMP = "%TEMP%/Quake3"
  
  local path, err = http.Download(target.browser_download_url, TMP)
  if err then
    error(err)
  end

  local err = archive.Unzip(path, DIR)
  file.Remove(TMP)
  if err then
    error(err)
  end
  
end