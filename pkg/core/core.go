package core

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	"golang.org/x/net/publicsuffix"
)

// Program structure for arkadiyt public bbp Program data
type Program struct {
	ID string `json:"id"` // yeswehack ID

	Name    string `json:"name"`
	URL     string `json:"url"`
	Targets struct {
		InScope []ProgramAsset `json:"in_scope"`
	} `json:"targets"`

	MaxPayout int         `json:"max_payout"` // bugcrowd payout
	MaxBounty interface{} `json:"max_bounty"` // intigriti,yeswehack payout

	OffersBounties bool `json:"offers_bounties"` // hackerone payout
	OffersSwag     bool `json:"offers_swag"`     // hackerone payout
}

type ProgramAsset struct {
	AssetType       string `json:"asset_type"`       // hackerone (URL)
	AssetIdentifier string `json:"asset_identifier"` // hackerone

	Type     string `json:"type"`     // bugcrowd,federacy,hackenproof,intigriti,yeswehack (website,api,Web,url,web-application)
	Target   string `json:"target"`   // bugcrowd,federacy,hackenproof,yeswehack
	Endpoint string `json:"endpoint"` // intigriti
}

// ChaosProgram json data item struct
type ChaosProgram struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

type IntigritiMaxBounty struct {
	Value float64 `json:"value"`
}

var ExcludeMap map[string]struct{}

func ReadExcludeList() {
	ExcludeMap = make(map[string]struct{})

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
			ExcludeMap[strings.ToLower(text)] = struct{}{}
		}
	}
}

func Process() error {
	chaosPrograms, err := ReadChaosBountyPrograms()
	if err != nil {
		return err
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

	var chaosSlice []ChaosProgram
	dataFiles := []string{"bugcrowd_data.json", "hackerone_data.json", "federacy_data.json", "hackenproof_data.json", "intigriti_data.json", "yeswehack_data.json"}
	for _, file := range dataFiles {
		log.Printf("[INFO] Reading %s data file\n", file)

		finalPath := filepath.Join(tempdir, "data", file)
		f, err := os.Open(finalPath)
		if err != nil {
			log.Printf("[WARN] Could not read %s file: %s\n", file, err)
			continue
		}
		var data []Program
		if err := json.NewDecoder(f).Decode(&data); err != nil {
			log.Printf("[WARN] Could not decode %s file: %s\n", file, err)
			f.Close()
			continue
		}
		f.Close()

		for _, item := range data {
			// Fix for blank yeswehack url field
			if item.URL == "" && file == "yeswehack_data.json" {
				item.URL = fmt.Sprintf("https://yeswehack.com/programs/%s", item.ID)
			}

			// Exclude if program name is in exclude.txt
			if _, ok := ExcludeMap[strings.ToLower(item.Name)]; ok {
				continue
			}
			// Only update if we get new domains from list if the program is already
			// in our list.
			if program, ok := chaosPrograms[item.Name]; ok {
				domains := ExtractDomainsFromItem(item)
				// Dedupe and update the program if we get new domains
				new := GetUniqueDomains(program.Domains, domains)
				if len(new) > len(program.Domains) {
					program.Domains = append(program.Domains, new...)
					log.Printf("[INFO] Updated program %s (%s): %v\n", item.Name, file, new)
					chaosSlice = append(chaosSlice, program)
					delete(chaosPrograms, item.Name)
				}
				continue
			}

			chaosItem := ChaosProgram{
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
				if value, ok := item.MaxBounty.(float64); ok && value > 0 {
					chaosItem.Bounty = true
				}
			case "intigriti_data.json":
				if value, ok := item.MaxBounty.(IntigritiMaxBounty); ok && value.Value > 0.0 {
					chaosItem.Bounty = true
				}
			}

			chaosItem.Domains = ExtractDomainsFromItem(item)
			if len(chaosItem.Domains) > 0 {
				log.Printf("[INFO] Added program %s [%s]\n", chaosItem.Name, file)
				chaosSlice = append(chaosSlice, chaosItem)
			}
		}
	}

	for _, v := range chaosPrograms {
		chaosSlice = append(chaosSlice, v)
	}
	newFile, err := os.Create("../chaos-bugbounty-list.json")
	if err != nil {
		return errors.Wrap(err, "could not create new bbp file")
	}
	defer newFile.Close()

	chaosData := ChaosList{
		Programs: chaosSlice,
	}
	marshalled, err := json.MarshalIndent(chaosData, " ", "  ")
	if err != nil {
		return errors.Wrap(err, "could not marshal chaos bbp data")
	}
	_, err = newFile.Write(marshalled)
	return err
}

func ExtractDomainsFromItem(item Program) []string {
	uniqMap := make(map[string]struct{})
	var domains []string

	extractDomain := func(hostname string) {
		if hostname == "" {
			return
		}
		if value := ExtractHostname(hostname); value != "" {
			if _, ok := uniqMap[value]; ok {
				return
			}
			uniqMap[value] = struct{}{}
			domains = append(domains, value)
		}
	}
	for _, asset := range item.Targets.InScope {
		// Handle hackerone and skip hackerone assets which are not URL
		if asset.AssetType != "" {
			if asset.AssetType != "URL" {
				continue
			}
			extractDomain(asset.AssetIdentifier)
		}
		if asset.Type != "" {
			if asset.Type != "website" && asset.Type != "api" && asset.Type != "Web" && asset.Type != "url" && asset.Type != "web-application" {
				continue
			}
			extractDomain(asset.Target)
			extractDomain(asset.Endpoint)
		}
	}
	return domains
}

type ChaosList struct {
	Programs []ChaosProgram `json:"programs"`
}

func ReadChaosBountyPrograms() (map[string]ChaosProgram, error) {
	log.Printf("[INFO] Reading chaos-bugbounty-list.json\n")

	file, err := os.Open("../chaos-bugbounty-list.json")
	if err != nil {
		return nil, errors.Wrap(err, "could not read chaos list")
	}
	defer file.Close()

	var list ChaosList
	if err := json.NewDecoder(file).Decode(&list); err != nil {
		return nil, errors.Wrap(err, "could not decode chaos list")
	}

	chaosMap := make(map[string]ChaosProgram)
	for _, value := range list.Programs {
		chaosMap[value.Name] = value
	}
	log.Printf("[INFO] Read %d programs from chaos list\n", len(chaosMap))
	return chaosMap, nil
}

func ExtractHostname(item string) string {
	item = strings.ToLower(item)

	validate := func(value string) string {
		tld, err := publicsuffix.EffectiveTLDPlusOne(value)
		if err != nil {
			return ""
		}
		// Exclude if program name is in exclude.txt
		if _, ok := ExcludeMap[tld]; ok {
			return ""
		}
		if govalidator.IsDNSName(tld) {
			return tld
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

func GetUniqueDomains(first, second []string) []string {
	uniq := make(map[string]struct{})
	for _, v := range first {
		uniq[v] = struct{}{}
	}
	var unique []string
	for _, v := range second {
		if _, ok := uniq[v]; !ok {
			unique = append(unique, v)
			uniq[v] = struct{}{}
		}
	}
	return unique
}
