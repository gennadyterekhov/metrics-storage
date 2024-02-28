package testhelper

import (
	"bytes"
	"compress/gzip"
	"github.com/gennadyterekhov/metrics-storage/internal/constants"
	"github.com/gennadyterekhov/metrics-storage/internal/container"
	"github.com/gennadyterekhov/metrics-storage/internal/helper/iohelpler"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var TestServer *httptest.Server

func Bootstrap(m *testing.M) {
	setUp(nil)
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func BootstrapWithServer(m *testing.M, server *httptest.Server) {
	setUp(server)
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp(server *httptest.Server) {
	container.MetricsRepository.Clear()
	if server != nil {
		TestServer = server
	}
}

func tearDown() {
	if TestServer != nil {
		TestServer.Close()
	}
}

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
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return response, string(respBody)
}

func SendAlreadyJSONedBody(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	requestBody *bytes.Buffer,
) (*http.Response, []byte) {
	//buf := bytes.NewBuffer(requestBody.Bytes())
	req, err := http.NewRequest(method, ts.URL+path, requestBody)
	require.NoError(t, err)
	req.Header.Set(constants.HeaderContentType, constants.ApplicationJSON)

	response, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody := iohelpler.ReadFromReadCloserOrDie(response.Body)

	return response, respBody
}

func SendGzipRequest(
	t *testing.T,
	ts *httptest.Server,
	method,
	path string,
	requestBody string,
) (*http.Response, []byte) {

	buf := bytes.NewBuffer([]byte(requestBody))
	gzipBodyWriter := gzip.NewWriter(buf)
	_, err := gzipBodyWriter.Write([]byte(requestBody))
	require.NoError(t, err)
	err = gzipBodyWriter.Close()
	require.NoError(t, err)

	request := httptest.NewRequest(method, ts.URL+path, buf)

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
	// i dont know why, but here it does not decompress automatically in contrast to compressor package
	respBody := iohelpler.ReadFromGzipReadCloserOrDie(response.Body)

	return response, respBody
}
