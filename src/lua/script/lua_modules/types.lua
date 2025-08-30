-- Copyright (c) Anthony Beaumont
-- This source code is licensed under the MIT License
-- found in the LICENSE file in the root directory of this source tree.

local types = {}

function types.isString(value)
  return type(value) == "string"
end

function types.isNumber(value)
  return type(value) == "number"
end

function types.isSafeInteger(value)
    local MAX_SAFE_INTEGER = 2^53 - 1
    local MIN_SAFE_INTEGER = -MAX_SAFE_INTEGER

    if not types.isNumber(value) then
        return false
    end

    return value >= MIN_SAFE_INTEGER and value <= MAX_SAFE_INTEGER and value % 1 == 0
end

function types.isBoolean(value)
  return type(value) == "boolean"
end

function types.isNil(value)
  return type(value) == "nil"
end

function types.isTable(value)
  return type(value) == "table"
end

function types.isFunction(value)
  return type(value) == "function"
end

function types.isUserData(value)
  return type(value) == "userdata"
end

function types.isThread(value)
  return type(value) == "thread"
end

function types.isEmpty(value)
  if value == nil then
    return true
  end

  if types.isString(value) and value == "" then
    return true
  end

  if types.isTable(value) and next(value) == nil then
    return true
  end

  return false
end

return types