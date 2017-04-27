package config

import (
	"os"

	"strings"

	"fmt"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/pivotal-cf/brokerapi"
)

type Config struct {
	MongoConfig  MongoConfig  `yaml:"mongod"`
	BrokerConfig BrokerConfig `yaml:"broker"`
}

type MongoConfig struct {
	Nodes    MongoNodes    `yaml:"nodes"`
	RootUser MongoRootUser `yaml:"root"`
}

type MongoRootUser struct {
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
}

type MongoNodes struct {
	IPs  []string `yaml:"ips"`
	Port string   `yaml:"port"`
}

type BrokerConfig struct {
	Host           string `yaml:"host"`
	BrokerUsername string `yaml:"security_user_name"`
	BrokerPassword string `yaml:"security_user_password"`

	ServiceID          string          `yaml:"id"`
	ServiceName        string          `yaml:"name"`
	ServiceDescription string          `yaml:"description"`
	ServiceBindable    bool            `yaml:"bindable"`
	PlanUpdateable     bool            `yaml:"plan_updateable"`
	Plans              []Plan          `yaml:"plans"`
	Tags               []string        `yaml:"tags"`
	ServiceMetadata    ServiceMetadata `yaml:"metadata"`
}

type ServiceMetadata struct {
	DisplayName         string `yaml:"displayName"`
	IconImage           string `yaml:"iconImage"`
	LongDescription     string `yaml:"longDescription"`
	ProviderDisplayName string `yaml:"providerDisplayName"`
	DocumentationUrl    string `yaml:"documentationUrl"`
	SupportUrl          string `yaml:"supportUrl"`
	//ImageUrl            string `yaml:"imageUrl"`
}

type Plan struct {
	ID          string       `yaml:"id"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Free        bool         `yaml:"free"`
	Metadata    PlanMetadata `yaml:"metadata"`
}

type PlanMetadata struct {
	DisplayName string             `yaml:"displayName,omitempty"`
	Bullets     []string           `yaml:"bullets,omitempty"`
	Costs       []PlanMetadataCost `yaml:"costs,omitempty"`
}

type PlanMetadataCost struct {
	Amount map[string]float64 `yaml:"amount"`
	Unit   string             `yaml:"unit"`
}

func ParseConfig(path string) (Config, error) {
	file, error := os.Open(path)
	if error != nil {
		return Config{}, error
	}

	var config Config
	if err := candiedyaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config *Config) MongoHosts() string {
	var hosts []string
	for _, host := range config.MongoConfig.Nodes.IPs {
		hostWithPort := host + ":" + config.MongoConfig.Nodes.Port
		hosts = append(hosts, hostWithPort)
	}

	return strings.Join(hosts, ",")
}

func (config *Config) MongoUsername() string {
	return config.MongoConfig.RootUser.Username
}

func (config *Config) MongoPassword() string {
	return config.MongoConfig.RootUser.Password
}

func (config *Config) Services() []brokerapi.Service {
	planList := []brokerapi.ServicePlan{}
	for _, plan := range config.Plans() {
		planList = append(planList, *plan)
	}

	brokerConfig := config.BrokerConfig
	serviceMetadata := brokerConfig.ServiceMetadata

	services := []brokerapi.Service{
		brokerapi.Service{
			ID:            brokerConfig.ServiceID,
			Name:          brokerConfig.ServiceName,
			Description:   brokerConfig.ServiceDescription,
			Bindable:      brokerConfig.ServiceBindable,
			Tags:          brokerConfig.Tags,
			Plans:         planList,
			PlanUpdatable: brokerConfig.PlanUpdateable,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName: serviceMetadata.DisplayName,
				//ImageUrl:            serviceMetadata.ImageUrl,
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", serviceMetadata.IconImage),
				LongDescription:     serviceMetadata.LongDescription,
				ProviderDisplayName: serviceMetadata.ProviderDisplayName,
				DocumentationUrl:    serviceMetadata.DocumentationUrl,
				SupportUrl:          serviceMetadata.SupportUrl,
			},
			//Requires
			//DashboardClient
		},
	}

	return services
}

func (config *Config) Plans() map[string]*brokerapi.ServicePlan {
	plans := map[string]*brokerapi.ServicePlan{}

	for _, plan := range config.BrokerConfig.Plans {
		plans[plan.Name] = &brokerapi.ServicePlan{
			ID:          plan.ID,
			Name:        plan.Name,
			Description: plan.Description,
			Free:        &plan.Free,
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: plan.Metadata.DisplayName,
				Bullets:     plan.Metadata.Bullets,
				Costs:       planCosts(plan),
			},
		}
	}

	return plans
}

func planCosts(plan Plan) []brokerapi.ServicePlanCost {
	costs := []brokerapi.ServicePlanCost{}

	for _, cost := range plan.Metadata.Costs {
		costs = append(costs, brokerapi.ServicePlanCost{
			Amount: cost.Amount,
			Unit:   cost.Unit,
		})
	}

	return costs
}
