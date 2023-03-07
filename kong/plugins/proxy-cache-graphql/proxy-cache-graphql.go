package main

import (
	"fmt"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/redis/go-redis/v9"
)

var gKong *pdk.PDK

type Config struct {
	TTLSeconds  int
	VaryHeaders []string
	svc         *Service
}

func New() interface{} {
	return &Config{
		svc: NewService(),
	}
}

func (c Config) Access(kong *pdk.PDK) {
	gKong = kong
	cacheKey, shouldCached, err := c.GenerateCacheKey(kong)
	if err != nil {
		kong.Log.Err(err.Error())
		return
	}
	if !shouldCached {
		kong.Response.SetHeader("X-Cache-Status", string(Bypass))
		return
	}

	if err := kong.Ctx.SetShared(CacheKey, cacheKey); err != nil {
		kong.Log.Err("err set shared context: ", err.Error())
		return
	}

	cacheVal, err := c.svc.GetCacheKey(cacheKey)
	if err != nil {
		kong.Response.SetHeader("X-Cache-Status", string(Miss))
		if err == redis.Nil {
			return
		}
		kong.Log.Err("error get redis key: %w", err)
	} else {
		kong.Response.SetHeader("Content-Type", "application/json")
		kong.Response.SetHeader("X-Cache-Key", cacheKey)
		kong.Response.SetHeader("X-Cache-Status", string(Hit))
		kong.Response.Exit(200, cacheVal, nil)
	}
}

func (c Config) GenerateCacheKey(kong *pdk.PDK) (cacheKey string, shouldCached bool, err error) {
	requestBody, err := kong.Request.GetRawBody()
	if err != nil {
		return "", false, fmt.Errorf("err get request body: %w", err)
	}

	var requestHeader string
	for _, header := range c.VaryHeaders {
		headerContent, _ := kong.Request.GetHeader(header)
		requestHeader += headerContent
	}

	requestPath, err := kong.Request.GetPath()
	if err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey get req path")
	}

	cacheKey, shouldCached, err = c.svc.GenerateCacheKey(string(requestBody), requestHeader, requestPath)
	if err != nil {
		return "", false, err
	}

	return cacheKey, shouldCached, nil
}

// Response automatically enables the buffered proxy mode.
func (c Config) Response(kong *pdk.PDK) {
}

func (c Config) Log(kong *pdk.PDK) {
	responseBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		kong.Log.Err("error get service response: ", err)
	}

	cacheKey, err := kong.Ctx.GetSharedString(CacheKey)
	if err != nil {
		kong.Log.Err("err get shared context: ", err.Error())
	}

	kong.Log.Notice("[Log]", cacheKey)
	kong.Log.Notice("[Log]", responseBody)

	if responseBody == "" {
		return
	}

	if err := c.svc.InsertCacheKey(cacheKey, responseBody, int64(c.TTLSeconds)*int64(NanoSecond)); err != nil {
		kong.Log.Err("error set redis key: ", err)
	}
}

func main() {
	Version := "1.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
