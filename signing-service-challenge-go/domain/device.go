package domain

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"strings"
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

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
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

// SignTransaction signs data with found devices, updates device's data and returns signed data
func SignTransaction(id string, data string, repo DevicesRepository) (SignatureResponse, error) {
	device, found := repo.Get(id)
	if !found {
		return SignatureResponse{}, fmt.Errorf("could not found signature device with id %q", id)
	}

	signer, err := device.Algorithm.Signer(device.PrivateKey)
	if err != nil {
		return SignatureResponse{}, err
	}

	securedDataToBeSigned := buildSecuredDataToBeSigned(device.SignatureCounter, data, device.LastSignature)
	signedData, err := signer.Sign([]byte(securedDataToBeSigned))
	if err != nil {
		return SignatureResponse{}, err
	}

	device.SignatureCounter += 1
	device.LastSignature = signedData
	err = repo.Update(device)
	if err != nil {
		return SignatureResponse{}, err
	}

	signedDataBase64 := base64.URLEncoding.EncodeToString(signedData)
	return SignatureResponse{
		Signature:  signedDataBase64,
		SignedData: string(signedData),
	}, nil
}

func buildSecuredDataToBeSigned(signatureCounter int, data string, lastSignature []byte) string {
	var sb strings.Builder
	sb.WriteString(string(rune(signatureCounter)))
	sb.WriteString("_")
	sb.WriteString(data)
	sb.WriteString("_")
	sb.WriteString(string(lastSignature))
	return sb.String()
}
