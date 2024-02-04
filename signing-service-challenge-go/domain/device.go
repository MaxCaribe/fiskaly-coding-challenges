package domain

import (
	"encoding/base64"
	"github.com/google/uuid"
)

type SignatureDevice struct {
	UUID             string
	Label            string
	PrivateKey       []byte
	PublicKey        []byte
	Algorithm        Algorithm
	SignatureCounter int
	LastSignature    []byte
}

type DevicesRepository interface {
	Get(uuid string) (SignatureDevice, bool)
	GetAll() []SignatureDevice
	Create(device SignatureDevice) error
	Update(device SignatureDevice) error
}

type CreateSignatureDeviceResponse struct {
	UUID             string    `json:"uuid"`
	Label            string    `json:"label"`
	PublicKey        []byte    `json:"public_key"`
	Algorithm        Algorithm `json:"algorithm"`
	SignatureCounter int       `json:"signature_counter"`
}

// CreateSignatureDevice creates SignatureDevice in store and returns serializable response
func CreateSignatureDevice(
	algorithm string,
	label string,
	repo DevicesRepository,
) (CreateSignatureDeviceResponse, error) {
	parsedAlgorithm, err := ParseAlgorithm(algorithm)
	if err != nil {
		return CreateSignatureDeviceResponse{}, err
	}

	keyPairInBytes, err := parsedAlgorithm.GenerateKeyPairsInBytes()
	if err != nil {
		return CreateSignatureDeviceResponse{}, err
	}

	id := uuid.NewString()
	lastSignature := base64.URLEncoding.EncodeToString([]byte(id))
	signatureDevice := SignatureDevice{
		UUID:             id,
		Label:            label,
		PrivateKey:       keyPairInBytes.PrivateKey,
		PublicKey:        keyPairInBytes.PublicKey,
		Algorithm:        parsedAlgorithm,
		SignatureCounter: 0,
		LastSignature:    []byte(lastSignature),
	}
	err = repo.Create(signatureDevice)
	if err != nil {
		return CreateSignatureDeviceResponse{}, err
	}

	return CreateSignatureDeviceResponse{
		UUID:             signatureDevice.UUID,
		Label:            signatureDevice.Label,
		PublicKey:        signatureDevice.PublicKey,
		Algorithm:        signatureDevice.Algorithm,
		SignatureCounter: signatureDevice.SignatureCounter,
	}, nil
}
