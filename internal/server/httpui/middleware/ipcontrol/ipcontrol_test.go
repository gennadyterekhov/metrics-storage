package ipcontrol

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test403IfSubnetDoesNotContainIP(t *testing.T) {
	tests := []struct {
		ip            string
		trustedSubnet string
		code          int
	}{
		{
			ip:            "1.1.1.1",
			trustedSubnet: "1.1.1.0/24",
			code:          200,
		},
		{
			ip:            "2.1.1.1",
			trustedSubnet: "1.1.1.0/24",
			code:          403,
		},
		{
			ip:            "",
			trustedSubnet: "1.1.1.0/24",
			code:          403,
		},
		{
			ip:            "1.1.1.1",
			trustedSubnet: "",
			code:          200,
		},
		{
			ip:            "",
			trustedSubnet: "",
			code:          200,
		},
	}

	for i, tt := range tests {

		_, ts, err := net.ParseCIDR(tt.trustedSubnet)
		if tt.trustedSubnet != "" {
			assert.NoError(t, err)
		}

		handler := New(ts).AllowOnlyTrustedSubnet(http.HandlerFunc(okHandler))

		t.Run(fmt.Sprintf("case%v", i), func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			request.Header.Set("X-Real-IP", tt.ip)

			responseWriter := httptest.NewRecorder()
			handler.ServeHTTP(responseWriter, request)
			assert.Equal(t, tt.code, responseWriter.Code)
		})
	}
}

func TestNilSubnet(t *testing.T) {
	handler := New(nil).AllowOnlyTrustedSubnet(http.HandlerFunc(okHandler))
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	request.Header.Set("X-Real-IP", "1.1.1.1")

	responseWriter := httptest.NewRecorder()
	handler.ServeHTTP(responseWriter, request)
	assert.Equal(t, 200, responseWriter.Code)
}

func okHandler(response http.ResponseWriter, _ *http.Request) { response.WriteHeader(200) }
