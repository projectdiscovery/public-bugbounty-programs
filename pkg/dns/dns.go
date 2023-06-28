package dns

import (
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/retryabledns"

	"github.com/asaskevich/govalidator"
	sliceutil "github.com/projectdiscovery/utils/slice"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"golang.org/x/net/publicsuffix"
)

// DefaultResolvers contains the default list of resolvers known to be good

var (
	ExcludeMap       map[string]struct{}
	DefaultResolvers = []string{
		"1.1.1.1:53",        // Cloudflare primary
		"1.0.0.1:53",        // Cloudflare secondary
		"8.8.8.8:53",        // Google primary
		"8.8.4.4:53",        // Google secondary
		"9.9.9.9:53",        // Quad9 Primary
		"9.9.9.10:53",       // Quad9 Secondary
		"77.88.8.8:53",      // Yandex Primary
		"77.88.8.1:53",      // Yandex Secondary
		"208.67.222.222:53", // OpenDNS Primary
		"208.67.220.220:53", // OpenDNS Secondary
	}
	client *retryabledns.Client
)

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

func ValidateFQDN(value string) bool {
	// check if domain can can be parsed
	tld, err := publicsuffix.EffectiveTLDPlusOne(value)
	if err != nil {
		// If the domain can't be parsed by publicsuffix,
		// then we attempt to resolve a DNS A record to determine if it's valid.
		resp, err := client.Resolve(value)
		if err != nil || (resp.A == nil && resp.AAAA == nil) {
			// DNS resolution also failed, so we conclude that the domain isn't valid.
			return false
		}
		return true
	}

	// check if top level domain is equal to original and it's a valid domain name
	return tld == value && govalidator.IsDNSName(tld)
}

func ExtractHostname(item string) string {
	item = strings.ToLower(item)

	// Exclude if program name is in exclude.txt
	if _, ok := ExcludeMap[item]; ok {
		return ""
	}

	trimmedStr := stringsutil.TrimPrefixAny(item, "http://", "https://", "*.")

	if ValidateFQDN(trimmedStr) {
		return trimmedStr
	}
	return ""
}

func GetUniqueDomains(first, second []string) []string {
	_, diff := sliceutil.Diff(first, second)
	return sliceutil.Dedupe(diff)
}

func init() {
	var err error
	client, err = retryabledns.New(DefaultResolvers, 3)
	if err != nil || client == nil {
		gologger.Fatal().Msgf("Could not create DNS client: %s\n", err)
	}
}
