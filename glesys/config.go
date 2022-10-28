package glesys

import (
	"github.com/glesys/glesys-go/v6"
)

// Config - Provider configuration
type Config struct {
	UserID    string
	Token     string
	UserAgent string
}

// Client - Setup new glesys client
func (c *Config) Client() (*glesys.Client, error) {
	client := glesys.NewClient(c.UserID, c.Token, "tf-glesys/0.7.1")

	return client, nil
}
