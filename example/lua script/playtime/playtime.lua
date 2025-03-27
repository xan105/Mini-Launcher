local time = require("time")
local process = require("process")

local started = time.Current()
process.On("will-quit", function() -- Use option `wait: true` in launcher.json
  local playtime = time.Current() - started
  console.log("You played for " .. time.HumanizeDuration(playtime))
end)