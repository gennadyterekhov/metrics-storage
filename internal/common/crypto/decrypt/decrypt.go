package decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/iohelpler"
)

func TryUsingKeyFileOrReturnPlainText(keyFilename string, data []byte) []byte {
	dec, err := usingKeyFile(keyFilename, data)
	if err != nil {
		return data
	}

	return dec
}

func usingKeyFile(keyFilename string, data []byte) ([]byte, error) {
	key, err := getPrivateKeyFromFile(keyFilename)
	if err != nil {
		return nil, err
	}

	return decrypt(key, data)
}

func decrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}

func getPrivateKeyFromFile(keyFilename string) (*rsa.PrivateKey, error) {
	rawKey, err := iohelpler.GetFileContents(keyFilename)
	if err != nil {
		return nil, err
	}

	spkiBlock, _ := pem.Decode(rawKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(spkiBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
