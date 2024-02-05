package domain

import (
	"reflect"
	"testing"
)

type testRepository struct {
	storage map[string]SignatureDevice
}

func (repo *testRepository) Get(uuid string) (SignatureDevice, bool) {
	return repo.storage[uuid], true
}
func (repo *testRepository) GetAll() []SignatureDevice {
	return nil
}
func (repo *testRepository) Create(device SignatureDevice) error {
	repo.storage[device.UUID] = device
	return nil
}
func (repo *testRepository) Update(device SignatureDevice) error {
	repo.storage[device.UUID] = device
	return nil
}

func TestCreateSignatureDeviceECC(t *testing.T) {
	repo := testRepository{storage: make(map[string]SignatureDevice)}
	device, err := CreateSignatureDevice(Algorithm(1), "", &repo)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(repo.storage[device.UUID].UUID) == 0 {
		t.Errorf("returned device is not stored")
	}
	if len(device.PublicKey) == 0 {
		t.Errorf("returned device has no public key")
	}
}

func TestCreateSignatureDeviceRSA(t *testing.T) {
	repo := testRepository{storage: make(map[string]SignatureDevice)}
	device, err := CreateSignatureDevice(Algorithm(2), "", &repo)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(repo.storage[device.UUID].UUID) == 0 {
		t.Errorf("returned device is not stored")
	}
	if len(device.PublicKey) == 0 {
		t.Errorf("returned device has no public key")
	}
}

func TestCreateSignatureDeviceInvalid(t *testing.T) {
	repo := testRepository{storage: make(map[string]SignatureDevice)}
	_, err := CreateSignatureDevice(Algorithm(0), "", &repo)
	if err == nil {
		t.Errorf("can't create signature device with invalid algorithm")
	}
}

func TestSignTransactionRSA(t *testing.T) {
	keyPairInBytes, err := Algorithm(2).GenerateKeyPairsInBytes()
	if err != nil {
		t.Errorf(err.Error())
	}

	device := SignatureDevice{
		UUID:       "uuid",
		Algorithm:  Algorithm(2), //RSA
		PrivateKey: keyPairInBytes.PrivateKey,
		PublicKey:  keyPairInBytes.PublicKey,
	}
	repo := testRepository{storage: make(map[string]SignatureDevice)}
	err = repo.Create(device)

	if err != nil {
		t.Errorf(err.Error())
	}

	dataToSign := "message"
	signedResponse, err := SignTransaction(device.UUID, dataToSign, &repo)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(signedResponse.SignedData) == 0 || len(signedResponse.Signature) == 0 {
		t.Errorf("invalid response")
	}

	savedDevice := repo.storage[device.UUID]
	if savedDevice.SignatureCounter == 0 || reflect.DeepEqual(savedDevice, device) {
		t.Errorf("device wasn't updated")
	}
}

func TestSignTransactionECC(t *testing.T) {
	keyPairInBytes, err := Algorithm(1).GenerateKeyPairsInBytes()
	if err != nil {
		t.Errorf(err.Error())
	}

	device := SignatureDevice{
		UUID:       "uuid",
		Algorithm:  Algorithm(1), //RSA
		PrivateKey: keyPairInBytes.PrivateKey,
		PublicKey:  keyPairInBytes.PublicKey,
	}
	repo := testRepository{storage: make(map[string]SignatureDevice)}
	err = repo.Create(device)

	if err != nil {
		t.Errorf(err.Error())
	}

	dataToSign := "message"
	signedResponse, err := SignTransaction(device.UUID, dataToSign, &repo)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(signedResponse.SignedData) == 0 || len(signedResponse.Signature) == 0 {
		t.Errorf("invalid response")
	}

	savedDevice := repo.storage[device.UUID]
	if savedDevice.SignatureCounter == 0 || reflect.DeepEqual(savedDevice, device) {
		t.Errorf("device wasn't updated")
	}
}
