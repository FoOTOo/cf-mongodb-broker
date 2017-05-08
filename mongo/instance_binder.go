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

func (instanceBinder *InstanceBinder) Bind(instanceID string, bindingID string, details brokerapi.BindDetails) (credentials bson.M, err error) {
	// TODO: ATOM
	credentials = bson.M{}

	databaseName := instanceID
	username := bindingID
	password := utiils.GenerateRandomString(25)

	// TODO check if user already exists in the DB

	err = instanceBinder.adminService.CreateUser(databaseName, username, password)

	if err != nil {
		return credentials, err
	}

	err = instanceBinder.repository.SaveInstanceBinding(instanceID, bindingID, details)

	if err != nil {
		return credentials, err
	}

	credentials["uri"] = instanceBinder.adminService.GetConnectionString(databaseName, username, password)

	return credentials, nil
}

func (instanceBinder *InstanceBinder) Unbind(instanceID string, bindingID string, details brokerapi.UnbindDetails) error {
	// TODO: ATOM
	databaseName := instanceID
	username := bindingID
	err := instanceBinder.adminService.DeleteUser(databaseName, username)

	if err != nil {
		return err
	}

	err = instanceBinder.repository.DeleteInstanceBinding(instanceID, bindingID, details)

	if err != nil {
		return err
	}

	return nil
}

func (instanceBinder *InstanceBinder) InstanceBindingExists(instanceID string, bindingID string) (bool, error) {
	instanceBindingExists, err := instanceBinder.repository.InstanceBindingExists(instanceID, bindingID)
	return instanceBindingExists, err
}
