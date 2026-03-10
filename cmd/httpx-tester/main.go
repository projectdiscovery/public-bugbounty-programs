package main

import (
	"flag"
	"log"
	"os"

	"github.com/projectdiscovery/public-bugbounty-programs/internal/data"

	"github.com/bytedance/sonic"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/httpx/runner"
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

	urls := goflags.StringSlice{}
	for _, program := range d.Programs {
		urls = append(urls, program.URL)
	}

	options := runner.Options{
		Methods:               "GET",
		InputTargetHost:       urls,
		Output:                *outputFile,
		OutputMatchStatusCode: "404",
		StatusCode:            true,
		NoColor:               true,
		Timeout:               10,
		DisableStdout:         true,
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
