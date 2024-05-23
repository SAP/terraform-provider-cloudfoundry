package mta

import (
	"net/http"
)

type Configuration struct {
	BasePath      string            `json:"basePath,omitempty"`
	Host          string            `json:"host,omitempty"`
	Scheme        string            `json:"scheme,omitempty"`
	DefaultHeader map[string]string `json:"defaultHeader,omitempty"`
	UserAgent     string            `json:"userAgent,omitempty"`
	HTTPClient    *http.Client
}

func NewConfiguration(basePath string, userAgent string, client *http.Client) *Configuration {
	cfg := &Configuration{
		BasePath:      basePath,
		DefaultHeader: make(map[string]string),
		UserAgent:     userAgent,
		HTTPClient:    client,
	}
	return cfg
}

func (c *Configuration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
