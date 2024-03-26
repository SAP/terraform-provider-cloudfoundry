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
	Origin            string
}

type Session struct {
	CFClient *client.Client
}

func (c *CloudFoundryProviderConfig) NewSession(httpClient *http.Client) (*Session, error) {
	var cfg *config.Config
	var err error
	var opts []config.Option

	if httpClient != nil {
		opts = append(opts, config.HttpClient(httpClient))
	}
	if c.SkipSslValidation {
		opts = append(opts, config.SkipTLSValidation())
	}
	switch {
	case c.User != "" && c.Password != "":
		opts = append(opts, config.UserPassword(c.User, c.Password))
		if c.Origin != "" {
			opts = append(opts, config.Origin(c.Origin))
		}
		cfg, err = config.New(c.Endpoint, opts...)
	case c.CFClientID != "" && c.CFClientSecret != "":
		opts = append(opts, config.ClientCredentials(c.CFClientID, c.CFClientSecret))
		cfg, err = config.New(c.Endpoint, opts...)
	default:
		cfg, err = config.NewFromCFHome(opts...)
	}
	if err != nil {
		return nil, err
	}
	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	s := Session{
		CFClient: cf,
	}
	return &s, nil
}
