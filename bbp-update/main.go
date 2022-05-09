package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

func main() {
	readExcludeList()

	if err := process(); err != nil {
		log.Fatalf("[FAIL] %s\n", err)
	}
}

// program structure for arkadiyt public bbp program data
type program struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Targets struct {
		InScope []programAsset `json:"in_scope"`
	} `json:"targets"`

	MaxPayout int         `json:"max_payout"` // bugcrowd payout
	MaxBounty interface{} `json:"max_bounty"` // intigriti,yeswehack payout

	OffersBounties bool `json:"offers_bounties"` // hackerone payout
	OffersSwag     bool `json:"offers_swag"`     // hackerone payout
}

type programAsset struct {
	AssetType       string `json:"asset_type"`       // hackerone (URL)
	AssetIdentifier string `json:"asset_identifier"` // hackerone

	Type     string `json:"type"`     // bugcrowd,federacy,hackenproof,intigriti,yeswehack (website,api,Web,url,web-application)
	Target   string `json:"target"`   // bugcrowd,federacy,hackenproof,yeswehack
	Endpoint string `json:"endpoint"` // intigriti
}

// chaosProgram json data item struct
type chaosProgram struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

type intigritiMaxBounty struct {
	Value float64 `json:"value"`
}

var excludeMap map[string]struct{}

func readExcludeList() {
	excludeMap = make(map[string]struct{})

	f, err := os.Open("exclude.txt")
	if err != nil {
		log.Printf("[WARN] Could not read exclude.txt: %s\n", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			excludeMap[text] = struct{}{}
		}
	}
}

func process() error {
	chaosPrograms, err := readChaosBountyPrograms()
	if err != nil {
		return err
	}
	var chaosList []chaosProgram
	for _, v := range chaosPrograms {
		chaosList = append(chaosList, v)
	}

	tempdir, err := ioutil.TempDir("", "bbp-*")
	if err != nil {
		return errors.Wrap(err, "could not create temporary directory")
	}
	defer os.RemoveAll(tempdir)

	log.Printf("[INFO] Cloning arkadiyt/bounty-targets-data repository\n")

	_, err = git.PlainClone(tempdir, false, &git.CloneOptions{
		URL:           "https://github.com/arkadiyt/bounty-targets-data",
		Progress:      os.Stdout,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.HEAD,
	})
	if err != nil {
		return errors.Wrap(err, "could not clone bounty targets data")
	}

	dataFiles := []string{"bugcrowd_data.json", "hackerone_data.json", "federacy_data.json", "hackenproof_data.json", "intigriti_data.json", "yeswehack_data.json"}
	for _, file := range dataFiles {
		log.Printf("[INFO] Reading %s data file\n", file)

		finalPath := filepath.Join(tempdir, "data", file)
		f, err := os.Open(finalPath)
		if err != nil {
			log.Printf("[WARN] Could not read %s file: %s\n", file, err)
			continue
		}
		var data []program
		if err := json.NewDecoder(f).Decode(&data); err != nil {
			log.Printf("[WARN] Could not decode %s file: %s\n", file, err)
			f.Close()
			continue
		}
		f.Close()

		for _, item := range data {
			// Exclude if program name is in exclude.txt
			if _, ok := excludeMap[item.Name]; ok {
				continue
			}
			// Skip if we already have a program with same name
			if _, ok := chaosPrograms[item.Name]; ok {
				continue
			}

			chaosItem := chaosProgram{
				Name: item.Name,
				URL:  item.URL,
			}

			// Parse the bounty and swag data from item
			switch file {
			case "hackerone_data.json":
				if item.OffersBounties {
					chaosItem.Bounty = true
				}
				if item.OffersSwag {
					chaosItem.Swag = true
				}
			case "bugcrowd_data.json", "federacy_data.json", "hackenproof_data.json":
				if item.MaxPayout > 0 {
					chaosItem.Bounty = true
				}
			case "yeswehack_data.json":
				if value, ok := item.MaxBounty.(int); ok && value > 0 {
					chaosItem.Bounty = true
				}
			case "intigriti_data.json":
				if value, ok := item.MaxBounty.(intigritiMaxBounty); ok && value.Value > 0.0 {
					chaosItem.Bounty = true
				}
			}

			for _, asset := range item.Targets.InScope {
				// Handle hackerone and skip hackerone assets which are not URL
				if asset.AssetType != "" {
					if asset.AssetType != "URL" {
						continue
					}
					if value := extractHostname(asset.AssetIdentifier); value != "" {
						chaosItem.Domains = append(chaosItem.Domains, value)
					}
				}
				if asset.Type != "" {
					if asset.Type != "website" && asset.Type != "api" && asset.Type != "Web" && asset.Type != "url" && asset.Type != "web-application" {
						continue
					}
					if value := extractHostname(asset.Target); value != "" {
						chaosItem.Domains = append(chaosItem.Domains, value)
					}
					if value := extractHostname(asset.Endpoint); value != "" {
						chaosItem.Domains = append(chaosItem.Domains, value)
					}
				}
			}
			if len(chaosItem.Domains) > 0 {
				log.Printf("[INFO] Added program %s [%s]\n", chaosItem.Name, file)
				chaosList = append(chaosList, chaosItem)
			}
		}
	}

	newFile, err := os.Create("../chaos-bugbounty-list.json")
	if err != nil {
		return errors.Wrap(err, "could not create new bbp file")
	}
	defer newFile.Close()

	marshalled, err := json.MarshalIndent(chaosList, " ", "    ")
	if err != nil {
		return errors.Wrap(err, "could not marshal chaos bbp data")
	}
	_, err = newFile.Write(marshalled)
	return err
}

type chaosList struct {
	Programs []chaosProgram `json:"programs"`
}

func readChaosBountyPrograms() (map[string]chaosProgram, error) {
	log.Printf("[INFO] Reading chaos-bugbounty-list.json\n")

	file, err := os.Open("../chaos-bugbounty-list.json")
	if err != nil {
		return nil, errors.Wrap(err, "could not read chaos list")
	}
	defer file.Close()

	var list chaosList
	if err := json.NewDecoder(file).Decode(&list); err != nil {
		return nil, errors.Wrap(err, "could not decode chaos list")
	}

	chaosMap := make(map[string]chaosProgram)
	for _, value := range list.Programs {
		chaosMap[value.Name] = value
	}
	log.Printf("[INFO] Read %d programs from chaos list\n", len(chaosMap))
	return chaosMap, nil
}

func extractHostname(item string) string {
	validate := func(value string) string {
		// Exclude if program name is in exclude.txt
		if _, ok := excludeMap[value]; ok {
			return ""
		}
		if govalidator.IsDNSName(value) {
			return value
		}
		return ""
	}
	if strings.HasPrefix(item, "http") {
		parsed, err := url.Parse(item)
		if err != nil {
			return ""
		}
		return validate(strings.TrimPrefix(parsed.Hostname(), "*."))
	}
	if strings.HasPrefix(item, "*.") {
		return validate(strings.TrimPrefix(item, "*."))
	}
	return validate(item)
}
