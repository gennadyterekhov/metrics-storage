package hasher

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCheckHash(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		body       string
		hashedBody string
		ok         bool
	}{
		{
			name:       "w key",
			key:        "key",
			body:       `{"id":"nm", "type":"counter", delta":1}`,
			hashedBody: "d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0509",
			ok:         true,
		},
		{
			name:       "w key wrong hash",
			key:        "key",
			body:       `{"id":"nm", "type":"counter", delta":1}`,
			hashedBody: "d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0508",
			ok:         false,
		},
		{
			name:       "w key wrong body",
			key:        "key",
			body:       `{"id":"nm", "type":"counter", delta":2}`,
			hashedBody: "d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0509",
			ok:         false,
		},
		{
			name:       "wo key",
			key:        "",
			body:       `{"id":"nm", "type":"counter", delta":1}`,
			hashedBody: "d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0509",
			ok:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader bytes.Buffer
			_, err := bodyReader.WriteString(tt.body)
			assert.NoError(t, err)

			request := httptest.NewRequest(
				http.MethodPost,
				"http://localhost:8080/",
				&bodyReader,
			)
			request.Header.Set("HashSHA256", tt.hashedBody)

			ok, err := IsBodyHashValid(request, tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.ok, ok)
		})
	}
}
