package main

type Config struct {
	TTLSeconds       uint
	ErrTTLSeconds    uint
	Headers          []string
	DisableNormalize bool

	LogFileSizeMaxMB               uint
	LogAgeMaxDays                  uint
	RedisHealthCheckIntervalSecond uint
}

type Plugin struct {
	ID     string
	Config Config `gorm:"serializer:json"`
	Name   string
}

type GraphQLRequest struct {
	Query     string
	Variables map[string]interface{}
}
