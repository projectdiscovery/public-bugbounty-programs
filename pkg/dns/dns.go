package dns

import (
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/publicsuffix"
)

var ExcludeMap map[string]struct{}

// ChaosProgram json data item struct
type ChaosProgram struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

type ChaosList struct {
	Programs []ChaosProgram `json:"programs"`
}

func ValidateFQDN(value string) string {
	// with publicsuffix package, it gets the top-level
	// domain of the input value
	tld, err := publicsuffix.EffectiveTLDPlusOne(value)

	if err != nil {
		return ""
	}

	// with govalidator package, it checks if value provided is a
	// valid domain name system (DNS) name
	if govalidator.IsDNSName(tld) {
		return tld
	}

	// If the value is not a valid DNS name, it returns an empty string
	return ""
}

func ExtractHostname(item string) string {
	item = strings.ToLower(item)

	// Exclude if program name is in exclude.txt
	if _, ok := ExcludeMap[item]; ok {
		return ""
	}

	if strings.HasPrefix(item, "http") {
		parsed, err := url.Parse(item)
		if err != nil {
			return ""
		}
		return ValidateFQDN(strings.TrimPrefix(parsed.Hostname(), "*."))
	}
	if strings.HasPrefix(item, "*.") {
		return ValidateFQDN(strings.TrimPrefix(item, "*."))
	}
	return ValidateFQDN(item)
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
