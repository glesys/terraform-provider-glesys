package glesys

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider - Setup new Terraform Provider resource
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// specify what is needed to configure the provider.
			"userid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GLESYS_USERID", nil),
				Description: "UserId for the Glesys API.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GLESYS_TOKEN", nil),
				Description: "User token for the Glesys API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"glesys_dnsdomain":             resourceGlesysDNSDomain(),
			"glesys_dnsdomain_record":      resourceGlesysDNSDomainRecord(),
			"glesys_loadbalancer":          resourceGlesysLoadBalancer(),
			"glesys_loadbalancer_backend":  resourceGlesysLoadBalancerBackend(),
			"glesys_loadbalancer_frontend": resourceGlesysLoadBalancerFrontend(),
			"glesys_loadbalancer_target":   resourceGlesysLoadBalancerTarget(),
			"glesys_network":               resourceGlesysNetwork(),
			"glesys_networkadapter":        resourceGlesysNetworkAdapter(),
			"glesys_server":                resourceGlesysServer(),
		},
		// this will be used to configure the client to communicate with the API
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		UserID: d.Get("userid").(string),
		Token:  d.Get("token").(string),
	}
	return config.Client()
}
