# Kong custom plugin

## Overview 
This repo contains plugins to enhance the experience of using GraphQL API with Kong API Gateway, including:
### Proxy cache plugin
+ Used to cached GraphQL query response.

### GraphQL schema version management
+ Not available yet.

## Command
+ Start docker containers 
    ```bash
    make up
    ```
+ Stop docker containers
    ```bash
    make down
    ```
+ Update the plugin & rerun containers
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