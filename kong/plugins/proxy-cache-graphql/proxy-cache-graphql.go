package main

import (
	"fmt"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

type Config struct {
	TTLSeconds       uint
	ErrTTLSeconds    uint
	Headers          []string
	DisableNormalize bool

	TurnOffRedis bool
}

var gKong *pdk.PDK
var gConf Config
var gSvc *Service

func New() interface{} {
	gConf = Config{}

	var err error
	gSvc, err = NewService()
	if err != nil {
		panic(err)
	}

	return &gConf
}

func (c Config) Access(kong *pdk.PDK) {
	// TODO: trung.bc - TD
	_ = kong.ServiceRequest.ClearHeader("Accept-Encoding")

	gKong = kong
	cacheKey, shouldCached, err := c.GenerateCacheKey(kong)
	if err != nil {
		_ = kong.Log.Err(err.Error())
		return
	}
	if !shouldCached {
		_ = kong.Response.SetHeader("X-Cache-Status", string(Bypass))
		return
	}

	if err := kong.Ctx.SetShared(CacheKey, cacheKey); err != nil {
		_ = kong.Log.Err("err set shared context: ", err.Error())
		return
	}

	cacheVal, err := gSvc.GetCacheKey(cacheKey)
	if err != nil {
		_ = kong.Response.SetHeader("X-Cache-Status", string(Miss))
		if err == redis.Nil {
			_ = kong.Response.SetHeader("X-Cache-Key", cacheKey)
			return
		}
		kong.Log.Err("error get redis key: %w", err)
	} else {
		_ = kong.Ctx.SetShared(ResponseAlreadyCached, true)

		kong.Response.Exit(200, cacheVal, map[string][]string{
			"Content-Type":                {"application/json"},
			"X-Cache-Key":                 {cacheKey},
			"X-Cache-Status":              {string(Hit)},
			"access-control-allow-origin": {"*"},
		})
	}
}

func (c Config) GenerateCacheKey(kong *pdk.PDK) (cacheKey string, shouldCached bool, err error) {
	requestBody, err := kong.Request.GetRawBody()
	if err != nil {
		return "", false, fmt.Errorf("err get request body: %w", err)
	}

	var requestHeader string
	for _, header := range c.Headers {
		headerContent, _ := kong.Request.GetHeader(header)
		requestHeader += headerContent
	}

	requestPath, err := kong.Request.GetPath()
	if err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey get req path")
	}

	cacheKey, shouldCached, err = gSvc.GenerateCacheKey(string(requestBody), requestHeader, requestPath)
	if err != nil {
		return "", false, err
	}

	return cacheKey, shouldCached, nil
}

// Response automatically enables the buffered proxy mode.
func (c Config) Response(kong *pdk.PDK) {
}

func (c Config) Log(kong *pdk.PDK) {
	isCache, err := kong.Ctx.GetSharedAny(ResponseAlreadyCached)
	if v, ok := isCache.(bool); ok && v {
		return
	}

	responseBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		_ = kong.Log.Err("error get service response: ", err)
		return
	}

	cacheKey, err := kong.Ctx.GetSharedString(CacheKey)
	if err != nil {
		_ = kong.Log.Err("err get shared context: ", err.Error())
	}
	if cacheKey == "" {
		return
	}

	//_ = kong.Log.Notice("[Log]", cacheKey)
	//_ = kong.Log.Notice("[Log]", responseBody)

	c.InsertCacheKey(kong, cacheKey, responseBody)
}

func (c Config) InsertCacheKey(kong *pdk.PDK, cacheKey string, cacheValue string) {
	statusCode, _ := kong.ServiceResponse.GetStatus()
	if statusCode >= http.StatusInternalServerError {
		return
	} else if statusCode >= http.StatusBadRequest && c.ErrTTLSeconds > 0 {
		if err := gSvc.InsertCacheKey(cacheKey, cacheValue, int64(c.ErrTTLSeconds)*int64(NanoSecond)); err != nil {
			_ = kong.Log.Err("error set redis key: ", err)
		}
	} else {
		var ttlSeconds uint

		ttlHeaderStr, _ := kong.Request.GetHeader(TTLSeconds)
		ttlHeader, err := strconv.Atoi(ttlHeaderStr)

		if ttlHeaderNotProvideOrInvalid := err != nil; ttlHeaderNotProvideOrInvalid {
			ttlSeconds = c.TTLSeconds
		} else if ttlHeader < 0 {
			return
		} else {
			ttlSeconds = uint(ttlHeader)
		}

		if err := gSvc.InsertCacheKey(cacheKey, cacheValue, int64(ttlSeconds)*int64(NanoSecond)); err != nil {
			_ = kong.Log.Err("error set redis key: ", err)
		}
	}
}

func main() {
	Version := "1.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
