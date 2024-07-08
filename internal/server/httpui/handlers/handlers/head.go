package handlers

import (
	"io"
	"net/http"
)

func HeadHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	_, err := io.WriteString(res, "ok")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
