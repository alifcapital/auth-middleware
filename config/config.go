package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module ...
var Module = fx.Provide(NewConfig)

// Config ...
type Config interface {
	LoadConfig()
	Get() *ConfigJSON
}

type config struct {
	configJSON *ConfigJSON
}

// Params ...
type Params struct {
	fx.In
}

// NewConfig ...
func NewConfig(p Params) Config {

	configItem := &config{
		configJSON: &ConfigJSON{},
	}

	configItem.LoadConfig()

	return configItem
}

// Db ...
type Db struct {
	Driver   string `json:"driver"`
	Username string `json:"username"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Params   string `json:"params"`
	Migrate  string `json:"migrate"`
}

// Sentry ...
type Sentry struct {
	Dsn         string `json:"dsn"`
	Environment string `json:"environment"`
}

// Minio ...
type Minio struct {
	Endpoint string `json:"endpoint"`
	URL      string `json:"url"`
}

// Aws ...
type Aws struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

// Rabbit ...
type Rabbit struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
	Protocol string `json:"protocol"`
	Exchange string `json:"exchange"`
}

// PubKey ...
type PubKey struct {
	URL  string `json:"url"`
	Data string `json:"data"`
}

// API ...
type API struct {
	FrontAlif APIData   `json:"frontAlif"`
	CRM       APIData   `json:"crm"`
	CBS       APIDataV2 `json:"cbs"`
}

// APIData ...
type APIData struct {
	Token   string `json:"token"`
	Service string `json:"service"`
}

// APIDataV2 ...
type APIDataV2 struct {
	Service           string `json:"service"`
	TokenExpressLimit string `json:"tokenExpressLimit"`
	TokenCredits      string `json:"tokenCredits"`
}

// Opentracing ...
type Opentracing struct {
	LogsEnabled      bool   `json:"logsEnabled"`
	ServiceName      string `json:"serviceName"`
	SamplerType      string `json:"samplerType"`
	SamplerParam     int    `json:"samplerParam"`
	ReporterLogSpans bool   `json:"reporterLogSpans"`
}

// Keycloak ...
type Keycloak struct {
	Issuer             string `json:"issuer"`
	Cache              int    `json:"cache"`
	SignatureAlgorithm string `json:"signatureAlgorithm"`
	DefaultScope       string `json:"defaultScope"`
}

// Redis ...
type Redis struct {
	Db   int    `json:"db"`
	Port int    `json:"port"`
	Host string `json:"hostAndPort"`
	Pass string `json:"pass"`
}
type DBs struct {
	Committee Db `json:"committee"`
	Crm       Db `json:"crm"`
}

// ConfigJSON ...
type ConfigJSON struct {
	BaseAPI                   string      `json:"baseAPI"`
	LogPath                   string      `json:"logPath"`
	ServerPort                string      `json:"serverPort"`
	ServiceName               string      `json:"serviceName"`
	WorkerIntervalAutoUpdate  int         `json:"workerIntervalAutoUpdate"`
	AuthMiddlewareServiceName string      `json:"authMiddlewareServiceName"`
	CRMToken                  string      `json:"crmToken"`
	Db                        DBs         `json:"db"`
	PubKey                    PubKey      `json:"pubKey"`
	Sentry                    Sentry      `json:"sentry"`
	Minio                     Minio       `json:"minio"`
	Aws                       Aws         `json:"aws"`
	Rabbit                    Rabbit      `json:"rabbit"`
	API                       API         `json:"api"`
	Opentracing               Opentracing `json:"opentracing"`
	Keycloak                  Keycloak    `json:"keycloak"`
	Redis                     Redis       `json:"redis"`
}

func (c *config) LoadConfig() {
	cfg := viper.New()

	cfg.AddConfigPath("./")
	cfg.AddConfigPath("../")
	cfg.AddConfigPath("../../")
	cfg.AddConfigPath("../../../")
	cfg.SetConfigName("config-dev")
	cfg.SetConfigType("json")
	cfg.AutomaticEnv()

	err := cfg.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("get - ConfigJSON - ReadInConfig: %w", err))
	}

	err = cfg.Unmarshal(&c.configJSON)
	if err != nil {
		panic(fmt.Errorf("get - ConfigJSON - Unmarshal: %w", err))
	}
}

func (c *config) Get() *ConfigJSON {
	return c.configJSON
}
