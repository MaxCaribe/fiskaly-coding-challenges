package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"testing"
)

func TestInMemoryDevicesRepository_GetSuccessful(t *testing.T) {
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	repo := seededRepo([]domain.SignatureDevice{device})

	foundDevice, _ := repo.Get(device.UUID)

	if foundDevice.UUID != device.UUID {
		t.Errorf("couldn't retrieve signature device")
	}
}

func TestInMemoryDevicesRepository_GetNotFound(t *testing.T) {
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	repo := seededRepo([]domain.SignatureDevice{})

	_, found := repo.Get(device.UUID)

	if found {
		t.Errorf("device should not present in repository")
	}
}

func TestInMemoryDevicesRepository_GetAll(t *testing.T) {
	devices := []domain.SignatureDevice{
		domain.SignatureDevice{UUID: uuid.NewString()},
		domain.SignatureDevice{UUID: uuid.NewString()},
		domain.SignatureDevice{UUID: uuid.NewString()},
		domain.SignatureDevice{UUID: uuid.NewString()},
		domain.SignatureDevice{UUID: uuid.NewString()},
	}

	repo := seededRepo(devices)
	foundDevices := repo.GetAll()
	if len(foundDevices) != len(devices) || len(foundDevices) == 0 {
		t.Errorf("incorrect amount of devices returned")
	}
}

func TestInMemoryDevicesRepository_CreateSuccessful(t *testing.T) {
	repo := seededRepo([]domain.SignatureDevice{})
	device := domain.SignatureDevice{UUID: uuid.NewString()}

	err := repo.Create(device)
	if err != nil {
		t.Errorf(err.Error())
	}
	_, found := repo.storage[device.UUID]
	if !found {
		t.Errorf("device wasn't saved")
	}
}

func TestInMemoryDevicesRepository_CreateDuplicateError(t *testing.T) {
	repo := seededRepo([]domain.SignatureDevice{})
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	err := repo.Create(device)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = repo.Create(device)
	if err == nil {
		t.Errorf("should not be allowed to save multiple devices with same uuid")
	}
}

func TestInMemoryDevicesRepository_UpdateSuccessful(t *testing.T) {
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	repo := seededRepo([]domain.SignatureDevice{device})
	device.Label = "label"

	err := repo.Update(device)
	if err != nil {
		t.Errorf(err.Error())
	}

	updatedDevice := repo.storage[device.UUID]
	if updatedDevice.Label != device.Label || updatedDevice.Label == "" {
		t.Errorf("device wasn't saved")
	}
}

func TestInMemoryDevicesRepository_UpdateNotFound(t *testing.T) {
	repo := seededRepo([]domain.SignatureDevice{})
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	device.Label = "label"

	err := repo.Update(device)
	if err == nil {
		t.Errorf("no such device to update")
	}
}

func seededRepo(devices []domain.SignatureDevice) *InMemoryDevicesRepository {
	repo := NewInMemoryDevicesRepository()
	for _, device := range devices {
		repo.storage[device.UUID] = device
	}

	return repo
}
