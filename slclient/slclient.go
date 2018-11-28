package main

import (
	"github.com/fifman/surlane/src/surlane"
)

func main() {
	config := surlane.GetClientConfig()
	surlane.RunClient(&surlane.RootContext, config)
}