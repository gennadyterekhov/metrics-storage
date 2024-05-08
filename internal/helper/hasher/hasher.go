package hasher

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/logger"
)

func HashBytes(target []byte, key string) []byte {
	keyBytes := []byte(key)

	sig := hmac.New(sha256.New, keyBytes)
	sig.Write(target)

	return sig.Sum(nil)
}

func IsBodyHashValid(req *http.Request, key string) bool {
	if key == "" {
		logger.ZapSugarLogger.Debugln("hash key is empty -> hash checking is successful")

		return true
	}
	logger.ZapSugarLogger.Debugln("hash key is not empty -> checking hash")

	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			logger.ZapSugarLogger.Errorln("could not read body to check hash", err.Error())
			return false
		}
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return isBytesHashValid(bodyBytes, req.Header.Get("HashSHA256"), key)
}

func isBytesHashValid(body []byte, hash string, key string) bool {
	keyBytes := []byte(key)

	sig := hmac.New(sha256.New, keyBytes)
	sig.Write(body)

	bodyHash := sig.Sum(nil)

	logger.ZapSugarLogger.Debugln(
		"hash of body:",
		hex.EncodeToString(bodyHash),
		"hash sent by client in the header:",
		hash,
	)

	return hex.EncodeToString(bodyHash) == hash
}
