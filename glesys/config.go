package glesys

import (
	"github.com/glesys/glesys-go/v3"
)

// Config - Provider configuration
type Config struct {
	UserID    string
	Token     string
	UserAgent string
}

// Client - Setup new glesys client
func (c *Config) Client() (*glesys.Client, error) {
	client := glesys.NewClient(c.UserID, c.Token, "tf-glesys/0.4.4")

	return client, nil
}
