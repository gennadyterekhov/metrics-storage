package cliargs

import (
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host           string
	Port           int
	wasInitialized bool
}

func (na *NetAddress) String() string {
	if !na.wasInitialized {
		return "http://localhost:8080"
	}
	return fmt.Sprintf("%v:%v", na.Host, na.Port)
}

func (na *NetAddress) Set(flagValue string) error {
	hostAndPort := strings.Split(flagValue, ":")

	if len(hostAndPort) != 2 {
		return fmt.Errorf("invalid format")
	}
	parsedPort, err := strconv.Atoi(hostAndPort[1])
	if err != nil {
		return err
	}
	na.wasInitialized = true
	na.Host = hostAndPort[0]
	na.Port = parsedPort
	return nil
}
