package yi

import (
	"net/http"
	"net/url"
)

type DeploymentConfig struct {
	ModelName   string   `yaml:"model_name" json:"model_name" mapstructure:"model_name"`    // corresponding model name in openai
	Endpoint    string   `yaml:"endpoint" json:"endpoint" mapstructure:"endpoint"`          // deployment endpoint
	ApiKey      string   `yaml:"api_key" json:"api_key" mapstructure:"api_key"`             // secrect key1 or 2
	ApiVersion  string   `yaml:"api_version" json:"api_version" mapstructure:"api_version"` // deployment version, not required
	EndpointUrl *url.URL // url.URL form deployment endpoint
}

type RequestConverter interface {
	Name() string
	Convert(req *http.Request, config *DeploymentConfig) (*http.Request, error)
}

type StripPrefixConverter struct {
	Prefix string
}
