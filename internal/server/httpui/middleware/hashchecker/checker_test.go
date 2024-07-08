package hashchecker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test400IfWrongHash(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		body       string
		hashedBody string
		code       int
	}{
		{
			name:       "w key",
			key:        "key",
			body:       `{"id":"nm", "type":"counter", delta":1}`,
			hashedBody: "d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0509",
			code:       200,
		},
		{
			name:       "w key wrong hash",
			key:        "key",
			body:       `{"id":"nm", "type":"counter", delta":1}`,
			hashedBody: "d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0508",
			code:       400,
		},
	}

	for _, tt := range tests {
		handler := New(tt.key).CheckHash(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(200)
			// code reaches this point only if hashes are valid
			assert.Equal(
				t,
				"d6c552977ca4ae14249a3ad8ac77f361265993e49376e8115a27daa4c67e0509",
				response.Header().Get("HashSHA256"),
			)
		}))
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

			responseWriter := httptest.NewRecorder()
			handler.ServeHTTP(responseWriter, request)
			assert.Equal(t, tt.code, responseWriter.Code)
		})
	}
}
