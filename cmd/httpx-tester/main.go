package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/httpx/runner"
	"github.com/projectdiscovery/public-bugbounty-programs/pkg/dns"
)

var bbListFile = flag.String("file", "../../chaos-bugbounty-list.json", "Chaos bugbounty list json file")

func main() {
	flag.Parse()

	rawJSON, err := os.ReadFile(*bbListFile)
	if err != nil {
		log.Fatalf("Failed to read initial JSON file: %v", err)
	}

	var chaosList dns.ChaosList
	err = json.Unmarshal(rawJSON, &chaosList)
	if err != nil {
		log.Fatalf("Failed to parse initial JSON file: %v", err)
	}

	allDomains := goflags.StringSlice{}

	for _, programs := range chaosList.Programs {
		allDomains = append(allDomains, programs.Domains...)
	}

	options := runner.Options{
		Methods:                "GET",
		InputTargetHost:        allDomains,
		Output:                 "invalid.txt",
		OutputFilterStatusCode: "200",
		StatusCode:             true,
		NoColor:                true,
		Timeout:                10,
	}

	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()
}
