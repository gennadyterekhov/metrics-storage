package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/iohelpler"
)

func TryUsingKeyFileOrReturnPlainText(keyFilename string, data []byte) []byte {
	enc, err := usingKeyFile(keyFilename, data)
	if err != nil {
		return data
	}
	return enc
}

func usingKeyFile(keyFilename string, data []byte) ([]byte, error) {
	key, err := getPublicKeyFromFile(keyFilename)
	if err != nil {
		return nil, err
	}

	return encrypt(key, data)
}

func encrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return nil, err
	}

	return encryptBytes, nil
}

func getPublicKeyFromFile(keyFilename string) (*rsa.PublicKey, error) {
	rawKey, err := iohelpler.GetFileContents(keyFilename)
	if err != nil {
		return nil, err
	}

	spkiBlock, _ := pem.Decode(rawKey)
	publicKey, err := x509.ParsePKCS1PublicKey(spkiBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
