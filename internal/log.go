package internal

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "uptrace: ", log.Lshortfile)
