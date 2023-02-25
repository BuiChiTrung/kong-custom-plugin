package main

import (
	"context"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/redis/go-redis/v9"
)

const NanoSecond = 1e9

var ctx = context.Background()
var requestHash string
var svc = NewService()
var gKong *pdk.PDK

type Config struct {
	TTLSeconds  int64
	VaryHeaders []string
}

func New() interface{} {
	return &Config{}
}

func (c Config) Access(kong *pdk.PDK) {
	gKong = kong
	c.GenerateCacheKey(kong)

	val, err := svc.GetCacheKey(requestHash)
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

// Response automatically enables the buffered proxy mode.
func (c Config) Response(kong *pdk.PDK) {
}

func (c Config) Log(kong *pdk.PDK) {
	responseBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		kong.Log.Err("error get service response: ", err)
	}

	kong.Log.Notice("[Log]", requestHash)
	kong.Log.Notice("[Log]", responseBody)

	if err := svc.InsertCacheKey(requestHash, responseBody, c.TTLSeconds*NanoSecond); err != nil {
		kong.Log.Err("error set redis key: ", err)
	}
}

func main() {
	Version := "1.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
