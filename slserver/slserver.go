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
		Config: surlane.Config{
			Password: "123456",
			Method:   surlane.CES_128_CFB,
			Port:     1190,
			Timeout:  time.Second * 60,
		},
	}
}
