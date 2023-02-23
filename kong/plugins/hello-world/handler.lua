local HelloWorldHandler = {
    VERSION  = "1.0.0",
    PRIORITY = 10,
}

function HelloWorldHandler:access(conf)
    kong.service.request.enable_buffering()
    local http = require("socket.http")
    local ltn12 = require 'ltn12'

    local body = {}
    local res, code, headers, status = http.request{
        url = "http://gin-app:8080/ping",
        sink = ltn12.sink.table(body)
    }

    if code == 200 then
        local response = table.concat(body)
        kong.response.exit(200, response)
    end

end

--function HelloWorldHandler:header_filter(conf)
--    kong.response.clear_header("Content-Length")
--end

--function HelloWorldHandler:response(conf)
--    kong.log("hello from body filter")
--    kong.log(kong.response.get_raw_body())
--    --kong.response.set_raw_body("{aaaaa}")
--    --kong.response.exit(200, "Hello gwen stacy")
--end

function HelloWorldHandler:log(conf)
    kong.log("Log handler")
    zlib = require('zlib')
    kong.log(zlib.uncompress(kong.service.response.get_raw_body(), 32))
    kong.log("Log handler")
end

return HelloWorldHandler