package yi

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/soulteary/yi-openai-proxy/define"
	"github.com/spf13/viper"
)

func (c *StripPrefixConverter) Name() string {
	return "StripPrefix"
}

func (c *StripPrefixConverter) Convert(req *http.Request, config *DeploymentConfig) (*http.Request, error) {
	req.Host = config.EndpointUrl.Host
	req.URL.Scheme = config.EndpointUrl.Scheme
	req.URL.Host = config.EndpointUrl.Host
	req.URL.RawPath = req.URL.EscapedPath()

	query := req.URL.Query()
	query.Add("api-version", config.ApiVersion)
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func NewStripPrefixConverter(prefix string) *StripPrefixConverter {
	return &StripPrefixConverter{
		Prefix: prefix,
	}
}

func GetOptionFromEnv(key string) string {
	return strings.TrimSpace(viper.GetString(key))
}

func GetInstance() (err error) {
	var config DeploymentConfig
	apiVersion := GetOptionFromEnv(define.ENV_YI_API_VER)
	endpoint := GetOptionFromEnv(define.ENV_YI_ENDPOINT)
	apikey := GetOptionFromEnv(define.ENV_YI_API_KEY)
	modelName := GetOptionFromEnv(define.ENV_YI_MODEL_NAME)

	if endpoint == "" {
		return errors.New("endpoint should not be empty")
	}

	if modelName == "" {
		modelName = "yi-34b-chat"
	}

	if apiVersion == "" {
		apiVersion = "2023-07-01-preview"
	}

	config.ApiKey = apikey
	config.Endpoint = endpoint
	config.ApiVersion = apiVersion
	config.ModelName = modelName

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("parse endpoint error: %w", err)
	}
	config.EndpointUrl = u

	ModelDeploymentConfig[modelName] = config

	return nil
}
