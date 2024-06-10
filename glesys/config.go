package glesys

import (
	"github.com/glesys/glesys-go/v8"
)

// Config - Provider configuration
type Config struct {
	UserID      string
	Token       string
	UserAgent   string
	APIEndpoint string
}

// Client - Setup new glesys client
func (c *Config) Client() (*glesys.Client, error) {
	client := glesys.NewClient(c.UserID, c.Token, "tf-glesys/0.11.3")

	err := client.SetBaseURL(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}
