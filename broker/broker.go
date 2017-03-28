package broker

import (
	"context"
	"errors"

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
}

func (mongoServiceBroker *MongoServiceBroker) Services(context context.Context) []brokerapi.Service {
	// TODO: read config

	planList := []brokerapi.ServicePlan{}
	for _, plan := range mongoServiceBroker.plans() {
		planList = append(planList, *plan)
	}

	services := []brokerapi.Service{
		brokerapi.Service{
			ID:            "SOME-UUID-98769-mongodb-service", // TODO: better uuid
			Name:          "footoo-mongodb",
			Description:   "A in development mongodb service",
			Bindable:      true,
			Tags:          []string{"FoOTOo", "mongodb"},
			PlanUpdatable: false,
			Plans:         planList,
			//Requires
			//Metadata
			//DashboardClient
		},
	}

	return services
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

	error := instanceCreator.Create(instanceID, details)
	if error != nil {
		return spec, error
	}

	return spec, nil
}

func (mongoServiceBroker *MongoServiceBroker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	for _, instanceCreator := range mongoServiceBroker.InstanceCreators {
		instanceExists, error := instanceCreator.InstanceExists(instanceID)

		if error != nil {
			return spec, error
		}

		if instanceExists {
			return spec, instanceCreator.Destroy(instanceID, details)
		}
	}

	return spec, nil
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

	error := instanceCreator.Update(instanceID, details)
	if error != nil {
		return spec, error
	}

	return spec, nil
}

func (mongoServiceBroker *MongoServiceBroker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	for key, instanceCreator := range mongoServiceBroker.InstanceCreators {
		instanceExists, error := instanceCreator.InstanceExists(instanceID)

		if error != nil {
			return binding, error
		}

		if instanceExists {
			instanceBinder, ok := mongoServiceBroker.InstanceBinders[key]
			if !ok {
				return binding, errors.New("instance binder not found for plan")
			}

			credentials, error := instanceBinder.Bind(instanceID, bindingID, details)
			binding.Credentials = credentials
			return binding, error
		}
	}

	return binding, brokerapi.ErrInstanceDoesNotExist
}

func (mongoServiceBroker *MongoServiceBroker) Unbind(context context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	for _, instanceBinder := range mongoServiceBroker.InstanceBinders {
		instanceExists, error := instanceBinder.InstanceBindingExists(instanceID, bindingID)

		if error != nil {
			return error
		}

		if instanceExists {
			error = instanceBinder.Unbind(instanceID, bindingID, details)
			return error
		}
	}

	return brokerapi.ErrInstanceDoesNotExist
}

func (mongoServiceBroker *MongoServiceBroker) plans() map[string]*brokerapi.ServicePlan {
	plans := map[string]*brokerapi.ServicePlan{}

	free := true

	plans["standard"] = &brokerapi.ServicePlan{
		ID:          "SOME-UUID-98769-standard", // TODO: better uuid
		Name:        "standard",
		Description: "Standard mongodb plan",
		Free:        &free,
		//Bindable:
		//Metadata:
	}

	return plans
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
