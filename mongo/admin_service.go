package mongo

import (
	"errors"
	"log"
	"strings"

	"time"

	"encoding/json"

	"os"

	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	logEnabled = true
)

type AdminService struct {
	hosts      string
	username   string
	password   string
	authSource string

	addresses []string

	session *mgo.Session
}

func NewAdminService(hosts, username, password, authSource string) (*AdminService, error) {
	if logEnabled {
		logger := log.New(os.Stdout, "mongo-broker-mongo:", 0)
		mgo.SetDebug(true)
		mgo.SetLogger(logger)
	}

	addresses, error := splitHosts(hosts, "27017")

	if error != nil {
		return nil, error
	}

	dialInfo := &mgo.DialInfo{
		Addrs:   addresses,
		Direct:  false,
		Timeout: 30 * time.Second,
		//Database: authSource,
		Source:    authSource,
		Mechanism: "SCRAM-SHA-1",
		Username:  username,
		Password:  password,
	}

	session, error := mgo.DialWithInfo(dialInfo)

	if error != nil {
		return nil, error
	}

	adminService := &AdminService{
		hosts:      hosts,
		username:   username,
		password:   password,
		authSource: authSource,
		addresses:  addresses,
		session:    session,
	}

	return adminService, nil
}

func (adminService *AdminService) DatabaseExists(databaseName string) (bool, error) {
	//session := adminService.session.Copy()
	//defer session.Close()
	session := adminService.session

	databaseNames, error := session.DatabaseNames()

	if error != nil {
		return false, error
	}

	for _, name := range databaseNames {
		if databaseName == name {
			return true, nil
		}
	}

	return false, nil
}

func (adminService *AdminService) DeleteDatabase(databaseName string) error {
	session := adminService.session
	database := session.DB(databaseName)

	if database.Name != databaseName {
		return errors.New("Database not exist: " + databaseName)
	}

	// TODO: Remove db owner role

	error := database.DropDatabase()

	if error != nil {
		return error
	}

	return nil
}

func (adminService *AdminService) CreateDatabase(databaseName string) (*mgo.Database, error) {
	error := adminService.addDBOwnerRole(databaseName)

	if error != nil {
		return nil, error
	}

	session := adminService.session
	database := session.DB(databaseName)
	collection := database.C("foo")
	error = collection.Insert(&bson.DocElem{"foo", "bar"})

	if error != nil {
		return nil, error
	}

	error = collection.DropCollection()

	if error != nil {
		return nil, error
	}

	error = session.Fsync(false)
	if error != nil {
		return nil, error
	}

	return database, nil
}

func (adminService *AdminService) addDBOwnerRole(databaseName string) error {
	session := adminService.session
	database := session.DB(adminService.authSource)

	// TODO: ??? Not sure if it's correct
	//roles := []bson.DocElem{{"role", "dbOwner"}, {"db", databaseName}}
	roles := []interface{}{map[string]string{"role": "dbOwner", "db": databaseName}}
	cmd := &bson.D{{"grantRolesToUser", adminService.username}, {"roles", roles}}

	//cmd := map[string]interface{}{
	//	"grantRolesToUser": adminService.username,
	//	"roles": []interface{}{
	//		map[string]string{
	//			"role": "dbOwner",
	//			"db":   databaseName,
	//		},
	//	},
	//}

	result := &bson.D{}
	fmt.Println("+++++++++++++")

	error := database.Run(cmd, result)

	fmt.Print("===========================")
	fmt.Println(result.Map())

	if error != nil {
		return error
	}

	ok := result.Map()["ok"]

	if ok != 1.0 {
		jsonStr, error := json.MarshalIndent(result.Map(), "", "  ")
		if error != nil {
			return error
		}

		return errors.New(string(jsonStr))
	}

	//user := &mgo.User{
	//	Username: adminService.username,
	//	//Roles: []mgo.Role{
	//	//	"dbOwner",
	//	//},
	//	OtherDBRoles: map[string][]mgo.Role{
	//		databaseName:                           {"dbOwner"},
	//		"a0c320a8-4108-4f6d-9a59-b91193a6073c": {"dbOwner"},
	//		"admin": {"root"},
	//	},
	//}
	//
	//error := database.UpsertUser(user)
	//
	//if error != nil {
	//	return error
	//}

	return nil
}

func (adminService *AdminService) CreateUser(databaseName, username, password string) error {
	session := adminService.session
	database := session.DB(databaseName)

	// TODO: ??? Not sure if it's correct
	roles := []bson.DocElem{{"role", "readWrite"}, {"db", databaseName}}
	cmd := &bson.D{{"createUser", username}, {"pwd", password}, {"roles", roles}}

	result := &bson.D{}
	error := database.Run(cmd, result)

	if error != nil {
		return error
	}

	fmt.Print("====== ")
	fmt.Println(result)

	ok := result.Map()["ok"]

	if ok != 1.0 {
		jsonStr, error := json.MarshalIndent(result.Map(), "", "  ")
		if error != nil {
			return error
		}

		return errors.New(string(jsonStr))
	}

	return nil
}

func (adminService *AdminService) DeleteUser(databaseName, username string) error {
	session := adminService.session
	database := session.DB(databaseName)

	cmd := &bson.D{{"dropUser", username}}

	result := &bson.D{}
	error := database.Run(cmd, result)

	if error != nil {
		return error
	}

	ok := result.Map()["ok"]

	if ok != 1.0 {
		jsonStr, error := json.MarshalIndent(result.Map(), "", "  ")
		if error != nil {
			return error
		}

		return errors.New(string(jsonStr))
	}

	return nil
}

func (adminService *AdminService) GetConnectionString(databaseName, username, password string) string {
	parts := []string{"mongodb://", username, ":", password, "@", adminService.GetServerAddresses(), "/", databaseName}
	return strings.Join(parts, "")
}

func (adminService *AdminService) GetServerAddresses() string {
	return strings.Join(adminService.addresses, ",")
}

func (adminService *AdminService) SaveDoc(doc interface{}, databaseName string, collectionName string) error {
	//databaseExists, error := adminService.DatabaseExists(databaseName)
	//
	//var database *mgo.Database
	//
	//if error != nil {
	//	return error
	//}
	//
	//if !databaseExists {
	//	database, error = adminService.CreateDatabase(databaseName)
	//
	//	if error != nil {
	//		return error
	//	}
	//}

	session := adminService.session

	database := session.DB(databaseName)
	collection := database.C(collectionName)
	error := collection.Insert(doc)

	if error != nil {
		return error
	}

	error = session.Fsync(false)
	if error != nil {
		return nil
	}

	return nil
}

func (adminService *AdminService) RemoveDoc(selector interface{}, databaseName string, collectionName string) error {
	//databaseExists, error := adminService.DatabaseExists(databaseName)
	//
	//var database *mgo.Database
	//
	//if error != nil {
	//	return error
	//}
	//
	//if !databaseExists {
	//	return errors.New("Database not exists")
	//}

	session := adminService.session

	database := session.DB(databaseName)
	collection := database.C(collectionName)

	error := collection.Remove(selector)

	if error != nil {
		return error
	}

	return nil
}

func (adminService *AdminService) DocExists(query *bson.M, databaseName string, collectionName string) (bool, error) {
	session := adminService.session
	database := session.DB(databaseName)
	collection := database.C(collectionName)

	result := &bson.D{}
	error := collection.Find(query).One(result)

	if error != nil {
		return false, error
	}

	return result != nil, nil
}

func splitHosts(hosts string, defaultPort string) ([]string, error) {
	var addresses []string
	for _, hostWithPort := range strings.Split(hosts, ",") {
		arr := strings.Split(hostWithPort, ":")
		if len(arr) == 0 || len(arr) > 2 {
			return addresses, errors.New("Bad hosts string: " + hosts)
		} else if len(arr) == 1 {
			addresses = append(addresses, strings.TrimSpace(arr[0])+":"+defaultPort)
		} else {
			addresses = append(addresses, strings.TrimSpace(arr[0])+":"+strings.TrimSpace(arr[1]))
		}
	}

	return addresses, nil
}
