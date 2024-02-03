package persistence

import (
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"testing"
)

func TestDevicesRepository_GetSuccessful(t *testing.T) {
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	repo := seededRepo([]domain.SignatureDevice{device})

	foundDevice, _ := repo.Get(device.UUID)

	if foundDevice.UUID != device.UUID {
		fmt.Errorf("couldn't retrieve signature device")
	}
}

func TestDevicesRepository_GetNotFound(t *testing.T) {
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	repo := seededRepo([]domain.SignatureDevice{})

	_, found := repo.Get(device.UUID)

	if found {
		fmt.Errorf("device should not present in repository")
	}
}

func TestDevicesRepository_GetAll(t *testing.T) {
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
		fmt.Errorf("incorrect amount of devices returned")
	}
}

func TestDevicesRepository_CreateSuccessful(t *testing.T) {
	repo := seededRepo([]domain.SignatureDevice{})
	device := domain.SignatureDevice{UUID: uuid.NewString()}

	err := repo.Create(device)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	_, found := repo.storage[device.UUID]
	if !found {
		fmt.Errorf("device wasn't saved")
	}
}

func TestDevicesRepository_CreateDuplicateError(t *testing.T) {
	repo := seededRepo([]domain.SignatureDevice{})
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	err := repo.Create(device)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	err = repo.Create(device)
	if err == nil {
		fmt.Errorf("should not be allowed to save multiple devices with same uuid")
	}
}

func TestDevicesRepository_UpdateSuccessful(t *testing.T) {
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	repo := seededRepo([]domain.SignatureDevice{device})
	device.Label = "label"

	err := repo.Update(device)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	updatedDevice := repo.storage[device.UUID]
	if updatedDevice.Label != device.Label || updatedDevice.Label == "" {
		fmt.Errorf("device wasn't saved")
	}
}

func TestDevicesRepository_UpdateNotFound(t *testing.T) {
	repo := seededRepo([]domain.SignatureDevice{})
	device := domain.SignatureDevice{UUID: uuid.NewString()}
	device.Label = "label"

	err := repo.Update(device)
	if err == nil {
		fmt.Errorf("no such device to update")
	}
}

func seededRepo(devices []domain.SignatureDevice) *DevicesRepository {
	repo := NewDevicesRepository()
	for _, device := range devices {
		repo.storage[device.UUID] = device
	}

	return repo
}
