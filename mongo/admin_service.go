package mongo

import (
	"errors"
	"log"
	"os"
	"strings"

	"time"

	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AdminService struct {
	hosts      string
	username   string
	password   string
	authSource string

	session *mgo.Session
}

func NewAdminService(hosts, username, password, authSource string) (*AdminService, error) {
	logger := log.New(os.Stdout, "mongo-broker-mongo:", 0)
	mgo.SetDebug(true)
	mgo.SetLogger(logger)

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
		Mechanism: "SCRAM_SHA_1",
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

	return database, nil
}

func (adminService *AdminService) addDBOwnerRole(databaseName string) error {
	session := adminService.session
	database := session.DB(databaseName)

	// TODO: ??? Not sure if it's correct
	roles := bson.D{{"role", "dbOwner"}, {"db", databaseName}}
	cmd := &bson.D{{"grantRolesToUser", adminService.username}, {"roles", roles}}

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
