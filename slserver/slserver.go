package main

import (
	"github.com/fifman/surlane/src/surlane"
)

func main() {
	config := surlane.GetServerConfig()
	surlane.RunServer(&surlane.RootContext, config, surlane.DialWeb)
}
