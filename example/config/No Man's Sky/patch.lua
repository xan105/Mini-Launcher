local memory = require("memory")

local patch = { -- Remove the "Mod Enabled" Warning screen on startup
  pattern = "48 8B 01 48 85 C0 74 08 0F B6 80 5A 46 00 00 C3", 
  offset = 0x06, 
  value = "EB",
}

function apply(patch)
  local address, err = memory.Find(patch.pattern)
  if err then
    error(err.message)
  end
  local success, err = memory.Write(address + patch.offset, patch.value)
  if err then
    error(err.message)
  end
  return success
end

if apply(patch) then
  console.log("Applied patch!")
end