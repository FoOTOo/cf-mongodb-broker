package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

type MongoServiceBroker struct{}

func (mongoServiceBroker *MongoServiceBroker) Services(context context.Context) []brokerapi.Service {
	// TODO: read config

	free := true

	plans := []brokerapi.ServicePlan{
		brokerapi.ServicePlan{
			ID:          "SOME-UUID-98769-standard", // TODO: better uuid
			Name:        "standard",
			Description: "Standard mongodb plan",
			Free:        &free,
			//Bindable:
			//Metadata:
		},
	}

	services := []brokerapi.Service{
		brokerapi.Service{
			ID:            "SOME-UUID-98769-mongodb-service", // TODO: better uuid
			Name:          "footoo-mongodb",
			Description:   "A in development mongodb service",
			Bindable:      true,
			Tags:          []string{"FoOTOo", "mongodb"},
			PlanUpdatable: false,
			Plans:         plans,
			//Requires
			//Metadata
			//DashboardClient
		},
	}

	return services
}

func (mongoServiceBroker *MongoServiceBroker) Provision(context context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	spec := brokerapi.ProvisionedServiceSpec{}

	// TODO:
	// 1. exist ?
	if serviceDetails.PlanID == "" {
		return spec, errors.New("plan_id required")
	}
	// 2. select plan based on planID
	// 3. create instance

	return spec, nil
}

func (mongoServiceBroker *MongoServiceBroker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	// TODO:

	return spec, nil
}

func (mongoServiceBroker *MongoServiceBroker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	// TODO:

	return binding, nil
}

func (mongoServiceBroker *MongoServiceBroker) Unbind(context context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	// TODO: nil

	return nil
}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (mongoServiceBroker *MongoServiceBroker) LastOperation(context context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

func (mongoServiceBroker *MongoServiceBroker) Update(context context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}
