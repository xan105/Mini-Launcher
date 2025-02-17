local regedit = require("regedit")
local random = require("random")

local path = "SOFTWARE/Electronic Arts/Electronic Arts/Red Alert 3/ergc"
local current = regedit.QueryStringValue("HKLM", path, "")
if current == "" or current == "%CDKEY%" then
    local key = random.AlphaNumString(20)
    regedit.WriteStringValue("HKLM", path, "", key)
end