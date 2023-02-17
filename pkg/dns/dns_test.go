package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFQDN(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "CorrectDomain",
			input: "example.com",
			want:  true,
		},
		{
			name:  "CorrectDomainSecondLevelDomain",
			input: "example.co.uk",
			want:  true,
		},
		{
			name:  "IncorrectDomainSecondLevelDomain",
			input: "docs.example.com",
			want:  false,
		},
		{
			name:  "ThirdLevelDomain",
			input: "a.a.example.com",
			want:  false,
		},
		{
			name:  "multiLevelDomain",
			input: "a.a.a.a.a.a.a.example.com",
			want:  false,
		},
		{
			name:  "WildcardDomain",
			input: "*.example.com",
			want:  false,
		},
		{
			name:  "MultiWildcardDomain",
			input: "*.aaaaa.*.example.com",
			want:  false,
		},
		{
			name:  "HttpUrlDomain",
			input: "http://example.com",
			want:  false,
		},
		{
			name:  "HttpsUrlDomain",
			input: "https://example.com",
			want:  false,
		},
		{
			name:  "EmptyDomain",
			input: "",
			want:  false,
		},
		{
			name:  "SpaceDomain",
			input: " ",
			want:  false,
		},
		{
			name:  "MultiSpaceDomain",
			input: "        ",
			want:  false,
		},
		{
			name:  "InvalidRegexDomain",
			input: "exa$mple.com",
			want:  false,
		},
		{
			name:  "InvalidRegexDomain2",
			input: "ex@$mpl&.com",
			want:  false,
		},
		{
			name:  "InvalidRegexSubomain",
			input: "some$$thing.example.com",
			want:  false,
		},
		{
			name:  "InvalidRegexSubomain2",
			input: "some$$thing.examp%&le.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFQDN(tt.input)
			require.Equalf(t, tt.want, got, "test %s => wanted %v but got %v", tt.name, tt.want, got)
		})

	}
}
