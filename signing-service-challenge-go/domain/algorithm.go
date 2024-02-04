package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"strings"
)

type Algorithm uint8

const (
	ECC Algorithm = iota + 1
	RSA
)

var algorithmName = map[uint8]string{
	1: "ECC",
	2: "RSA",
}

var algorithmValue = map[string]uint8{
	"ecc": 1,
	"rsa": 2,
}

// String representation of Algorithm object
func (algorithm Algorithm) String() string {
	return algorithmName[uint8(algorithm)]
}

// ParseAlgorithm from string
func ParseAlgorithm(s string) (Algorithm, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, found := algorithmValue[s]
	if !found {
		return Algorithm(0), fmt.Errorf("%q is not a valid algorithm", s)
	}
	return Algorithm(value), nil
}

// MarshalJSON converts Algorithm to string representation for client
func (algorithm Algorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(algorithm.String())
}

// UnmarshalJSON converts user's string input to Algorithm object
func (algorithm *Algorithm) UnmarshalJSON(data []byte) (err error) {
	var stringValue string
	if err = json.Unmarshal(data, &stringValue); err != nil {
		return err
	}
	if *algorithm, err = ParseAlgorithm(stringValue); err != nil {
		return err
	}
	return nil
}

type KeyPairInBytes struct {
	PrivateKey []byte
	PublicKey  []byte
}

// GenerateKeyPairsInBytes depending on Algorithm type returns KeyPairInBytes for storing
func (algorithm Algorithm) GenerateKeyPairsInBytes() (*KeyPairInBytes, error) {
	switch algorithm {
	case ECC:
		return eccKeyPairInBytesGenerator{
			marshaler: crypto.NewECCMarshaler(),
			generator: &crypto.ECCGenerator{},
		}.generateECCKeyPairInBytes()
	case RSA:
		return rsaKeyPairInBytesGenerator{
			marshaler: crypto.NewRSAMarshaler(),
			generator: &crypto.RSAGenerator{},
		}.generateRSAKeyPairInBytes()
	default:
		return nil, errors.New("invalid algorithm")
	}
}

type eccKeyPairInBytesGenerator struct {
	marshaler crypto.ECCMarshaler
	generator *crypto.ECCGenerator
}

func (keyPairGenerator eccKeyPairInBytesGenerator) generateECCKeyPairInBytes() (*KeyPairInBytes, error) {
	eccKeyPair, err := keyPairGenerator.generator.Generate()
	if err != nil {
		panic("error during ECC key generation")
	}

	publicKey, privateKey, err := keyPairGenerator.marshaler.Encode(*eccKeyPair)
	if err != nil {
		return nil, err
	}
	return &KeyPairInBytes{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

type rsaKeyPairInBytesGenerator struct {
	marshaler crypto.RSAMarshaler
	generator *crypto.RSAGenerator
}

func (keyPairGenerator rsaKeyPairInBytesGenerator) generateRSAKeyPairInBytes() (*KeyPairInBytes, error) {
	rsaKeyPair, err := keyPairGenerator.generator.Generate()
	if err != nil {
		panic("error during RSA key generation")
	}

	publicKey, privateKey, err := keyPairGenerator.marshaler.Marshal(*rsaKeyPair)
	if err != nil {
		return nil, err
	}
	return &KeyPairInBytes{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

func (algorithm Algorithm) Signer(privateKey []byte) (crypto.Signer, error) {
	switch algorithm {
	case ECC:
		marshaler := crypto.NewECCMarshaler()
		return crypto.NewSignerECDSA(privateKey, &marshaler), nil
	case RSA:
		marshaler := crypto.NewRSAMarshaler()
		return crypto.NewSignerRSA(privateKey, &marshaler), nil
	default:
		return nil, errors.New("invalid algorithm")
	}
}
