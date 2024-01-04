package managers

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	config "github.com/cloudfoundry-community/go-cfclient/v3/config"
)

type CloudFoundryProviderConfig struct {
	Endpoint                  string
	User                      string
	Password                  string
	SSOPasscode               string
	CFClientID                string
	CFClientSecret            string
	UaaClientID               string
	UaaClientSecret           string
	SkipSslValidation         bool
	AppLogsMax                int64
	DefaultQuotaName          string
	PurgeWhenDelete           bool
	StoreTokensPath           string
	ForceNotFailBrokerCatalog bool
}

type Session struct {
	cfclient *client.Client
}

func (c *CloudFoundryProviderConfig) NewSession() (*Session, error) {
	var cfg *config.Config
	var err error
	switch {
	case c.User != "" && c.Password != "":
		cfg, err = config.NewUserPassword(c.Endpoint, c.User, c.Password)
		if err != nil {
			return nil, err
		}
	}
	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	s := Session{
		cfclient: cf,
	}
	return &s, nil
}
