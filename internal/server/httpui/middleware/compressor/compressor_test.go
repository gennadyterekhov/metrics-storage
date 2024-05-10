package compressor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type args struct {
	accepts          []string
	acceptEncodings  []string
	contentEncodings []string
	contenttype      string
}

func TestIsGzipAvailableForThisRequest(t *testing.T) {
	tests := []struct {
		name string
		args args
		ok   bool
	}{
		{
			name: "only json",
			args: args{contenttype: "application/json"},
			ok:   false,
		},
		{
			name: "only html/text",
			args: args{accepts: []string{"asdf", "html/text", "asdf"}},
			ok:   false,
		},
		{
			name: " html/text gzip",
			args: args{
				accepts:         []string{"asdf", "html/text", "asdf"},
				acceptEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
		{
			name: " text/html gzip",
			args: args{
				accepts:         []string{"asdf", "text/html", "asdf"},
				acceptEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
		{
			name: " json gzip",
			args: args{
				accepts:         []string{"asdf", "application/json", "asdf"},
				acceptEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
		{
			name: "json gzip gzip",
			args: args{
				accepts:          []string{"asdf", "application/json", "asdf"},
				acceptEncodings:  []string{"asdf", "gzip", "asdf"},
				contentEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
		{
			name: "only contentEncodings",
			args: args{contentEncodings: []string{"asdf", "gzip", "asdf"}},
			ok:   false,
		},
		{
			name: " application/json accept gzip",
			args: args{
				accepts: []string{"asdf", "application/json", "asdf"},
				// contenttype:     "application/json",// not required
				acceptEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
		{
			name: " text/html accept gzip",
			args: args{
				accepts: []string{"asdf", "text/html", "asdf"},
				// contenttype:     "text/html", // not required
				acceptEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
		{
			name: " html/text accept gzip",
			args: args{
				accepts: []string{"asdf", "html/text", "asdf"},
				// contenttype:     "html/text",// not required
				acceptEncodings: []string{"asdf", "gzip", "asdf"},
			},
			ok: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
			addHeaders(request, tt.args)

			ok := isGzipAvailableForThisRequest(request)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func addHeaders(request *http.Request, args args) {
	request.Header.Set("Accept", strings.Join(args.accepts, ","))
	request.Header.Set("Accept-Encoding", strings.Join(args.acceptEncodings, ","))
	request.Header.Set("Content-Encoding", strings.Join(args.contentEncodings, ","))

	if args.contenttype != "" {
		request.Header.Set("Content-Type", args.contenttype)
	}
}
