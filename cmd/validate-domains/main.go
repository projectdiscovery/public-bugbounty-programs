package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/projectdiscovery/public-bugbounty-programs/internal/data"
	"github.com/projectdiscovery/public-bugbounty-programs/internal/dns"

	"github.com/bytedance/sonic"
)

var (
	bbListFile = flag.String("file", "dist/data.json", "Chaos bugbounty list json file")
	outputFile = flag.String("output", "invalid.txt", "Output file for invalid URLs")
	d          data.Data
)

func main() {
	flag.Parse()

	bbList, err := os.ReadFile(*bbListFile)
	if err != nil {
		log.Fatal(err)
	}

	err = sonic.Unmarshal(bbList, &d)
	if err != nil {
		log.Fatal(err)
	}

	var domains []string
	for _, program := range d.Programs {
		if len(program.Domains) > 0 {
			for _, domain := range program.Domains {
				domains = append(domains, domain)
			}
		}
	}

	var invalidDomains []string
	for _, domain := range domains {
		if !dns.ValidateFQDN(domain) {
			invalidDomains = append(invalidDomains, domain)
		}
	}

	if len(invalidDomains) == 0 {
		return
	}

	output := strings.Join(invalidDomains, "\n")

	err = os.WriteFile(*outputFile, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
