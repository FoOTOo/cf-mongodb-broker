package mongo

import (
	"testing"
)

//func TestSplitHosts(t *testing.T) {
//	addresses, err := splitHosts("127.0.0.1:1234, 172.16.0.1,192.168.1.1", "9999")
//
//	fmt.Println("==============")
//	fmt.Println(err)
//	fmt.Println(addresses)
//	fmt.Println("--------------")
//}

const (
	DB1 = "TestDatabase1"
)

func TestAdminService(t *testing.T) {
	t.Log("Creating admin service")
	adminService, error := NewAdminService("172.16.0.156", "rootusername", "rootpassword", "admin") // TODO: change

	if error != nil {
		t.Fatal("Error: ", error)
	}

	databaseExists, error := adminService.DatabaseExists(DB1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	if databaseExists {
		t.Fatal("Database %s should not exist", DB1)
	}

	_, error = adminService.CreateDatabase(DB1)

	if error != nil {
		t.Fatal("Error: ", error)
	}

	//error = adminService.DeleteDatabase(DB1)
	//
	//if error != nil {
	//	t.Fatal("Error: ", error)
	//}
}
