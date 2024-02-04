package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

type SignerRSA struct {
	privateKey []byte
	marshaler  *RSAMarshaler
}

func NewSignerRSA(privateKey []byte, marshaler *RSAMarshaler) *SignerRSA {
	return &SignerRSA{
		privateKey,
		marshaler,
	}
}

// Sign implementation for RSA algorithm
func (signer *SignerRSA) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := signer.marshaler.Unmarshal(signer.privateKey)
	if err != nil {
		return nil, err
	}

	hashedData := sha256.Sum256(dataToBeSigned)
	signedData, err := rsa.SignPKCS1v15(nil, keyPair.Private, crypto.SHA256, hashedData[:])
	if err != nil {
		return nil, err
	}

	return signedData, nil
}
