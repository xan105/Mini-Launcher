local regedit = require("regedit")
local random = require("random")

function randAlphaNumString(length)
    local charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    local result = {}

    for i = 1, length do
        local randIndex = math.random(1, #charset)
        result[i] = charset:sub(randIndex, randIndex)
    end

    return table.concat(result)
end

local path = "SOFTWARE/Electronic Arts/Electronic Arts/Red Alert 3/ergc"
local current = regedit.QueryStringValue("HKLM", path, "")
if current == "" or current == "%CDKEY%" then
    -- randAlphaNumString(): Pure Lua implementation
    -- random.AlphaNumString(): Better "randomness"
    -- You can use either one
    local key = random.AlphaNumString(20)
    regedit.WriteStringValue("HKLM", path, "", key)
end