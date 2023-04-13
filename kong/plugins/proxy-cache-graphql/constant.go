package main

// CacheKey kong.Ctx shared key
const CacheKey = "CacheKey"
const ResponseIsCached = "ResponseIsCached"
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
