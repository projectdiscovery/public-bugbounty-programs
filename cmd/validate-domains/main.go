package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/asaskevich/govalidator"
)

var bbListFile = flag.String("file", "../../chaos-bugbounty-list.json", "Chaos bugbounty list json file")

type program struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

type programs struct {
	Programs []program `json:"programs"`
}

func main() {
	flag.Parse()

	rawJSON, err := ioutil.ReadFile(*bbListFile)
	if err != nil {
		log.Fatalf("Failed to read initial JSON file: %v", err)
	}

	var ps programs
	err = json.Unmarshal(rawJSON, &ps)
	if err != nil {
		log.Fatalf("Failed to parse initial JSON file: %v", err)
	}

	var invalidDomains []string
	for _, p := range ps.Programs {
		for _, domain := range p.Domains {
			if !govalidator.IsDNSName(domain) || !govalidator.IsURL(domain) {
				invalidDomains = append(invalidDomains, domain)
			}

			if strings.Count(domain, ".") > 1 {
				invalidDomains = append(invalidDomains, domain)
			}
		}
	}

	output := strings.Join(invalidDomains, "\n")

	err = ioutil.WriteFile("invalid_domains.txt", []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
