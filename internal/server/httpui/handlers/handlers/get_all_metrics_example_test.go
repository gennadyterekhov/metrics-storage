package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

func isNeededForExampleToShowFullFile() {}

func ExampleGetAllMetricsHandler() {
	router, controllers := SetUpExampleRouter()
	router.Get("/", GetAllMetricsHandler(controllers.GetController).ServeHTTP)
	server := httptest.NewServer(router)

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(nil, err)

	response, err := server.Client().Do(req)
	require.NoError(nil, err)

	readBytes, err := io.ReadAll(response.Body)
	require.NoError(nil, err)

	err = response.Body.Close()
	require.NoError(nil, err)

	fmt.Println(string(readBytes))
	// Output:
	// <!DOCTYPE html>
	// <html>
	//   <head></head>
	//   <body>
	//     <h2>gauge</h2>
	//     <ul>
	//
	//     </ul>
	//     <h2>counter</h2>
	//     <ul>
	//
	//     </ul>
	//   </body>
	// </html>
}
