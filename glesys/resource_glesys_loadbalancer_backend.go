package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysLoadBalancerBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysLoadBalancerBackendCreate,
		Read:   resourceGlesysLoadBalancerBackendRead,
		Update: resourceGlesysLoadBalancerBackendUpdate,
		Delete: resourceGlesysLoadBalancerBackendDelete,

		Schema: map[string]*schema.Schema{
			"connecttimeout": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"loadbalancerid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"responsetimeout": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"stickysessions": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGlesysLoadBalancerBackendCreate(d *schema.ResourceData, m interface{}) error {
	// Add frontend to glesys_loadbalancer resource
	client := m.(*glesys.Client)

	params := glesys.AddBackendParams{
		ConnectTimeout:  d.Get("connecttimeout").(int),
		Mode:            d.Get("mode").(string),
		Name:            d.Get("name").(string),
		ResponseTimeout: d.Get("responsetimeout").(int),
		StickySession:   d.Get("stickysessions").(string),
	}

	loadbalancerID := d.Get("loadbalancerid").(string)

	_, err := client.LoadBalancers.AddBackend(context.Background(), loadbalancerID, params)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Backend: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceGlesysLoadBalancerBackendRead(d, m)
}

func resourceGlesysLoadBalancerBackendRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)
	lb, err := client.LoadBalancers.Details(context.Background(), loadbalancerid)
	if err != nil {
		fmt.Errorf("Loadbalancer not found: %s\n", err)
		d.SetId("")
		return nil
	}

	for _, n := range lb.BackendsList {
		if n.Name == d.Get("name").(string) {
			d.Set("mode", n.Mode)
			d.Set("connecttimeout", n.ConnectTimeout)
			d.Set("responsetimeout", n.ResponseTimeout)
			d.Set("status", n.Status)

			var targets []string
			for _, t := range n.Targets {
				targets = append(targets, t.Name)
			}
			d.Set("targets", targets)
		}
	}

	return nil
}

func resourceGlesysLoadBalancerBackendUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.EditBackendParams{
		Name: d.Get("name").(string),
	}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	if d.HasChange("connecttimeout") {
		params.ConnectTimeout = d.Get("connecttimeout").(int)
	}

	if d.HasChange("mode") {
		params.Mode = d.Get("mode").(string)
	}

	if d.HasChange("responsetimeout") {
		params.ResponseTimeout = d.Get("responsetimeout").(int)
	}

	if d.HasChange("stickysessions") {
		params.StickySession = d.Get("stickysessions").(string)
	}

	_, err := client.LoadBalancers.EditBackend(context.Background(), loadbalancerid, params)
	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer Backend: %s", err)
	}

	return resourceGlesysLoadBalancerBackendRead(d, m)
}

func resourceGlesysLoadBalancerBackendDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.RemoveBackendParams{
		Name: d.Get("name").(string),
	}

	err := client.LoadBalancers.RemoveBackend(context.Background(), loadbalancerid, params)
	if err != nil {
		return fmt.Errorf("Error deleting LoadBalancer Backend: %s", err)
	}

	d.SetId("")
	return nil
}
