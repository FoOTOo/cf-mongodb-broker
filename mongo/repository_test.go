package mongo

import (
	"testing"

	"github.com/pivotal-cf/brokerapi"
)

const (
	InstanceID1       = "InstanceID1"
	BindingID1        = "BindingID1"
	ServiceID1        = "ServiceID1"
	ServiceID2        = "ServiceID2"
	PlanID1           = "PlanID1"
	PlanID2           = "PlanID2"
	OrganizationGUID1 = "OrganizationGUID1"
	SpaceGUID1        = "SpaceGUID1"
	AppGUID1          = "AppGUID1"
)

func TestRepository(t *testing.T) {
	//-------------------------
	t.Log("Create repo")
	adminService, error := NewAdminService("172.16.0.156", "rootusername", "rootpassword", "admin") // TODO: change

	if error != nil {
		t.Fatal("Error: ", error)
	}

	repository := NewRepository(adminService)

	//-------------------------
	t.Log("Create instance")
	instanceExists, error := repository.InstanceExists(InstanceID1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if instanceExists {
		t.Fatal("Instance should NOT exist")
	}

	provisionDetails := brokerapi.ProvisionDetails{
		ServiceID:        ServiceID1,
		PlanID:           PlanID1,
		OrganizationGUID: OrganizationGUID1,
		SpaceGUID:        SpaceGUID1,
	}

	error = repository.SaveInstance(InstanceID1, provisionDetails)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	instanceExists, error = repository.InstanceExists(InstanceID1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if !instanceExists {
		t.Fatal("Instance should exist")
	}

	//-------------------------
	t.Log("Update instance binding")
	updateDetails := brokerapi.UpdateDetails{
		ServiceID: ServiceID2,
		PlanID:    PlanID2,
	}

	error = repository.UpdateInstance(InstanceID1, updateDetails)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	//-------------------------
	t.Log("Save instance binding")

	instanceBindingExists, error := repository.InstanceBindingExists(InstanceID1, BindingID1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if instanceBindingExists {
		t.Fatal("Instance binding should NOT exist")
	}

	bindDetails := brokerapi.BindDetails{
		AppGUID: AppGUID1,
	}
	error = repository.SaveInstanceBinding(InstanceID1, BindingID1, bindDetails)
	if error != nil {
		t.Fatal("Error: ", error)
	}

	instanceBindingExists, error = repository.InstanceBindingExists(InstanceID1, BindingID1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if !instanceBindingExists {
		t.Fatal("Instance binding should exist")
	}

	//-------------------------
	t.Log("Delete instance binding")

	unbindDetails := brokerapi.UnbindDetails{}
	error = repository.DeleteInstanceBinding(InstanceID1, BindingID1, unbindDetails)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	instanceBindingExists, error = repository.InstanceBindingExists(InstanceID1, BindingID1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if instanceBindingExists {
		t.Fatal("Instance binding should NOT exist")
	}

	//-------------------------
	t.Log("Delete instance")

	deprovisionDetails := brokerapi.DeprovisionDetails{}
	error = repository.DeleteInstance(InstanceID1, deprovisionDetails)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	instanceExists, error = repository.InstanceExists(InstanceID1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if instanceExists {
		t.Fatal("Instance should NOT exist")
	}
}
