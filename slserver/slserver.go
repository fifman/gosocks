package main

import (
	"github.com/fifman/gosocks/src/surlane"
	"time"
)

func main() {
	config := getServerConfig()
	surlane.RunServer(&surlane.RootContext, config, surlane.DialWeb)
}

func getServerConfig() surlane.ServerConfig {
	return surlane.ServerConfig{
		surlane.Config{
			"123456",
			surlane.CES_128_CFB,
			1190,
			time.Second * 150,
		},
	}
}
