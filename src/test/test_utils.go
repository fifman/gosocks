package test

import (
	"log"
	"os"
	"github.com/fifman/gosocks/src/surlane"
)

var (
	ERROR = log.New(os.Stderr, "ERROR: ", surlane.FLAG)
)
