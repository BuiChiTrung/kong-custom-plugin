package main

import (
	"fmt"
	"github.com/BuiChiTrung/go-pdk"
	"github.com/BuiChiTrung/go-pdk/server"
	"github.com/gin-gonic/gin"
	"time"
)

type Config struct {
	WaitTime int
}

func New() interface{} {
	return &Config{}
}

var requests = make(map[string]time.Time)

func (config Config) Access(kong *pdk.PDK) {
	_ = kong.Response.SetHeader("x-wait-time", fmt.Sprintf("%d seconds", config.WaitTime))

	host, _ := kong.Request.GetHost()
	lastRequest, exists := requests[host]

	if exists && time.Now().Sub(lastRequest) < time.Duration(config.WaitTime)*time.Second {
		kong.Log.Err("aádasdf")
		kong.Response.Exit(400, []byte("Maximum Requests Reached"), make(map[string][]string))
	} else {
		requests[host] = time.Now()
	}
}

func (config Config) AdminAPI(r *gin.RouterGroup) {
	r.GET("/myendpoint", func(c *gin.Context) {
		// Add your custom logic here
		c.JSON(200, gin.H{"message": "Hello from my custom endpoint"})
	})
}

func main() {
	// plugin version
	Version := "1.1"
	// plugin priority
	Priority := 1000
	_ = server.StartServer(New, Version, Priority)
}
