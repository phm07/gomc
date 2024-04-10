package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
)

var (
	PublicKeyBytes []byte
	privateKey     *rsa.PrivateKey
)

type publicKeyInfo struct {
	Algorithm algorithmIdentifier
	PublicKey asn1.BitString
}

type algorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue
}

func GenerateKeypair() error {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	pubBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	PublicKeyBytes, err = asn1.Marshal(publicKeyInfo{
		Algorithm: algorithmIdentifier{
			Algorithm:  asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1},
			Parameters: asn1.RawValue{Tag: asn1.TagNull},
		},
		PublicKey: asn1.BitString{
			Bytes:     pubBytes,
			BitLength: len(pubBytes) * 8,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func Decrypt(encrypted []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
}
