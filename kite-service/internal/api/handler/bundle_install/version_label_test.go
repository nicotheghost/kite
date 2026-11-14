package bundleinstall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextTemplateVersionLabel(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		revision                  int
		previous, requested, want string
		invalid                   bool
	}{
		{name: "legacy default", revision: 1},
		{name: "first custom label", revision: 1, requested: "1.01", want: "1.01"},
		{name: "second major", revision: 2, requested: "2.5", want: "2.5"},
		{name: "retains formatting", revision: 2, previous: "1.01", requested: "1.020", want: "1.020"},
		{name: "defaults from displayed version", revision: 4, previous: "2.5", want: "3"},
		{name: "compares decimals numerically", previous: "1.9", requested: "1.10", invalid: true},
		{name: "rejects equivalent number", previous: "2.5", requested: "2.50", invalid: true},
		{name: "rejects earlier major", revision: 2, requested: "1.99", invalid: true},
		{name: "rejects invalid format", revision: 1, requested: "1.0.1", invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := nextTemplateVersionLabel(tc.revision, tc.previous, tc.requested)
			if tc.invalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
