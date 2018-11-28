package test

import (
	"log"
	"os"
	"github.com/fifman/surlane/src/surlane"
)

var (
	ERROR = log.New(os.Stderr, "ERROR: ", surlane.FLAG)
)
