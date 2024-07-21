package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

func isNeededForExampleToShowFullFile2() {}

func ExampleGetMetricHandler() {
	router, controllers := SetUpExampleRouter()
	router.Get("/value/{metricType}/{metricName}", GetMetricHandler(controllers.GetController).ServeHTTP)
	server := httptest.NewServer(router)

	req, err := http.NewRequest("GET", server.URL+"/value/counter/nm", nil)
	require.NoError(nil, err)

	response, err := server.Client().Do(req)
	require.NoError(nil, err)

	readBytes, err := io.ReadAll(response.Body)
	require.NoError(nil, err)

	err = response.Body.Close()
	require.NoError(nil, err)

	fmt.Println(string(readBytes))
}
