package decryptor

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gennadyterekhov/metrics-storage/internal/common/crypto/decrypt"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

type Decryptor struct {
	PrivateKeyFilePath string
}

func New(privateKeyFilePath string) *Decryptor {
	return &Decryptor{
		PrivateKeyFilePath: privateKeyFilePath,
	}
}

func (dc *Decryptor) TryToDecryptUsingPrivateKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var reqBody []byte
		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Custom.Errorln("could not read body", err.Error())
			next.ServeHTTP(res, req)
			return
		}

		decryptedBody := decrypt.TryUsingKeyFileOrReturnPlainText(dc.PrivateKeyFilePath, reqBody)
		req.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

		next.ServeHTTP(res, req)
	})
}
