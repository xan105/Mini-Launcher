local process = require("process")
local regedit = require("regedit")

-- EAX (DSOAL)

regedit.WriteStringValue("HKCU", "SOFTWARE/Classes/CLSID/{3901CC3F-84B5-4FA4-BA35-AA8172B8A09B}/InprocServer32", "", "dsound.dll")
regedit.WriteStringValue("HKCU", "SOFTWARE/Classes/CLSID/{47D4D946-62E8-11CF-93BC-444553540000}/InprocServer32", "", "dsound.dll")
regedit.WriteStringValue("HKCU", "SOFTWARE/Classes/WOW6432Node/CLSID/{3901CC3F-84B5-4FA4-BA35-AA8172B8A09B}/InprocServer32", "", "dsound.dll")
regedit.WriteStringValue("HKCU", "SOFTWARE/Classes/WOW6432Node/CLSID/{47D4D946-62E8-11CF-93BC-444553540000}/InprocServer32", "", "dsound.dll")

process.On("will-quit", function() 
  regedit.Delete("HKCU", "SOFTWARE/Classes/CLSID/{3901CC3F-84B5-4FA4-BA35-AA8172B8A09B}")
  regedit.Delete("HKCU", "SOFTWARE/Classes/CLSID/{47D4D946-62E8-11CF-93BC-444553540000}")
  regedit.Delete("HKCU", "SOFTWARE/Classes/WOW6432Node/CLSID/{3901CC3F-84B5-4FA4-BA35-AA8172B8A09B}")
  regedit.Delete("HKCU", "SOFTWARE/Classes/WOW6432Node/CLSID/{47D4D946-62E8-11CF-93BC-444553540000}")
end)