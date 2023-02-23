local endpoints = require "kong.api.endpoints"

local consumers_schema = kong.db.consumers.schema

return {
    ["/hello-world/test"] = {
        schema = consumers_schema,
        methods = {
            GET = function(self, db, helpers)
                return kong.response.exit(200, { message = "Lua is a stupid language" })
            end,
        },
    },
}