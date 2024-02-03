package persistence

import (
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"sync"
)

type DevicesRepository struct {
	storage map[string]domain.SignatureDevice
	mutex   sync.Mutex
}

func NewDevicesRepository() *DevicesRepository {
	repo := DevicesRepository{storage: make(map[string]domain.SignatureDevice)}
	return &repo
}

func (repository *DevicesRepository) Get(uuid string) (domain.SignatureDevice, bool) {
	device, found := repository.storage[uuid]
	return device, found
}

func (repository *DevicesRepository) GetAll() []domain.SignatureDevice {
	devices := make([]domain.SignatureDevice, 0, len(repository.storage))
	for _, device := range repository.storage {
		devices = append(devices, device)
	}
	return devices
}

func (repository *DevicesRepository) Create(device domain.SignatureDevice) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	if _, found := repository.storage[device.UUID]; found {
		return fmt.Errorf(`device with UUID "%q" already exists`, device.UUID)
	}
	repository.storage[device.UUID] = device
	return nil
}

func (repository *DevicesRepository) Update(device domain.SignatureDevice) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	if _, found := repository.storage[device.UUID]; !found {
		return fmt.Errorf(`device with UUID "%q" doesn't exists`, device.UUID)
	}
	repository.storage[device.UUID] = device
	return nil
}
