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
	adminService, err := NewAdminService("172.16.0.156", "rootusername", "rootpassword", "admin")

	if err != nil {
		t.Fatal("Error: ", err)
	}

	//-------------------------
	t.Log("Creating database")
	databaseExists, err := adminService.DatabaseExists(DB1)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	if databaseExists {
		t.Fatal("Database %s should not exist", DB1)
	}

	_, err = adminService.CreateDatabase(DB1)

	if err != nil {
		t.Fatal("Error: ", err)
	}

	//err = adminService.CreateUser(DB1, User1, Pwd1)
	//
	//if err != nil {
	//	t.Fatal("Error: ", err)
	//}
	//
	//err = adminService.DeleteUser(DB1, User1)
	//
	//if err != nil {
	//	t.Fatal("Error: ", err)
	//}

	//-------------------------
	t.Log("Save Doc")
	doc := map[string]string{"_id": "123", "v": "456"}
	err = adminService.SaveDoc(doc, DB1, Collection1)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	query := &bson.M{"_id": "123"}

	docExists, err := adminService.DocExists(query, DB1, Collection1)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if !docExists {
		t.Fatal("Doc should exists")
	}

	//-------------------------
	t.Log("Update Doc")
	update := map[string]string{"v": "789"}
	err = adminService.UpdateDoc(query, update, DB1, Collection1)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	result, err := adminService.GetOneDoc(query, DB1, Collection1)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if result == nil {
		t.Fatal("Updated doc should exist")
	}

	if result["v"] != "789" {
		t.Fatal("Doc no updated correctly", result)
	}

	//-------------------------
	t.Log("Remove Doc")
	err = adminService.RemoveDoc(query, DB1, Collection1)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	docExists, err = adminService.DocExists(query, DB1, Collection1)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if docExists {
		t.Fatal("Doc should NOT exists")
	}

	//-------------------------
	t.Log("Delete Database")
	err = adminService.DeleteDatabase(DB1)

	if err != nil {
		t.Fatal("Error: ", err)
	}
}
