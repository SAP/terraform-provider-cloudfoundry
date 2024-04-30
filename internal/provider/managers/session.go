package managers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SAP/terraform-provider-cloudfoundry/internal/version"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	config "github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/hashicorp/terraform-plugin-framework/provider"
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

func (c *CloudFoundryProviderConfig) NewSession(httpClient *http.Client, req provider.ConfigureRequest) (*Session, error) {
	var cfg *config.Config
	var err error
	var opts []config.Option
	var finalAgent string

	cfUserAgent := os.Getenv("CF_APPEND_USER_AGENT")

	if len(strings.TrimSpace(cfUserAgent)) == 0 {
		finalAgent = fmt.Sprintf("Terraform/%s terraform-provider-cloudfoundry/%s", req.TerraformVersion, version.ProviderVersion)
	} else {
		finalAgent = fmt.Sprintf("Terraform/%s terraform-provider-cloudfoundry/%s %s", req.TerraformVersion, version.ProviderVersion, cfUserAgent)
	}
	opts = append(opts, config.UserAgent(finalAgent))

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
