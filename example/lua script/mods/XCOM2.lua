local file = require("file")

-- Mods

mods = file.Glob("XCom2-WarOfTheChosen/XComGame/Mods/", "*.XComMod", { recursive = true })
list = { "[Engine.XComModOptions]" }
if mods then
  for i, value in ipairs(mods) do
      table.insert(list, "ActiveMods=\"" .. file.Basename(value, false) .. "\"")
  end
end
config = table.concat(list, "\r\n")
file.Write("XCom2-WarOfTheChosen/XComGame/Config/DefaultModOptions.ini", config)