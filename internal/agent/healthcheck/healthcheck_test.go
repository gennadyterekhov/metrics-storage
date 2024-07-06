package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/server/httpui/router"

	"github.com/gennadyterekhov/metrics-storage/internal/common/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestCanSendHeadRequest(t *testing.T) {
	ts := httptest.NewServer(router.GetRouter())
	defer ts.Close()

	type want struct {
		code int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "head",
			url:  "/",
			want: want{code: http.StatusOK},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isOk := MakeHealthcheck(ts.URL)
			assert.Equal(t, true, isOk)

			resp, _ := testhelper.SendRequest(t, ts, http.MethodHead, tt.url)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}
