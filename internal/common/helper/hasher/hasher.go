package hasher

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

func HashBytes(target []byte, key string) ([]byte, error) {
	keyBytes := []byte(key)

	sig := hmac.New(sha256.New, keyBytes)
	_, err := sig.Write(target)
	if err != nil {
		return nil, err
	}
	return sig.Sum(nil), nil
}

func IsBodyHashValid(req *http.Request, key string) (bool, error) {
	if key == "" {
		logger.Custom.Debugln("hash key is empty -> hash checking is successful")

		return true, nil
	}
	logger.Custom.Debugln("hash key is not empty -> checking hash")

	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			logger.Custom.Errorln("could not read body to check hash", err.Error())
			return false, nil
		}
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return isBytesHashValid(bodyBytes, req.Header.Get("HashSHA256"), key)
}

func isBytesHashValid(body []byte, hash string, key string) (bool, error) {
	keyBytes := []byte(key)

	sig := hmac.New(sha256.New, keyBytes)
	_, err := sig.Write(body)
	if err != nil {
		return false, err
	}

	bodyHash := sig.Sum(nil)

	logger.Custom.Debugln(
		"hash of body:",
		hex.EncodeToString(bodyHash),
		"hash sent by client in the header:",
		hash,
	)

	return hex.EncodeToString(bodyHash) == hash, nil
}
