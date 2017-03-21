package main

import (
	"net/http"
	"os"

	"github.com/pivotal-cf/brokerapi"

	"code.cloudfoundry.org/lager"
	"github.com/ultragtx/mongo-broker/broker"
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

	serviceBroker := &broker.MongoServiceBroker{}

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
