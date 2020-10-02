package glesys

import (
	"github.com/glesys/glesys-go/v2"
)

type Config struct {
	UserId    string
	Token     string
	UserAgent string
}

func (c *Config) Client() (*glesys.Client, error) {
	client := glesys.NewClient(c.UserId, c.Token, "tf-glesys/0.0.1")

	return client, nil
}
