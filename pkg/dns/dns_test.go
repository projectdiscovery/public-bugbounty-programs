package dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		isInvalid bool
	}{
		{
			name:      "CorrectDomain",
			input:     "example.com",
			want:      "example.com",
			isInvalid: false,
		},
		{
			name:      "CorrectDomainSecondLevelDomain",
			input:     "example.co.uk",
			want:      "example.co.uk",
			isInvalid: false,
		},
		{
			name:      "IncorrectDomainSecondLevelDomain",
			input:     "docs.example.com",
			want:      "example.com",
			isInvalid: true,
		},
		{
			name:      "ThirdLevelDomain",
			input:     "a.a.example.com",
			want:      "example.com",
			isInvalid: true,
		},
		{
			name:      "multiLevelDomain",
			input:     "a.a.a.a.a.a.a.example.com",
			want:      "example.com",
			isInvalid: true,
		},
		{
			name:      "WildcardDomain",
			input:     "*.example.com",
			want:      "example.com",
			isInvalid: true,
		},
		{
			name:      "MultiWildcardDomain",
			input:     "*.aaaaa.*.example.com",
			want:      "example.com",
			isInvalid: true,
		},
		{
			name:      "HttpUrlDomain",
			input:     "http://example.com",
			want:      "",
			isInvalid: true,
		},
		{
			name:      "HttpsUrlDomain",
			input:     "https://example.com",
			want:      "",
			isInvalid: true,
		},
		{
			name:      "EmptyDomain",
			input:     "",
			want:      "",
			isInvalid: true,
		},
		{
			name:      "SpaceDomain",
			input:     " ",
			want:      "",
			isInvalid: true,
		},
		{
			name:      "MultiSpaceDomain",
			input:     "        ",
			want:      "",
			isInvalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateFQDN(tt.input)
			// compare output
			require.Equal(t, got, tt.want, tt.name)
			invalid := got != tt.input || tt.input == ""
			require.Equal(t, invalid, tt.isInvalid)
		})

	}
}
