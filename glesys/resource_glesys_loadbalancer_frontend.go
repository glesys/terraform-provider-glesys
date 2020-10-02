package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysLoadBalancerFrontend() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysLoadBalancerFrontendCreate,
		Read:   resourceGlesysLoadBalancerFrontendRead,
		Update: resourceGlesysLoadBalancerFrontendUpdate,
		Delete: resourceGlesysLoadBalancerFrontendDelete,

		Schema: map[string]*schema.Schema{
			"backend": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"clienttimeout": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"loadbalancerid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"maxconnections": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"sslcertificate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGlesysLoadBalancerFrontendCreate(d *schema.ResourceData, m interface{}) error {
	// Add frontend to glesys_loadbalancer resource
	client := m.(*glesys.Client)

	params := glesys.AddFrontendParams{
		Backend:        d.Get("backend").(string),
		ClientTimeout:  d.Get("clienttimeout").(int),
		MaxConnections: d.Get("maxconnections").(int),
		Name:           d.Get("name").(string),
		Port:           d.Get("port").(int),
		SSLCertificate: d.Get("sslcertificate").(string),
	}

	loadbalancerID := d.Get("loadbalancerid").(string)

	_, err := client.LoadBalancers.AddFrontend(context.Background(), loadbalancerID, params)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Frontend: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceGlesysLoadBalancerFrontendRead(d, m)
}

func resourceGlesysLoadBalancerFrontendRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)
	lb, err := client.LoadBalancers.Details(context.Background(), loadbalancerid)
	if err != nil {
		fmt.Errorf("LoadBalancer not found: %s\n", err)
		d.SetId("")
		return nil
	}

	for _, n := range lb.FrontendsList {
		if n.Name == d.Get("name").(string) {
			d.Set("backend", n.Backend)
			d.Set("clienttimeout", n.ClientTimeout)
			d.Set("maxconnections", n.MaxConnections)
			d.Set("port", n.Port)
			d.Set("sslcertificate", n.SSLCertificate)
			d.Set("status", n.Status)
		}
	}

	return nil
}

func resourceGlesysLoadBalancerFrontendUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.EditFrontendParams{
		Name: d.Get("name").(string),
	}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	if d.HasChange("clienttimeout") {
		params.ClientTimeout = d.Get("clienttimeout").(int)
	}

	if d.HasChange("maxconnections") {
		params.MaxConnections = d.Get("maxconnections").(int)
	}

	if d.HasChange("port") {
		params.Port = d.Get("port").(int)
	}

	if d.HasChange("sslcertificate") {
		params.SSLCertificate = d.Get("sslcertificate").(string)
	}

	_, err := client.LoadBalancers.EditFrontend(context.Background(), loadbalancerid, params)
	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer Frontend: %s", err)
	}

	return resourceGlesysLoadBalancerFrontendRead(d, m)
}

func resourceGlesysLoadBalancerFrontendDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.RemoveFrontendParams{
		Name: d.Get("name").(string),
	}

	err := client.LoadBalancers.RemoveFrontend(context.Background(), loadbalancerid, params)
	if err != nil {
		return fmt.Errorf("Error deleting LoadBalancer Frontend: %s", err)
	}

	d.SetId("")
	return nil
}
