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
	adminService, err := NewAdminService("172.16.0.156", "rootusername", "rootpassword", "admin") // TODO: change

	if err != nil {
		t.Fatal("Error: ", err)
	}

	repository := NewRepository(adminService)

	//-------------------------
	t.Log("Create instance")
	instanceExists, err := repository.InstanceExists(InstanceID1)

	if err != nil {
		t.Fatal("Error: ", err)
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

	err = repository.SaveInstance(InstanceID1, provisionDetails)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	instanceExists, err = repository.InstanceExists(InstanceID1)

	if err != nil {
		t.Fatal("Error: ", err)
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

	err = repository.UpdateInstance(InstanceID1, updateDetails)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	//-------------------------
	t.Log("Save instance binding")

	instanceBindingExists, err := repository.InstanceBindingExists(InstanceID1, BindingID1)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	if instanceBindingExists {
		t.Fatal("Instance binding should NOT exist")
	}

	bindDetails := brokerapi.BindDetails{
		AppGUID: AppGUID1,
	}
	err = repository.SaveInstanceBinding(InstanceID1, BindingID1, bindDetails)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	instanceBindingExists, err = repository.InstanceBindingExists(InstanceID1, BindingID1)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	if !instanceBindingExists {
		t.Fatal("Instance binding should exist")
	}

	//-------------------------
	t.Log("Delete instance binding")

	unbindDetails := brokerapi.UnbindDetails{}
	err = repository.DeleteInstanceBinding(InstanceID1, BindingID1, unbindDetails)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	instanceBindingExists, err = repository.InstanceBindingExists(InstanceID1, BindingID1)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	if instanceBindingExists {
		t.Fatal("Instance binding should NOT exist")
	}

	//-------------------------
	t.Log("Delete instance")

	deprovisionDetails := brokerapi.DeprovisionDetails{}
	err = repository.DeleteInstance(InstanceID1, deprovisionDetails)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	instanceExists, err = repository.InstanceExists(InstanceID1)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	if instanceExists {
		t.Fatal("Instance should NOT exist")
	}
}
