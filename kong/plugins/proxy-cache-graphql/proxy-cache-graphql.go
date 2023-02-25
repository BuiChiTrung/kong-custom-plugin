package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/redis/go-redis/v9"
)

const NanoSecond = 1e9

var rdb *redis.Client
var ctx context.Context
var requestHash string

type Config struct {
	TTLSeconds  int
	VaryHeaders []string
}

func New() interface{} {
	return &Config{}
}

type GraphQLRequest struct {
	Query string
}

func (c Config) Access(kong *pdk.PDK) {
	ctx := context.Background()

	requestBody, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Log.Err("error get request body: %v", err)
		return
	}

	graphQLAstBytes, err := Az(string(requestBody))
	if err != nil {
		kong.Log.Err(err)
		return
	}
	kong.Log.Notice(string(graphQLAstBytes))

	var requestHeader string
	for _, header := range c.VaryHeaders {
		headerContent, err := kong.Request.GetHeader(header)
		if err != nil {
			kong.Log.Notice(header, " header is not provided")
		}
		requestHeader += headerContent
	}

	request := append([]byte(requestHeader), graphQLAstBytes...)
	requestHashByte := md5.Sum(request)
	requestHash = fmt.Sprintf("%x", string(requestHashByte[:]))

	val, err := rdb.Get(ctx, requestHash).Result()
	if err == redis.Nil {
		return
	} else if err != nil {
		kong.Log.Err("error get redis")
	} else {
		kong.Response.SetHeader("Content-Type", "application/json")
		kong.Response.Exit(200, val, nil)
	}
}

func (c Config) Response(kong *pdk.PDK) {
	// The presence of the Response handler automatically enables the buffered proxy mode.
}

func (c Config) Log(kong *pdk.PDK) {
	currentBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		kong.Log.Err("error get service response: %v", err)
	}

	kong.Log.Notice("ahah")
	kong.Log.Notice(requestHash)
	kong.Log.Notice(currentBody)
	kong.Log.Notice("ahah")

	_, err = rdb.Get(ctx, requestHash).Result()
	if err == redis.Nil {
		if err := rdb.Set(ctx, requestHash, currentBody, time.Duration(c.TTLSeconds*NanoSecond)); err != nil {
			kong.Log.Err("error set redis: %v", err)
		}
	}
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "kong-redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx = context.Background()

	Version := "1.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
