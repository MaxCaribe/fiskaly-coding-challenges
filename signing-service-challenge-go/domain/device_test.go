package domain

import (
	"testing"
)

type testRepository struct {
	storage map[string]SignatureDevice
}

func (repo *testRepository) Get(uuid string) (SignatureDevice, bool) {
	return SignatureDevice{}, false
}
func (repo *testRepository) GetAll() []SignatureDevice {
	return nil
}
func (repo *testRepository) Create(device SignatureDevice) error {
	repo.storage[device.UUID] = device
	return nil
}
func (repo *testRepository) Update(device SignatureDevice) error {
	return nil
}

func TestCreateSignatureDeviceECC(t *testing.T) {
	repo := testRepository{storage: make(map[string]SignatureDevice)}
	device, err := CreateSignatureDevice("ECC", "", &repo)
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
	device, err := CreateSignatureDevice("RSA", "", &repo)
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
	_, err := CreateSignatureDevice("SHA", "", &repo)
	if err == nil {
		t.Errorf("can't create signature device with invalid algorithm")
	}
}
