package managers

import (
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	config "github.com/cloudfoundry-community/go-cfclient/v3/config"
)

type CloudFoundryProviderConfig struct {
	Endpoint          string
	User              string
	Password          string
	CFClientID        string
	CFClientSecret    string
	SkipSslValidation bool
}

type Session struct {
	CFClient *client.Client
}

func (c *CloudFoundryProviderConfig) NewSession(httpClient *http.Client) (*Session, error) {
	var cfg *config.Config
	var err error
	switch {
	case c.User != "" && c.Password != "":
		cfg, err = config.NewUserPassword(c.Endpoint, c.User, c.Password)
	case c.CFClientID != "" && c.CFClientSecret != "":
		cfg, err = config.NewClientSecret(c.Endpoint, c.CFClientID, c.CFClientSecret)
	default:
		cfg, err = config.NewFromCFHome()
	}
	if err != nil {
		return nil, err
	}
	if httpClient != nil {
		cfg.WithHTTPClient(httpClient)
	}
	cfg.WithSkipTLSValidation(c.SkipSslValidation)
	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	s := Session{
		CFClient: cf,
	}
	return &s, nil
}
