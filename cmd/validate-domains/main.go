package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/tidwall/gjson"
	"golang.org/x/net/publicsuffix"
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

	rawJSON, err := os.ReadFile(*bbListFile)
	if err != nil {
		log.Fatalf("Failed to read initial JSON file: %v", err)
	}

	var ps programs
	err = json.Unmarshal(rawJSON, &ps)
	if err != nil {
		log.Fatalf("Failed to parse initial JSON file: %v", err)
	}

	var invalidDomains []string
	gdata := gjson.ParseBytes(rawJSON)
	gdata.Get("programs.#.domains|@flatten").ForEach(func(key, value gjson.Result) bool {
		domain := value.String()
		tld, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			invalidDomains = append(invalidDomains, domain)
		}
		if (!govalidator.IsDNSName(domain) || tld != domain) && err == nil {
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
