package main

import (
	"fmt"
	"github.com/BuiChiTrung/go-pdk"
	"github.com/BuiChiTrung/go-pdk/server"
)

type Config struct {
}

func (config Config) Access(kong *pdk.PDK) {
	body, err := kong.Request.GetRawBody()
	var responseBody string
	if err != nil {
		responseBody = fmt.Sprintf("error get request body: %v", err)
	} else {
		responseBody = string(body) + "huhu"
	}
	kong.Response.Exit(200, []byte(responseBody), make(map[string][]string))
}

func New() interface{} {
	return &Config{}
}

func main() {
	err := server.StartServer(New, "1.0", 1001)
	if err != nil {
		fmt.Printf("error start server: %v", err)
	}
}
