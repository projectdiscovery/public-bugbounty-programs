package main

import (
	"log"

	"github.com/projectdiscovery/public-bugbounty-programs/pkg/core"
)

func main() {
	core.ReadExcludeList()

	if err := core.Process(); err != nil {
		log.Fatalf("[FAIL] %s\n", err)
	}
}
