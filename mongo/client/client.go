package client

import (
	"strings"

	"errors"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

type Client struct {
	hosts      string
	username   string
	password   string
	authSource string
}

func NewClient(hosts, username, password, authSource string) (Client, error) {
	//credential := mgo.Credential{
	//	Username:  username,
	//	Password:  password,
	//	Source:    authSource,
	//	Mechanism: "SCRAM_SHA_1",
	//}

	logger := log.New(os.Stdout, "mongo-broker-mongo:", 0)
	mgo.SetDebug(true)
	mgo.SetLogger(logger)

	addresses, error := splitHosts(hosts)

	if error != nil {
		return nil, error
	}

	dialInfo := mgo.DialInfo{
		Addrs: addresses,
		Direct: false,
		Timeout: 30,


	}

	return Client{}, nil
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
