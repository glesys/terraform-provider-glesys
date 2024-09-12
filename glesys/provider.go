package glesys

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider - Setup new Terraform Provider resource
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// specify what is needed to configure the provider.
			"userid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GLESYS_USERID", nil),
				Description: "UserId for the Glesys API. Alternatively, this can be set using the `GLESYS_USERID` environment variable",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GLESYS_TOKEN", nil),
				Description: "User token for the Glesys API. Alternatively, this can be set using the `GLESYS_TOKEN` environment variable",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GLESYS_API_URL", "https://api.glesys.com"),
				Description: "The base URL to use for the GleSYS API requests. (Defaults to the value of the `GLESYS_API_URL` environment variable or `https://api.glesys.com` if unset.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"glesys_dnsdomain":      dataSourceGlesysDNSDomain(),
			"glesys_network":        dataSourceGlesysNetwork(),
			"glesys_networkadapter": dataSourceGlesysNetworkAdapter(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"glesys_dnsdomain":                resourceGlesysDNSDomain(),
			"glesys_dnsdomain_record":         resourceGlesysDNSDomainRecord(),
			"glesys_emailaccount":             resourceGlesysEmailAccount(),
			"glesys_emailalias":               resourceGlesysEmailAlias(),
			"glesys_loadbalancer":             resourceGlesysLoadBalancer(),
			"glesys_loadbalancer_backend":     resourceGlesysLoadBalancerBackend(),
			"glesys_loadbalancer_frontend":    resourceGlesysLoadBalancerFrontend(),
			"glesys_loadbalancer_target":      resourceGlesysLoadBalancerTarget(),
			"glesys_network":                  resourceGlesysNetwork(),
			"glesys_networkadapter":           resourceGlesysNetworkAdapter(),
			"glesys_server":                   resourceGlesysServer(),
			"glesys_server_disk":              resourceGlesysServerDisk(),
			"glesys_objectstorage_instance":   resourceGlesysObjectStorageInstance(),
			"glesys_objectstorage_credential": resourceGlesysObjectStorageCredential(),
			"glesys_ip":                       resourceGlesysIP(),
		},
		// this will be used to configure the client to communicate with the API
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		UserID:      d.Get("userid").(string),
		Token:       d.Get("token").(string),
		APIEndpoint: d.Get("api_endpoint").(string),
	}
	return config.Client()
}
