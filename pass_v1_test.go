package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name          string
		password      string
		serviceNames  []string
		uniqueResults bool
	}{
		{
			name:          "simple password",
			password:      "qwerty",
			serviceNames:  []string{"github.com", "gitlab.com", "google.com"},
			uniqueResults: true,
		},
		{
			name:          "the same services",
			password:      "qwerty",
			serviceNames:  []string{"github.com", "github.com"},
			uniqueResults: false,
		},
		{
			name:          "many the same services",
			password:      "qwerty",
			serviceNames:  []string{"github.com", "github.com", "github.com", "github.com", "github.com", "github.com", "github.com", "github.com", "github.com", "github.com"},
			uniqueResults: false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			results := make([]string, 0, len(tt.serviceNames))
			for _, serviceName := range tt.serviceNames {
				result, err := generatePassword(parameters{
					ServiceName:    serviceName,
					MasterPassword: tt.password,
				})
				require.NoError(t, err)

				results = append(results, result)
			}

			if tt.uniqueResults {
				require.Len(t, unique(results), len(results))
			} else {
				require.Len(t, unique(results), 1)
			}
		})
	}
}

func unique(collection []string) []string {
	result := make([]string, 0, len(collection))
	seen := make(map[string]struct{}, len(collection))

	for i := range collection {
		if _, ok := seen[collection[i]]; ok {
			continue
		}

		seen[collection[i]] = struct{}{}
		result = append(result, collection[i])
	}

	return result

}
