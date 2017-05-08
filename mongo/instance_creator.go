package mongo

import (
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

type InstanceCreator struct {
	adminService *AdminService
	repository   *Repository
}

func NewInstanceCreator(adminService *AdminService, repository *Repository) *InstanceCreator {
	return &InstanceCreator{
		adminService,
		repository,
	}
}

func (instanceCreator *InstanceCreator) Create(instanceID string, details brokerapi.ProvisionDetails) error {
	// TODO: ATOM
	databaseExists, err := instanceCreator.adminService.DatabaseExists(instanceID)

	if err != nil {
		return err
	}

	// ensure the instance is empty
	if databaseExists {
		err := instanceCreator.adminService.DeleteDatabase(instanceID)

		if err != nil {
			return err
		}
	}

	database, err := instanceCreator.adminService.CreateDatabase(instanceID)

	if err != nil {
		return err
	}

	if database == nil {
		return errors.New("Failed to create new DB instance: " + instanceID)
	}

	err = instanceCreator.repository.SaveInstance(instanceID, details)

	if err != nil {
		return err
	}

	return nil
}

func (instanceCreator *InstanceCreator) Destroy(instanceID string, details brokerapi.DeprovisionDetails) error {
	// TODO: ATOM
	err := instanceCreator.adminService.DeleteDatabase(instanceID)

	if err != nil {
		return err
	}

	err = instanceCreator.repository.DeleteInstance(instanceID, details)

	if err != nil {
		return err
	}

	return nil
}

func (instanceCreator *InstanceCreator) InstanceExists(instanceID string) (bool, error) {
	instanceExists, err := instanceCreator.repository.InstanceExists(instanceID)
	return instanceExists, err
}

func (instanceCreator *InstanceCreator) Update(instanceID string, details brokerapi.UpdateDetails) error {
	err := instanceCreator.repository.UpdateInstance(instanceID, details)

	if err != nil {
		return err
	}

	return nil
}
