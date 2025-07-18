package password

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func Test_Generate(t *testing.T) {
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
				result, err := Generate(tt.password, serviceName, 0)
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

func Test_GenerateWithLenght(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name        string
		master      string
		lenght      uint8
		expErr      error
		serviceName string
	}{
		{
			name:        "short password",
			master:      "qwerty",
			serviceName: "github.com",
			lenght:      2,
			expErr:      ErrInvalidLength,
		},
		{
			name:        "default password",
			master:      "qwerty",
			serviceName: "github.com",
			lenght:      32,
		},
		{
			name:        "max available password",
			master:      "qwerty",
			serviceName: "github.com",
			lenght:      40,
		},
		{
			name:        "too long password",
			master:      "qwerty",
			serviceName: "github.com",
			lenght:      41,
			expErr:      ErrInvalidLength,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Generate(tt.master, tt.serviceName, tt.lenght)
			if err != nil {
				require.ErrorIs(t, err, tt.expErr)
			} else {
				require.Len(t, result, int(tt.lenght))
			}

		})
	}
}
