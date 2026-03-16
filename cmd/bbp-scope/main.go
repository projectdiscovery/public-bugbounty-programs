package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/projectdiscovery/public-bugbounty-programs/internal/data"
	"github.com/projectdiscovery/public-bugbounty-programs/internal/dns"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	outData         = flag.String("out", "src/data.yaml", "Output file path for generated bounty-targets data")
	excludeFilePath = flag.String("exclude", "src/exclude.txt", "Path to newline-delimited program exclusion list")
	repoURL         = flag.String("repo", "https://github.com/arkadiyt/bounty-targets-data", "Source repository URL for bounty target data")
)

var srcDataFiles = []string{
	"bugcrowd_data.json",
	"hackerone_data.json",
	"federacy_data.json",
	// "hackenproof_data.json", // NOTE(dwisiswant0): hackenproof data file is currently unavailable.
	"intigriti_data.json",
	"yeswehack_data.json",
}

var webAssetTypes = map[string]struct{}{
	"website":         {},
	"api":             {},
	"web":             {},
	"url":             {},
	"web-application": {},
}

type IntigritiMaxBounty struct {
	Value float64 `json:"value"`
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}

func run() error {
	loadExcludeList(*excludeFilePath)

	chaosPrograms, err := readBBPrograms(*outData)
	if err != nil {
		return err
	}

	mergedPrograms, err := buildProgramList(chaosPrograms)
	if err != nil {
		return err
	}

	return writeBBPrograms(*outData, mergedPrograms)
}

func loadExcludeList(path string) {
	dns.ExcludeMap = make(map[string]struct{})

	f, err := os.Open(path)
	if err != nil {
		log.Printf("Could not read %s: %s\n", path, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			dns.ExcludeMap[strings.ToLower(text)] = struct{}{}
		}
	}
}

func buildProgramList(existing map[string]data.Program) ([]data.Program, error) {
	tempDir, err := os.MkdirTemp("", "bbp-*")
	if err != nil {
		return nil, errors.Wrap(err, "could not create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	if err := dlDataFiles(tempDir); err != nil {
		return nil, err
	}

	var merged []data.Program
	for _, file := range srcDataFiles {
		log.Printf("Reading %s data file\n", file)

		dataPath := filepath.Join(tempDir, "data", file)
		items, err := readSourcePrograms(dataPath)
		if err != nil {
			log.Printf("Could not parse %s file: %s\n", file, err)
			continue
		}

		for _, item := range items {
			normalizeProgramURL(&item, file)
			if isExcludedProgram(item.Name) {
				continue
			}

			domains := extractDomainsFromItem(item)

			if program, ok := existing[item.Name]; ok {
				if updated := mergeProgramDomains(&program, domains); updated {
					log.Printf("Updated program %s (%s): %v\n", item.Name, file, program.Domains)
					merged = append(merged, program)
					delete(existing, item.Name)
				}
				continue
			}

			chaosItem := data.Program{
				Name:    item.Name,
				URL:     item.URL,
				Domains: domains,
			}
			setProgramRewards(&chaosItem, item, file)

			if len(chaosItem.Domains) == 0 {
				continue
			}

			log.Printf("Added program %s [%s]\n", chaosItem.Name, file)
			merged = append(merged, chaosItem)
		}
	}

	for _, program := range existing {
		merged = append(merged, program)
	}

	sort.Slice(merged, func(i, j int) bool {
		return strings.ToLower(merged[i].Name) < strings.ToLower(merged[j].Name)
	})

	return merged, nil
}

func dlDataFiles(tempDir string) error {
	log.Printf("Downloading bounty-targets-data source files\n")

	dataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return errors.Wrap(err, "could not create temporary data directory")
	}

	downloaded := 0
	for _, file := range srcDataFiles {
		if err := dlSrcDataFile(dataDir, file); err != nil {
			log.Printf("Could not download %s: %s\n", file, err)
			continue
		}
		downloaded++
	}

	if downloaded == 0 {
		return errors.New("could not download any source data file")
	}

	return nil
}

func dlSrcDataFile(destDir, file string) error {
	rawURL, err := buildGitHubRawURL(*repoURL, "main", "data/"+file)
	if err != nil {
		return err
	}

	if err := dlFile(rawURL, filepath.Join(destDir, file)); err != nil {
		return errors.Wrapf(err, "could not download %s", file)
	}

	log.Printf("Downloaded %s (main)\n", file)
	return nil
}

func buildGitHubRawURL(repo, branch, path string) (string, error) {
	trimmed := strings.TrimSpace(repo)
	trimmed = strings.TrimSuffix(trimmed, ".git")
	trimmed = strings.TrimPrefix(trimmed, "https://github.com/")
	trimmed = strings.TrimPrefix(trimmed, "http://github.com/")
	trimmed = strings.TrimPrefix(trimmed, "github.com/")
	trimmed = strings.Trim(trimmed, "/")

	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", errors.Errorf("unsupported repository URL %q; expected github owner/repo", repo)
	}

	owner, repoName := parts[0], parts[1]
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repoName, branch, path), nil
}

func dlFile(url, destPath string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("request failed: %s", resp.Status)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return nil
}

func readSourcePrograms(path string) ([]data.SourceProgram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data []data.SourceProgram
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func normalizeProgramURL(item *data.SourceProgram, file string) {
	if item.URL == "" && file == "yeswehack_data.json" {
		item.URL = fmt.Sprintf("https://yeswehack.com/programs/%s", item.ID)
	}
}

func isExcludedProgram(name string) bool {
	_, ok := dns.ExcludeMap[strings.ToLower(name)]
	return ok
}

func mergeProgramDomains(program *data.Program, domains []string) bool {
	newDomains := dns.GetUniqueDomains(program.Domains, domains)
	if len(newDomains) == 0 {
		return false
	}
	program.Domains = append(program.Domains, newDomains...)
	return true
}

func setProgramRewards(program *data.Program, item data.SourceProgram, sourceFile string) {
	switch sourceFile {
	case "hackerone_data.json":
		program.Bounty = item.OffersBounties
		program.Swag = item.OffersSwag
	case "bugcrowd_data.json", "federacy_data.json", "hackenproof_data.json":
		program.Bounty = item.MaxPayout > 0
	case "yeswehack_data.json":
		if value, ok := item.MaxBounty.(float64); ok {
			program.Bounty = value > 0
		}
	case "intigriti_data.json":
		program.Bounty = parseIntigritiBounty(item.MaxBounty) > 0
	}
}

func parseIntigritiBounty(value interface{}) float64 {
	switch v := value.(type) {
	case IntigritiMaxBounty:
		return v.Value
	case map[string]interface{}:
		if raw, ok := v["value"]; ok {
			if amount, ok := raw.(float64); ok {
				return amount
			}
		}
	}
	return 0
}

func writeBBPrograms(path string, programs []data.Program) error {
	newFile, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "could not create new bbp file")
	}
	defer newFile.Close()

	chaosData := data.Data{
		Programs: programs,
	}

	encoder := yaml.NewEncoder(newFile)
	encoder.SetIndent(2)
	if err := encoder.Encode(chaosData); err != nil {
		return errors.Wrap(err, "could not marshal chaos bbp data")
	}

	return encoder.Close()
}

func extractDomainsFromItem(item data.SourceProgram) []string {
	uniqMap := make(map[string]struct{})
	var domains []string

	extractDomain := func(hostname string) {
		if hostname == "" {
			return
		}
		if value := dns.ExtractHostname(hostname); value != "" {
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
			if _, ok := webAssetTypes[strings.ToLower(asset.Type)]; !ok {
				continue
			}
			extractDomain(asset.Target)
			extractDomain(asset.Endpoint)
		}
	}
	return domains
}

func readBBPrograms(path string) (map[string]data.Program, error) {
	log.Printf("Reading %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not read chaos list")
	}
	defer file.Close()

	var list data.Data
	if err := yaml.NewDecoder(file).Decode(&list); err != nil {
		return nil, errors.Wrap(err, "could not decode chaos list")
	}

	chaosMap := make(map[string]data.Program)
	for _, value := range list.Programs {
		chaosMap[value.Name] = value
	}
	log.Printf("Read %d programs from chaos list\n", len(chaosMap))
	return chaosMap, nil
}
