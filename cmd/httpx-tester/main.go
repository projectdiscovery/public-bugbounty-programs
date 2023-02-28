package main

import (
	"flag"
	"log"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/httpx/runner"
	fileutil "github.com/projectdiscovery/utils/file"
)

var bbListFile = flag.String("file", "../../urls.txt", "Chaos bugbounty list json file")

func main() {
	flag.Parse()

	urlFile, err := fileutil.ReadFile(*bbListFile)
	if err != nil {
		log.Fatal(err)
	}

	allUrls := goflags.StringSlice{}

	for url := range urlFile {
		allUrls = append(allUrls, url)
	}

	options := runner.Options{
		Methods:               "GET",
		InputTargetHost:       allUrls,
		Output:                "invalid.txt",
		OutputMatchStatusCode: "404",
		StatusCode:            true,
		NoColor:               true,
		Timeout:               10,
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
