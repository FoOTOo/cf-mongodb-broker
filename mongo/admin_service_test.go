package mongo

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
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
	DB1         = "TestDatabase1"
	User1       = "user1"
	Pwd1        = "pwd1"
	Collection1 = "collection1"
)

func TestAdminService(t *testing.T) {
	//-------------------------
	t.Log("Creating admin service")
	adminService, error := NewAdminService("172.16.0.156", "rootusername", "rootpassword", "admin") // TODO: change

	if error != nil {
		t.Fatal("Error: ", error)
	}

	//-------------------------
	t.Log("Creating database")
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

	//error = adminService.CreateUser(DB1, User1, Pwd1)
	//
	//if error != nil {
	//	t.Fatal("Error: ", error)
	//}

	//error = adminService.DeleteUser(DB1, User1)
	//
	//if error != nil {
	//	t.Fatal("Error: ", error)
	//}

	//-------------------------
	t.Log("Save Doc")
	doc := map[string]string{"_id": "123", "v": "456"}
	error = adminService.SaveDoc(doc, DB1, Collection1)
	if error != nil {
		t.Fatal("Error: ", error)
	}

	query := &bson.M{"_id": "123"}

	docExists, error := adminService.DocExists(query, DB1, Collection1)
	if error != nil {
		t.Fatal("Error: ", error)
	}

	if !docExists {
		t.Fatal("Doc should exists")
	}

	//-------------------------
	t.Log("Remove Doc")
	error = adminService.RemoveDoc(query, DB1, Collection1)
	if error != nil {
		t.Fatal("Error: ", error)
	}

	docExists, error = adminService.DocExists(query, DB1, Collection1)
	if error != nil {
		t.Fatal("Error: ", error)
	}

	if docExists {
		t.Fatal("Doc should NOT exists")
	}

	//-------------------------
	t.Log("Delete Database")
	error = adminService.DeleteDatabase(DB1)

	if error != nil {
		t.Fatal("Error: ", error)
	}
}
