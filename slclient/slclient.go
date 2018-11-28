package main

import (
	"github.com/fifman/gosocks/src/surlane"
)

func main() {
	config := surlane.GetClientConfig()
	surlane.RunClient(&surlane.RootContext, config)
}
