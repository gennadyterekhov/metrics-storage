package bodypreparer

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCanHashBodyAndSaveInHeader(t *testing.T) {

	type want struct {
		hash string
	}
	tests := []struct {
		name string
		key  string
		want want
	}{
		{
			name: "w key",
			key:  "key",
			want: want{"d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0509"},
		},
		{
			name: "wo key",
			key:  "",
			want: want{""},
		},
	}
	body := `{"id":"nm", "type":"counter", delta":1}`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := PrepareRequest(
				resty.New(),
				[]byte(body),
				false,
				tt.key,
			)
			require.NoError(t, err)

			fmt.Println(request.Header.Get("HashSHA256"))
			assert.Equal(t,
				tt.want.hash,
				request.Header.Get("HashSHA256"),
			)
		})
	}
}
