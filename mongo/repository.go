package mongo

import (
	"github.com/pivotal-cf/brokerapi"
	"gopkg.in/mgo.v2/bson"
)

const (
	DatabaseName                         = "brokerData"
	ServiceInstanceCollectionName        = "serviceInstance"
	ServiceInstanceBindingCollectionName = "serviceInstanceBinding"
	ID                                   = "_id"
)

type ServiceInstance struct {
	ServiceInstanceID   string `bson:"_id"`
	ServiceDefinitionID string `bson:"serviceDefinitionID"`
	PlanID              string `bson:"planID"`
	OrganizationGUID    string `bson:"organizationGUID"`
	SpaceGUID           string `bson:"spaceGUID"`
	DashboardUrl        string `bson:"dashboardUrl,omitempty"`
}

type ServiceInstanceUpdate struct {
	ServiceDefinitionID string `bson:"serviceDefinitionID"`
	PlanID              string `bson:"planID"`
}

type ServiceInstanceBinding struct {
	BindingID         string `bson:"_id"`
	ServiceInstanceID string `bson:"serviceInstanceID"`
	SyslogDrainUrl    string `bson:"syslogDrainUrl,omitempty"`
	AppGUID           string `bson:"appGUID"`
}

type Repository struct {
	adminService *AdminService
}

func NewRepository(adminService *AdminService) *Repository {
	repository := &Repository{
		adminService,
	}

	return repository
}

func (repository *Repository) InstanceExists(instanceID string) (bool, error) {
	docExists, error := repository.adminService.DocExists(&bson.M{ID: instanceID}, DatabaseName, ServiceInstanceCollectionName)
	return docExists, error
}

func (repository *Repository) SaveInstance(instanceID string, details brokerapi.ProvisionDetails) error {
	serviceInstance := &ServiceInstance{
		ServiceInstanceID:   instanceID,
		ServiceDefinitionID: details.ServiceID,
		PlanID:              details.PlanID,
		OrganizationGUID:    details.OrganizationGUID,
		SpaceGUID:           details.SpaceGUID,
	}

	error := repository.adminService.SaveDoc(serviceInstance, DatabaseName, ServiceInstanceCollectionName)

	if error != nil {
		return error
	}

	return nil
}

func (repository *Repository) DeleteInstance(instanceID string, details brokerapi.DeprovisionDetails) error {
	error := repository.adminService.RemoveDoc(&bson.M{ID: instanceID}, DatabaseName, ServiceInstanceCollectionName)

	if error != nil {
		return error
	}

	return nil
}

func (repository *Repository) UpdateInstance(instanceID string, details brokerapi.UpdateDetails) error {
	update := &ServiceInstanceUpdate{
		ServiceDefinitionID: details.ServiceID,
		PlanID:              details.PlanID,
	}

	error := repository.adminService.UpdateDoc(&bson.M{ID: instanceID}, update, DatabaseName, ServiceInstanceCollectionName)

	if error != nil {
		return error
	}

	return nil
}

func (repository *Repository) SaveInstanceBinding(instanceID, bindingID string, details brokerapi.BindDetails) error {
	serviceInstanceBinding := &ServiceInstanceBinding{
		BindingID:         bindingID,
		ServiceInstanceID: instanceID,
		AppGUID:           details.AppGUID,
	}

	error := repository.adminService.SaveDoc(serviceInstanceBinding, DatabaseName, ServiceInstanceBindingCollectionName)

	if error != nil {
		return error
	}

	return nil
}

func (repository *Repository) DeleteInstanceBinding(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	error := repository.adminService.RemoveDoc(&bson.M{ID: bindingID}, DatabaseName, ServiceInstanceBindingCollectionName)

	if error != nil {
		return error
	}

	return nil
}

func (repository *Repository) InstanceBindingExists(instanceID, bindingID string) (bool, error) {
	docExists, error := repository.adminService.DocExists(&bson.M{ID: bindingID}, DatabaseName, ServiceInstanceBindingCollectionName)
	return docExists, error
}
