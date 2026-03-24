local process = require("process")
local file = require("file")

-- Delete logs
process.On("will-quit", function() 
  file.Remove("Sandfall/Binaries/Win64/ClairObscurFix.log")
end)