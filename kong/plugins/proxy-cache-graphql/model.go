package main

type Config struct {
	TTLSeconds       uint `json:"ttl_seconds"`
	ErrTTLSeconds    uint `json:"err_ttl_seconds"`
	Headers          []string
	DisableNormalize bool `json:"disable_normalize"`

	LogFileSizeMaxMB               uint `json:"log_file_size_max_mb"`
	LogAgeMaxDays                  uint `json:"log_age_max_days"`
	RedisHealthCheckIntervalSecond uint `json:"redis_health_check_interval_second"`
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
