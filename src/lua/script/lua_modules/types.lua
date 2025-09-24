local types = {}

local test = {}
test["string"]   = function(v) return type(v) == "string" end
test["str"]      = test["string"]
test["number"]   = function(v) return type(v) == "number" end
test["nbr"]      = test["number"]
test["int"]      = function(v) return type(v) == "number" and v % 1 == 0 end
test["integer"]  = test["int"] 
test["uint"]     = function(v) return test["int"](v) and v >= 0 end
test["boolean"]  = function(v) return type(v) == "boolean" end
test["bool"]     = test["boolean"]
test["table"]    = function(v) return type(v) == "table" end
test["array"]    = test["table"]
test["arr"]      = test["table"]
test["function"] = function(v) return type(v) == "function" end
test["func"]     = test["function"]
test["fn"]       = test["function"]
test["userdata"] = function(v) return type(v) == "userdata" end
test["thread"]   = function(v) return type(v) == "thread" end

local function isArrayOf(tbl, test, length)
  if type(tbl) ~= "table" then return false end
  if length and #tbl ~= length then return false end
  for _, v in ipairs(tbl) do
    if not test(v) then return false end
  end
  return true
end

local function parse(typeString)
  assert(type(typeString) == "string" and #typeString > 0,
         "typeString must be a non-empty string")

  local typeName, lenStr = typeString:match("^(%w+)%s*%[?(%d*)%]?$")
  if not typeName then
    error("Unable to parse type string: " .. typeString)
  end

  local array = typeString:find("%[") ~= nil
  local length = tonumber(lenStr)

  return typeName:lower(), array, length
end

local function translate(name)
  local fn = test[name]
  if not fn then
    error("Unknown type: " .. name)
  end
  return fn
end

function types.is(typeString, value)
  local name, array, length = parse(typeString)
  local test = translate(name)

  if array then
    return isArrayOf(value, test, length)
  else
    return test(value)
  end
end

function types.as(typeString, value)
  return types.is(typeString, value) and value or nil
end

function types.should(typeString, value)
  if not types.is(typeString, value) then
    error("Expected value of type: " .. typeString)
  end
  return value
end

return types