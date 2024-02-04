package agent

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAgent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Agent()
			// because server is not started
			require.Error(t, err)
		})
	}
}
