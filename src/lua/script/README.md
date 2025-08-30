Lua modules **written in Lua** located in the `lua_modules/` folder are embedded into the Go executable using [Go's embed](https://pkg.go.dev/embed) feature.

- The filename _(without `.lua`)_ is used as the module name:

  ```lua
  foo.lua → local foo = require("foo")
  ```
  
- Go modules are preloaded first, followed by embedded Lua modules. Any dependencies on Go-preloaded modules will work.

### Lua module example:

```lua
local foo = {}

function foo.bar()
  -- do something
end

return foo
```