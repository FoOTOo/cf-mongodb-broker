package mongo

import (
	"errors"

	"github.com/FoOTOo/mongo-service-broker-golang-ultragtx/utils"
	"github.com/pivotal-cf/brokerapi"
)

type InstanceBinder struct {
	adminService *AdminService
	repository   *Repository
}

func NewInstanceBinder(adminService *AdminService, repository *Repository) *InstanceBinder {
	return &InstanceBinder{
		adminService,
		repository,
	}
}

func (instanceBinder *InstanceBinder) Bind(instanceID string, bindingID string, details brokerapi.BindDetails) error {
	// TODO: ATOM
	instanceBindingExists, error := instanceBinder.repository.InstanceBindingExists(instanceID, bindingID)

	if error != nil {
		return error
	}

	if instanceBindingExists {
		return instanceBindingExistsError(instanceID, bindingID, details)
	}

	// TODO check if user already exists in the DB

	databaseName := instanceID
	username := bindingID
	password := utiils.GenerateRandomString(25)

	error = instanceBinder.adminService.CreateUser(databaseName, username, password)

	if error != nil {
		return error
	}

	error = instanceBinder.repository.SaveInstanceBinding(instanceID, bindingID, details)

	if error != nil {
		return error
	}

	return nil
}

func (instanceBinder *InstanceBinder) Unbind(instanceID string, bindingID string, details brokerapi.UnbindDetails) error {
	// TODO: ATOM
	instanceBindingExists, error := instanceBinder.repository.InstanceBindingExists(instanceID, bindingID)

	if error != nil {
		return error
	}

	if !instanceBindingExists {
		return instanceBindingDoesNotExistError(instanceID, bindingID, details)
	}

	databaseName := instanceID
	username := bindingID
	error = instanceBinder.adminService.DeleteUser(databaseName, username)

	if error != nil {
		return error
	}

	error = instanceBinder.repository.DeleteInstanceBinding(instanceID, bindingID, details)

	if error != nil {
		return error
	}

	return nil
}

func (instanceBinder *InstanceBinder) InstanceBindingExists(instanceID string, bindingID string) (bool, error) {
	instanceBindingExists, error := instanceBinder.repository.InstanceBindingExists(instanceID, bindingID)
	return instanceBindingExists, error
}

func instanceBindingExistsError(instanceID, bindingID string, details brokerapi.BindDetails) error {
	return errors.New("Instance binding exists, incetanceID: " + instanceID + ", bindingID: " + bindingID)
}

func instanceBindingDoesNotExistError(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return errors.New("Instance binding doesn't exist, incetanceID: " + instanceID + ", bindingID: " + bindingID)
}
