package main

import (
	"context"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/redis/go-redis/v9"
)

const NanoSecond = 1e9

var ctx = context.Background()
var requestHash string
var svc Service
var gKong *pdk.PDK

type Config struct {
	TTLSeconds  int
	VaryHeaders []string
}

func New() interface{} {
	return &Config{}
}

func (c Config) Access(kong *pdk.PDK) {
	gKong = kong
	c.GenerateCacheKey(kong)

	val, err := rdb.Get(ctx, requestHash).Result()
	if err == redis.Nil {
		return
	} else if err != nil {
		kong.Log.Err("error get redis key: %w", err)
	} else {
		kong.Response.SetHeader("Content-Type", "application/json")
		kong.Response.Exit(200, val, nil)
	}
}

func (c Config) GenerateCacheKey(kong *pdk.PDK) {
	requestBody, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Log.Err("err get request body: ", err)
		return
	}

	var requestHeader string
	for _, header := range c.VaryHeaders {
		headerContent, err := kong.Request.GetHeader(header)
		if err != nil {
			kong.Log.Notice(header, " header is not provided")
		}
		requestHeader += headerContent
	}

	requestHash, err = svc.GenerateCacheKey(requestBody, []byte(requestHeader))
	if err != nil {
		kong.Log.Err(err)
	}
}

func (c Config) Response(kong *pdk.PDK) {
	// The presence of the Response handler automatically enables the buffered proxy mode.
}

func (c Config) Log(kong *pdk.PDK) {
	responseBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		kong.Log.Err("error get service response: %v", err)
	}

	kong.Log.Notice("[Log]", requestHash)
	kong.Log.Notice("[Log]", responseBody)

	_, err = rdb.Get(ctx, requestHash).Result()
	if err == redis.Nil {
		if err := rdb.Set(ctx, requestHash, responseBody, time.Duration(c.TTLSeconds*NanoSecond)); err != nil {
			kong.Log.Err("error set redis key: %w", err)
		}
	}
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "kong-redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	Version := "1.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
