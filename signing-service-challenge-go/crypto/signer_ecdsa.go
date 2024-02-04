package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
)

type SignerECDSA struct {
	privateKey []byte
	marshaler  *ECCMarshaler
}

func NewSignerECDSA(privateKey []byte, marshaler *ECCMarshaler) *SignerECDSA {
	return &SignerECDSA{
		privateKey,
		marshaler,
	}
}

// Sign implementation for ECC algorithm
func (signer *SignerECDSA) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := signer.marshaler.Decode(signer.privateKey)
	if err != nil {
		return nil, err
	}

	hashedData := sha256.Sum256(dataToBeSigned)
	signedData, err := ecdsa.SignASN1(nil, keyPair.Private, hashedData[:])
	if err != nil {
		return nil, err
	}

	return signedData, nil
}
