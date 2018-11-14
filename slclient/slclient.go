package main

import (
	"github.com/fifman/gosocks/src/surlane"
	"time"
)

func main() {
	config := getConfig()
	surlane.RunClient(&surlane.RootContext, config)
}

func getConfig() surlane.ClientConfig {
	return surlane.ClientConfig{
		surlane.Config{
			"123456",
			surlane.CES_128_CFB,
			1180,
			time.Second * 150,
		},
		"45.78.12.43:1190",
	}
}
