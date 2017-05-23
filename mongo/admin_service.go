package mongo

import (
	"errors"
	"log"
	"strings"

	"time"

	"encoding/json"

	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	logEnabled = false
)

type AdminService struct {
	hosts       string
	username    string
	password    string
	replSetName string
	authSource  string

	addresses []string

	dialInfo *mgo.DialInfo
}

func NewAdminService(hosts, username, password, replSetName, authSource string) (*AdminService, error) {
	if logEnabled {
		logger := log.New(os.Stdout, "mongo-broker-mongo:", 0)
		mgo.SetDebug(true)
		mgo.SetLogger(logger)
	}

	addresses, err := splitHosts(hosts, "27017")

	if err != nil {
		return nil, err
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

	adminService := &AdminService{
		hosts:       hosts,
		username:    username,
		password:    password,
		replSetName: replSetName,
		authSource:  authSource,
		addresses:   addresses,
		dialInfo:    dialInfo,
	}

	return adminService, nil
}

func (adminService *AdminService) newSession() (*mgo.Session, error) {
	session, err := mgo.DialWithInfo(adminService.dialInfo)

	return session, err
}

func (adminService *AdminService) DatabaseExists(databaseName string) (bool, error) {
	session, err := adminService.newSession()

	if err != nil {
		return false, err
	}

	defer session.Close()

	databaseNames, err := session.DatabaseNames()

	if err != nil {
		return false, err
	}

	for _, name := range databaseNames {
		if databaseName == name {
			return true, nil
		}
	}

	return false, nil
}

func (adminService *AdminService) DeleteDatabase(databaseName string) error {
	session, err := adminService.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	database := session.DB(databaseName)

	if database.Name != databaseName {
		return errors.New("Database not exist: " + databaseName)
	}

	// TODO: Remove db owner role

	err = database.DropDatabase()

	if err != nil {
		return err
	}

	return nil
}

func (adminService *AdminService) CreateDatabase(databaseName string) (*mgo.Database, error) {
	session, err := adminService.newSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	err = adminService.addDBOwnerRole(session, databaseName)

	if err != nil {
		return nil, err
	}

	database := session.DB(databaseName)
	collection := database.C("foo")
	err = collection.Insert(&bson.DocElem{"foo", "bar"})

	if err != nil {
		return nil, err
	}

	err = collection.DropCollection()

	if err != nil {
		return nil, err
	}

	err = session.Fsync(false)
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (adminService *AdminService) addDBOwnerRole(session *mgo.Session, databaseName string) error {
	database := session.DB(adminService.authSource)

	roles := []interface{}{map[string]string{"role": "dbOwner", "db": databaseName}}
	cmd := &bson.D{{"grantRolesToUser", adminService.username}, {"roles", roles}}

	result := &bson.D{}

	err := database.Run(cmd, result)

	//fmt.Print("===========================")
	//fmt.Println(result.Map())

	if err != nil {
		return err
	}

	ok := result.Map()["ok"]

	if ok != 1.0 {
		jsonStr, err := json.MarshalIndent(result.Map(), "", "  ")
		if err != nil {
			return err
		}

		return errors.New(string(jsonStr))
	}

	return nil
}

func (adminService *AdminService) CreateUser(databaseName, username, password string) error {
	session, err := adminService.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	database := session.DB(databaseName)

	roles := []interface{}{map[string]string{"role": "readWrite", "db": databaseName}}
	cmd := &bson.D{{"createUser", username}, {"pwd", password}, {"roles", roles}}

	result := &bson.D{}
	err = database.Run(cmd, result)

	if err != nil {
		return err
	}

	//fmt.Print("====== ")
	//fmt.Println(result)

	ok := result.Map()["ok"]

	if ok != 1.0 {
		jsonStr, err := json.MarshalIndent(result.Map(), "", "  ")
		if err != nil {
			return err
		}

		return errors.New(string(jsonStr))
	}

	return nil
}

func (adminService *AdminService) DeleteUser(databaseName, username string) error {
	session, err := adminService.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	database := session.DB(databaseName)

	cmd := &bson.D{{"dropUser", username}}

	result := &bson.D{}
	err = database.Run(cmd, result)

	if err != nil {
		return err
	}

	ok := result.Map()["ok"]

	if ok != 1.0 {
		jsonStr, err := json.MarshalIndent(result.Map(), "", "  ")
		if err != nil {
			return err
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
	session, err := adminService.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	database := session.DB(databaseName)
	collection := database.C(collectionName)
	err = collection.Insert(doc)

	if err != nil {
		return err
	}

	err = session.Fsync(false)
	if err != nil {
		return nil
	}

	return nil
}

func (adminService *AdminService) RemoveDoc(selector interface{}, databaseName string, collectionName string) error {
	session, err := adminService.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	database := session.DB(databaseName)
	collection := database.C(collectionName)

	err = collection.Remove(selector)

	if err != nil {
		return err
	}

	return nil
}

func (adminService *AdminService) DocExists(query *bson.M, databaseName string, collectionName string) (bool, error) {
	session, err := adminService.newSession()

	if err != nil {
		return false, err
	}

	defer session.Close()

	database := session.DB(databaseName)
	collection := database.C(collectionName)

	result := &bson.D{}
	err = collection.Find(query).One(result)

	//fmt.Println("================")
	//fmt.Println(result.Map())

	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}

	return len(result.Map()) > 0, nil
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

func (adminService *AdminService) UpdateDoc(selector interface{}, update interface{}, databaseName string, collectionName string) error {
	session, err := adminService.newSession()

	if err != nil {
		return err
	}

	defer session.Close()

	database := session.DB(databaseName)
	collection := database.C(collectionName)
	_, err = collection.Upsert(selector, update)

	if err != nil {
		return err
	}

	return nil
}

func (adminService *AdminService) GetOneDoc(query *bson.M, databaseName string, collectionName string) (bson.M, error) {
	session, err := adminService.newSession()

	if err != nil {
		return nil, err
	}

	defer session.Close()

	database := session.DB(databaseName)
	collection := database.C(collectionName)

	result := &bson.D{}
	err = collection.Find(query).One(result)

	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result.Map(), nil
}
