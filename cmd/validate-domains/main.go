package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/projectdiscovery/public-bugbounty-programs/pkg/dns"
	"github.com/tidwall/gjson"
)

var bbListFile = flag.String("file", "../../chaos-bugbounty-list.json", "Chaos bugbounty list json file")

func main() {
	flag.Parse()

	rawJSON, err := os.ReadFile(*bbListFile)
	if err != nil {
		log.Fatalf("Failed to read initial JSON file: %v", err)
	}

	var ps dns.ChaosList
	err = json.Unmarshal(rawJSON, &ps)
	if err != nil {
		log.Fatalf("Failed to parse initial JSON file: %v", err)
	}

	var invalidDomains []string
	gdata := gjson.ParseBytes(rawJSON)
	gdata.Get("programs.#.domains|@flatten").ForEach(func(key, value gjson.Result) bool {
		domain := value.String()
		tld := dns.ValidateFQDN(domain)
		if tld != domain || domain == "" {
			invalidDomains = append(invalidDomains, domain)
		}
		return true
	})

	output := strings.Join(invalidDomains, "\n")

	err = os.WriteFile("invalid_domains.txt", []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
