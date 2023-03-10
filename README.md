# Kong custom plugin
This repo contains plugins to enhance the experience of using GraphQL API with Kong API Gateway, including:

## Proxy cache plugin
Used to cached GraphQL query response. The plugin support most of the features provided by Kong Enterprise [GraphQL Proxy Caching Advanced](https://docs.konghq.com/hub/kong-inc/graphql-proxy-cache-advanced/) plugin along with additional features:
+ Recognize mutilple requests are the same and returned the cached result already created when making the first request instead of making query to upstream service again. Support cases:
  + Ex1: All requests contain the same fields but the order are not identical:
    ```
    # 1st query
    query Country {
      country(code: "VN") {
          emoji,    
          name,
          native,
          capital
	    }
    }
        
    # 2nd query: field order are changed, contain comment
    query Country {
      country(code: "VN") {
          emoji, name,
          capital, # graphql cmt
          native
      }
    }
    
    # 3rd query: operation name, type are omitted
    {
        country(code: "VN") {
            native, emoji, name, capital
        } 
    }
    ```
  
  + Ex2: All requests use the same variables but in declare, use in different order
  
    ```
    # 1st query
    query Repository($name: String!, $owner: String!, $followRenames: Boolean) {
      repository(name: $name, owner: $owner, followRenames: $followRenames) {
        allowUpdateBranch
        autoMergeAllowed
        id
        createdAt
        owner {
          avatarUrl
          id
          login
          url
        }
        isPrivate
      }
    }
    
    # variables
    {
      "name": "kong-custom-plugin",
      "owner": "BuiChiTrung",
      "followRenames": true
    }
    
    # 2nd query
    query Repository($followRenames: Boolean, $owner: String!, $name: String!) {
      repository(name: $name, owner: $owner, followRenames: $followRenames) {
        owner {
          avatarUrl
          id
          login
          url
        }
        allowUpdateBranch
        autoMergeAllowed
        createdAt
        id
        isPrivate
      }
    }
    
    # variables
    {
      "name": "kong-custom-plugin",
      "followRenames": true,
      "owner": "BuiChiTrung"
    }
    ```

### Cache key

`MD5(requestHeader | graphQLAST | graphQL Variable)`

+ `requestHeader`: the content of headers specified in plugin config
+ `graphQLAST`: the AST of graphQL query

### Plugin config

+ `ttlseconds` - `int`: the cache key time to live seconds
+  `varyheaders` - `[]string`: the list of header used to create cache key. 

### Plugin APIs

Written using Fiber framework: `./fiber-app`. Cache key used for these APIs can be retrieved in the `X-Cache-Key` in response.

+ GET cache entity:

  ```bash
  curl --location 'http://localhost:9080/proxy-cache/<cache_key>' \
  ```

+ UPSERT cache entity:

  ```bash
  curl --location 'http://localhost:9080/proxy-cache' \
  --data '{
      "cacheKey": "<cache-key>",
      "value": "<new-cache-value>"
  }'
  ```

+ DELETE cache entity:

  ```bash
  curl --location --request DELETE 'http://localhost:9080/proxy-cache/<cache-key>'
  ```

+ PURGE cache entity: remove all cache keys

  ```bash
  curl --location --request DELETE 'http://localhost:9080/proxy-cache'
  ```

## GraphQL schema version management
+ Coming soon. 

## Installation

### Setup kong and the plugin
+ The docker-compose file follows [Kong installation document via Docker](https://docs.konghq.com/gateway/latest/install/docker/).
+ Start docker containers
    ```bash
    make up
    ```
+ **Note**: At the second time of running these containers you should comment out the `kong-gateway-db-bootstrap` service in `docker-compose.yaml` as it only used as the first time you run these container to migration Kong initial db to the psql volumn: `./db/psql`.

### Enable the plugin
+ Globally
```bash
curl --location 'http://localhost:8001/plugins' \
--data '{
    "name": "proxy-cache-graphql",
    "config": {
        "ttlseconds":  5
    }
}'
```

+ For a single service
```bash
curl --location 'http://localhost:8001/services/SERVICE_NAME|SERVICE_ID/plugins' \
--data '{
    "name": "proxy-cache-graphql",
    "config": {}
}'
```

+ For a single route
```bash
curl --location 'http://localhost:8001/routes/ROUTE_NAME|ROUTE_ID/plugins' \
--data '{
    "name": "proxy-cache-graphql",
    "config": {}
}'
```


### Other commands
+ Stop docker containers
    ```bash
    make down
    ```
+ Update the plugin source code & rerun containers
    ```bash
    make restart
    ```
+ Run fiber app separately
    ```bash
    make run fiber 
    ```

## Folder structure
+ bin: binary built Go file.
+ db: docker volume to store Kong db schema (it will be removed in next versions).
+ fiber-app: an application server written using fiber framework which provide API to all plugins.
+ kong/plugins: plugin source code
  + go-wait: test plugin written in Go.
  + hello-world: test plugin written in Lua.
  + proxy-cache-graphql: proxy-cache-graphql plugin.