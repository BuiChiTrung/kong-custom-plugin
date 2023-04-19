package main

// CacheKey kong.Ctx shared key
const CacheKey = "CacheKey"
const ResponseAlreadyCached = "ResponseAlreadyCached"
const NanoSecond = 1e9
const TTLSeconds = "TTL-Seconds"

type CacheStatus string

const (
	Hit    CacheStatus = "Hit"
	Miss   CacheStatus = "Miss"
	Bypass CacheStatus = "Bypass"
)

type OperationType string

const (
	Query        OperationType = "query"
	Mutation     OperationType = "mutation"
	Subscription OperationType = "subscription"
)

const (
	HeaderXCacheStatus             = "X-Cache-Status"
	HeaderAcceptEncoding           = "Accept-Encoding"
	HeaderContentType              = "Content-Type"
	HeaderXCacheKey                = "X-Cache-Key"
	HeaderAccessControlAllowOrigin = "access-control-allow-origin"
)

const (
	EnvRedisMasterHost   = "KONG_REDIS_MASTER_HOST"
	EnvRedisMasterPort   = "KONG_REDIS_MASTER_PORT"
	EnvRedisReplicasHost = "KONG_REDIS_REPLICAS_HOST"
	EnvRedisReplicasPort = "KONG_REDIS_REPLICAS_PORT"
)
