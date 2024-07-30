package crypto

import (
	"testing"

	"github.com/gennadyterekhov/metrics-storage/internal/common/crypto/decrypt"

	"github.com/gennadyterekhov/metrics-storage/internal/common/crypto/encrypt"

	"github.com/stretchr/testify/assert"
)

func TestCanEncryptAndDecrypt(t *testing.T) {
	publicKeyFilePath := "../../../keys/public.test"
	privateKeyFilePath := "../../../keys/private.test"

	plainText := "hello world"

	encrypted := encrypt.TryUsingKeyFileOrReturnPlainText(publicKeyFilePath, []byte(plainText))

	decrypted := decrypt.TryUsingKeyFileOrReturnPlainText(privateKeyFilePath, encrypted)

	assert.Equal(t, plainText, string(decrypted))
}

func TestPlaintextReturnedWithoutFile(t *testing.T) {
	publicKeyFilePath := "../../../keys/wrong_public"
	privateKeyFilePath := "../../../keys/wrong_private"

	plainText := "hello world"

	encrypted := encrypt.TryUsingKeyFileOrReturnPlainText(publicKeyFilePath, []byte(plainText))
	assert.Equal(t, plainText, string(encrypted))

	decrypted := decrypt.TryUsingKeyFileOrReturnPlainText(privateKeyFilePath, encrypted)
	assert.Equal(t, plainText, string(decrypted))
}
