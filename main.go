package main

import (
	"github.com/fifman/surlane/src/surlane"
	"fmt"
	"github.com/fifman/kcptun/src"
)

func main() {
	config := surlane.GetConfig()
	kcpConfig := config.Kcptun
	if config.Server == "" {
		fmt.Println("running as server")
		if config.Kcp {
			kcpConfig.Listen = fmt.Sprintf(":%d", config.Port)
			kcpConfig.Target = ":29901"
			kcpConfig.Key = config.Password
			src.PostLoadConfig(&kcpConfig)
			config.Port = 29901
			serverConfig := src.ServerConfig{
				false, kcpConfig,
			}
			go src.RunServer(serverConfig)
		}
		surlane.RunServer(&surlane.RootContext, config, surlane.DialWeb)
	} else {
		fmt.Println("running as client")
		if config.Kcp {
			kcpConfig.Listen = ":29901"
			kcpConfig.Target = config.Server
			config.Server = ":29901"
			kcpConfig.Key = config.Password
			src.PostLoadConfig(&kcpConfig)
			go src.RunClient(src.ClientConfig{
				1, 0, 600, kcpConfig,
			})
		}
		surlane.RunClient(&surlane.RootContext, config)
	}
}