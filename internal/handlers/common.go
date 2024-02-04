package handlers

import (
	"errors"
	"net/url"
	"strings"
)

func parseURL(url *url.URL) (metricType string, name string, value string, err error) {
	parameters := strings.Split(url.Path, "/")

	if len(parameters) != 5 { // empty string before first slash, update and 3 params
		return "", "", "", errors.New("expected exactly 3 parameters")
	}
	return parameters[2], parameters[3], parameters[4], nil
}
