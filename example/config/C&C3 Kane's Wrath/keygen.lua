local regedit = require("regedit")

function randAlphaNumString(length)
    local charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    local result = {}

    for i = 1, length do
        local randIndex = math.random(1, #charset)
        result[i] = charset:sub(randIndex, randIndex)
    end

    return table.concat(result)
end

local path = "SOFTWARE/Electronic Arts/Electronic Arts/Command and Conquer 3 Kanes Wrath/ergc"
local current = regedit.QueryStringValue("HKLM", path , "")
if current == "" or current == "%CDKEY%" then
    local key = randAlphaNumString(20)
    regedit.WriteStringValue("HKLM", path, "", key)
end