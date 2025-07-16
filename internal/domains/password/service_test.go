package password_test

import (
	"testing"

	"github.com/hd-passgen/core/internal/domains/password"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestService_Generate(t *testing.T) {
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

			s := password.NewService()

			results := make([]string, 0, len(tt.serviceNames))
			for _, serviceName := range tt.serviceNames {
				result, err := s.Generate(tt.password, serviceName)
				require.NoError(t, err)

				results = append(results, result)
			}

			if tt.uniqueResults {
				require.Len(t, lo.Uniq(results), len(results))
			} else {
				require.Len(t, lo.Uniq(results), 1)
			}
		})
	}
}
