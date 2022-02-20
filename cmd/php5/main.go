package main

import (
	"log"
	"os"

	"github.com/Lofanmi/yam/internal/core/php5"
)

var _ = php5.ImportMe

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {}
