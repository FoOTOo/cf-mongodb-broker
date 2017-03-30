package mongo

import (
	"github.com/FoOTOo/cf-mongodb-broker/utils"
	"github.com/pivotal-cf/brokerapi"
	"gopkg.in/mgo.v2/bson"
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

func (instanceBinder *InstanceBinder) Bind(instanceID string, bindingID string, details brokerapi.BindDetails) (credentials bson.M, error error) {
	// TODO: ATOM
	credentials = bson.M{}

	databaseName := instanceID
	username := bindingID
	password := utiils.GenerateRandomString(25)

	// TODO check if user already exists in the DB

	error = instanceBinder.adminService.CreateUser(databaseName, username, password)

	if error != nil {
		return credentials, error
	}

	error = instanceBinder.repository.SaveInstanceBinding(instanceID, bindingID, details)

	if error != nil {
		return credentials, error
	}

	credentials["uri"] = instanceBinder.adminService.GetConnectionString(databaseName, username, password)

	return credentials, nil
}

func (instanceBinder *InstanceBinder) Unbind(instanceID string, bindingID string, details brokerapi.UnbindDetails) error {
	// TODO: ATOM
	databaseName := instanceID
	username := bindingID
	error := instanceBinder.adminService.DeleteUser(databaseName, username)

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
