package main

import (
	"net/http"
	"os"

	"github.com/FoOTOo/mongo-service-broker-golang-ultragtx/broker"
	"github.com/pivotal-cf/brokerapi"

	"code.cloudfoundry.org/lager"
	"github.com/FoOTOo/mongo-service-broker-golang-ultragtx/mongo"
)

const (
	BrokerName     = "mongodb-broker"
	BrokerUsername = "mongodb-broker-user"
	BrokerPassword = "mongodb-broker-password"
)

func main() {
	brokerLogger := lager.NewLogger(BrokerName)
	//brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Starting Mongodb broker")

	// TODO: config file

	adminService, error := mongo.NewAdminService("172.16.0.156", "rootusername", "rootpassword", "admin") // TODO: change

	if error != nil {
		brokerLogger.Fatal("mongo-admin-service", error)
		return
	}

	repository := mongo.NewRepository(adminService)

	instanceCreator := mongo.NewInstanceCreator(adminService, repository)
	instanceBinder := mongo.NewInstanceBinder(adminService, repository)

	serviceBroker := &broker.MongoServiceBroker{
		InstanceCreators: map[string]broker.InstanceCreator{
			"standard": instanceCreator,
		},
		InstanceBinders: map[string]broker.InstanceBinder{
			"standard": instanceBinder,
		},
	}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: BrokerUsername,
		Password: BrokerPassword,
	}

	// broker
	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	// authWrapper := auth.NewWrapper(brokerCredentials.Username, brokerCredentials.Password)
	// TODO: ??? /instance /debug

	http.Handle("/", brokerAPI)

	brokerLogger.Fatal("http-listen", http.ListenAndServe(":"+"9876", nil)) // TODO: config
}
