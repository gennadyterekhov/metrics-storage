package compressor

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type args struct {
	accepts          []string
	acceptEncodings  []string
	contentEncodings []string
	contenttype      string
}

func TestIsGzipAvailableForThisRequest(t *testing.T) {
	// in tests this wants gzip
	//        Accept: application/json
	//        Accept-Encoding: gzip
	//        Content-Type: application/json
	// and sometimes
	// Accept: html/text

	tests := []struct {
		name string
		args args
		ok   bool
	}{
		{
			name: "only json",
			args: args{contenttype: "json"},
			ok:   false,
		},
		{
			name: "only html/text",
			args: args{accepts: []string{"asdf", "html/text", "asdf"}},
			ok:   false,
		},
		{
			name: " html/text gzip",
			args: args{accepts: []string{"asdf", "html/text", "asdf"},
				acceptEncodings: []string{"asdf", "gzip", "asdf"}},
			ok: true,
		},
		{
			name: " text/html gzip",
			args: args{accepts: []string{"asdf", "text/html", "asdf"},
				acceptEncodings: []string{"asdf", "gzip", "asdf"}},
			ok: true,
		},
		{
			name: " json gzip",
			args: args{accepts: []string{"asdf", "application/json", "asdf"},
				acceptEncodings: []string{"asdf", "gzip", "asdf"}},
			ok: true,
		},
		{
			name: "json gzip gzip",
			args: args{accepts: []string{"asdf", "application/json", "asdf"},
				acceptEncodings:  []string{"asdf", "gzip", "asdf"},
				contentEncodings: []string{"asdf", "gzip", "asdf"}},
			ok: true,
		},
		{
			name: "only contentEncodings",
			args: args{contentEncodings: []string{"asdf", "gzip", "asdf"}},
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			addHeaders(request, tt.args)

			ok := IsGzipAvailableForThisRequest(request)
			assert.Equal(t, tt.ok, ok)

		})
	}
}

func addHeaders(request *http.Request, args args) {
	for _, val := range args.accepts {
		request.Header.Set("Accept", val)
	}
	for _, val := range args.acceptEncodings {
		request.Header.Set("Accept-Encoding", val)
	}
	for _, val := range args.contentEncodings {
		request.Header.Set("Content-Encoding", val)
	}
	if args.contenttype != "" {
		request.Header.Set("Content-Type", args.contenttype)
	}
}
