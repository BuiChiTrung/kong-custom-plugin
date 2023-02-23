local typedefs = require "kong.db.schema.typedefs"

return {
    name = "<plugin-name>",
    fields = {
        {
            -- this plugin will only be applied to Services or Routes
            consumer = typedefs.no_consumer
        },
        {
            -- this plugin will only run within Nginx HTTP module
            protocols = typedefs.protocols_http
        },
        {
            config = {
                type = "record",
                fields = {
                    -- Describe your plugin's configuration's schema here.
                },
            },
        },
    },
    entity_checks = {
        -- Describe your plugin's entity validation rules
    },
}
