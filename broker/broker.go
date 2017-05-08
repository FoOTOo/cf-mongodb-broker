package broker

import (
	"context"
	"errors"

	"github.com/FoOTOo/cf-mongodb-broker/config"
	"github.com/pivotal-cf/brokerapi"
	"gopkg.in/mgo.v2/bson"
)

type InstanceCredentials struct {
	Host     string
	Port     int
	Password string
}

type InstanceCreator interface {
	Create(instanceID string, serviceDetails brokerapi.ProvisionDetails) error
	Destroy(instanceID string, details brokerapi.DeprovisionDetails) error
	Update(instanceID string, details brokerapi.UpdateDetails) error
	InstanceExists(instanceID string) (bool, error)
}

type InstanceBinder interface {
	Bind(instanceID string, bindingID string, details brokerapi.BindDetails) (bson.M, error)
	Unbind(instanceID string, bindingID string, details brokerapi.UnbindDetails) error
	InstanceBindingExists(instanceID, bindingID string) (bool, error)
}

type MongoServiceBroker struct {
	InstanceCreators map[string]InstanceCreator
	InstanceBinders  map[string]InstanceBinder
	Config           config.Config
}

func (mongoServiceBroker *MongoServiceBroker) Services(context context.Context) []brokerapi.Service {
	return mongoServiceBroker.Config.Services()
}

func (mongoServiceBroker *MongoServiceBroker) Provision(context context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	spec := brokerapi.ProvisionedServiceSpec{}

	if details.PlanID == "" {
		return spec, errors.New("plan_id required")
	}

	planIdentifier := ""
	for key, plan := range mongoServiceBroker.plans() {
		if plan.ID == details.PlanID {
			planIdentifier = key
			break
		}
	}

	if planIdentifier == "" {
		return spec, errors.New("plan_id not recognized")
	}

	instanceCreator, ok := mongoServiceBroker.InstanceCreators[planIdentifier]
	if !ok {
		return spec, errors.New("instance creator not found for plan")
	}

	err := instanceCreator.Create(instanceID, details)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (mongoServiceBroker *MongoServiceBroker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	for _, instanceCreator := range mongoServiceBroker.InstanceCreators {
		instanceExists, err := instanceCreator.InstanceExists(instanceID)

		if err != nil {
			return spec, err
		}

		if instanceExists {
			return spec, instanceCreator.Destroy(instanceID, details)
		}
	}

	return spec, brokerapi.ErrInstanceDoesNotExist
}

func (mongoServiceBroker *MongoServiceBroker) Update(context context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	spec := brokerapi.UpdateServiceSpec{}

	if details.PlanID == "" {
		return spec, errors.New("plan_id required")
	}

	planIdentifier := ""
	for key, plan := range mongoServiceBroker.plans() {
		if plan.ID == details.PlanID {
			planIdentifier = key
			break
		}
	}

	if planIdentifier == "" {
		return spec, errors.New("plan_id not recognized")
	}

	instanceCreator, ok := mongoServiceBroker.InstanceCreators[planIdentifier]
	if !ok {
		return spec, errors.New("instance creator not found for plan")
	}

	err := instanceCreator.Update(instanceID, details)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (mongoServiceBroker *MongoServiceBroker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	for key, instanceCreator := range mongoServiceBroker.InstanceCreators {
		instanceExists, err := instanceCreator.InstanceExists(instanceID)

		if err != nil {
			return binding, err
		}

		if instanceExists {
			instanceBinder, ok := mongoServiceBroker.InstanceBinders[key]
			if !ok {
				return binding, errors.New("instance binder not found for plan")
			}

			credentials, err := instanceBinder.Bind(instanceID, bindingID, details)
			binding.Credentials = credentials
			return binding, err
		}
	}

	return binding, brokerapi.ErrInstanceDoesNotExist
}

func (mongoServiceBroker *MongoServiceBroker) Unbind(context context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	for _, instanceBinder := range mongoServiceBroker.InstanceBinders {
		instanceExists, err := instanceBinder.InstanceBindingExists(instanceID, bindingID)

		if err != nil {
			return err
		}

		if instanceExists {
			err = instanceBinder.Unbind(instanceID, bindingID, details)
			return err
		}
	}

	return brokerapi.ErrInstanceDoesNotExist
}

func (mongoServiceBroker *MongoServiceBroker) plans() map[string]*brokerapi.ServicePlan {
	return mongoServiceBroker.Config.Plans()
}

//func (mongoServiceBroker *MongoServiceBroker) instanceExists(instanceID string) bool {
//	for _, instanceCreator := range mongoServiceBroker.InstanceCreators {
//		instanceExists, _ := instanceCreator.InstanceExists(instanceID)
//		if instanceExists {
//			return true
//		}
//	}
//	return false
//}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (mongoServiceBroker *MongoServiceBroker) LastOperation(context context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
