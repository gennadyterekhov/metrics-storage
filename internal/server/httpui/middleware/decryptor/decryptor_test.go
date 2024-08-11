package decryptor

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/crypto/encrypt"

	"github.com/stretchr/testify/assert"
)

func TestCanDecrypt(t *testing.T) {
	publicKeyFilePath := "../../../keys/public.test"
	privateKeyFilePath := "../../../keys/private.test"
	body := `{"id":"nm","type":"counter","delta":1,"value":0}`
	encrypted := encrypt.TryUsingKeyFileOrReturnPlainText(publicKeyFilePath, []byte(body))
	assertCanDecrypt(t, privateKeyFilePath, encrypted)
}

func TestReturnsPlaintextIfNoPrivateKeyFile(t *testing.T) {
	privateKeyFilePath := "../../../keys/wrong_private"
	body := `{"id":"nm","type":"counter","delta":1,"value":0}`

	assertCanDecrypt(t, privateKeyFilePath, []byte(body))
}

func assertCanDecrypt(t *testing.T, privateKeyFilePath string, body []byte) {
	handler := New(privateKeyFilePath).TryToDecryptUsingPrivateKey(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			decrypted, err := io.ReadAll(req.Body)
			assert.NoError(t, err)

			assert.Equal(
				t,
				body,
				decrypted,
			)

			res.WriteHeader(200)
		}))

	var bodyReader bytes.Buffer
	_, err := bodyReader.Write(body)
	assert.NoError(t, err)

	request := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		&bodyReader,
	)

	responseWriter := httptest.NewRecorder()
	handler.ServeHTTP(responseWriter, request)

	assert.Equal(t, http.StatusOK, responseWriter.Code)
}
