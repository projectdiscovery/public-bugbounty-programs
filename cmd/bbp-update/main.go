package main

import (
	"log"

	"github.com/projectdiscovery/public-bugbounty-programs/pkg/dns"
)

func main() {
	dns.ReadExcludeList()

	if err := dns.Process(); err != nil {
		log.Fatalf("[FAIL] %s\n", err)
	}
}
