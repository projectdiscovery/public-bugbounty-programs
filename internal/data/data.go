package data

type Program struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Domains []string `json:"domains"`
	Swag    bool     `json:"swag,omitempty"`
}

type Data struct {
	Programs []Program `json:"programs"`
}

// SourceProgram mirrors entries from arkadiyt/bounty-targets-data JSON files.
type SourceProgram struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	URL     string        `json:"url"`
	Targets SourceTargets `json:"targets"`

	MaxPayout int `json:"max_payout"`
	MaxBounty any `json:"max_bounty"`

	OffersBounties bool `json:"offers_bounties"`
	OffersSwag     bool `json:"offers_swag"`
}

type SourceTargets struct {
	InScope []SourceProgramAsset `json:"in_scope"`
}

type SourceProgramAsset struct {
	AssetType       string `json:"asset_type"`
	AssetIdentifier string `json:"asset_identifier"`

	Type     string `json:"type"`
	Target   string `json:"target"`
	Endpoint string `json:"endpoint"`
}
