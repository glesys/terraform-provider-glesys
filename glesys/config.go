package glesys

import (
	"github.com/glesys/glesys-go/v4"
)

// Config - Provider configuration
type Config struct {
	UserID    string
	Token     string
	UserAgent string
}

// Client - Setup new glesys client
func (c *Config) Client() (*glesys.Client, error) {
	client := glesys.NewClient(c.UserID, c.Token, "tf-glesys/0.5.0")

	return client, nil
}
