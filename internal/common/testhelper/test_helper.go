package testhelper

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/iohelpler"
	"github.com/stretchr/testify/require"
)

func SendRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	response, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody := iohelpler.ReadFromReadCloserOrDie(response.Body)
	err = response.Body.Close()
	if err != nil {
		return nil, ""
	}
	return response, string(respBody)
}

func SendAlreadyJSONedBody(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	requestBody *bytes.Buffer,
) (*http.Response, []byte) {
	// buf := bytes.NewBuffer(requestBody.Bytes())
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	require.NoError(t, err)
	req.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)

	response, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody := iohelpler.ReadFromReadCloserOrDie(response.Body)
	err = response.Body.Close()
	if err != nil {
		return nil, nil
	}

	return response, respBody
}

func SendGzipRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	requestBody string,
) (*http.Response, []byte) {
	var buf bytes.Buffer

	gzipBodyWriter := gzip.NewWriter(&buf)
	_, err := gzipBodyWriter.Write([]byte(requestBody))
	require.NoError(t, err)
	err = gzipBodyWriter.Close()
	require.NoError(t, err)

	request := httptest.NewRequest(method, ts.URL+path, &buf)

	request.RequestURI = ""
	u, err := url.Parse(ts.URL + path)
	require.NoError(t, err)
	request.URL = u
	request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
	request.Header.Set("Accept", constants.ApplicationJSON)
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")

	response, err := ts.Client().Do(request)

	require.NoError(t, err)
	// I don't know why, but here it does not decompress automatically in contrast to compressor package
	respBody := iohelpler.ReadFromGzipReadCloserOrDie(response.Body)
	err = response.Body.Close()
	if err != nil {
		return nil, nil
	}

	return response, respBody
}

func SendGzipNoBodyRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
) (*http.Response, []byte) {
	request := httptest.NewRequest(method, ts.URL+path, nil)

	request.RequestURI = ""
	u, err := url.Parse(ts.URL + path)
	require.NoError(t, err)
	request.URL = u
	request.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)
	request.Header.Set("Accept", "html/text")
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")

	response, err := ts.Client().Do(request)

	require.NoError(t, err)
	// I don't know why, but here it does not decompress automatically in contrast to compressor package
	respBody := iohelpler.ReadFromGzipReadCloserOrDie(response.Body)
	err = response.Body.Close()
	if err != nil {
		return nil, nil
	}

	return response, respBody
}
