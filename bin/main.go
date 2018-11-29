package main

import (
	"github.com/fifman/surlane/src/surlane"
	"fmt"
)

func main() {
	config := surlane.GetConfig()
	if config.Server == "" {
		fmt.Println("running as server")
		surlane.RunServer(&surlane.RootContext, config, surlane.DialWeb)
	} else {
		fmt.Println("running as client")
		surlane.RunClient(&surlane.RootContext, config)
	}
}