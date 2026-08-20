package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysLoadBalancerBackend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysLoadBalancerBackendCreate,
		ReadContext:   resourceGlesysLoadBalancerBackendRead,
		UpdateContext: resourceGlesysLoadBalancerBackendUpdate,
		DeleteContext: resourceGlesysLoadBalancerBackendDelete,

		Description: "LoadBalancer Backend for a glesys_loadbalancer",

		Schema: map[string]*schema.Schema{
			"connecttimeout": {
				Description: "Connection timeout to backend target. `milliseconds`",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"loadbalancerid": {
				Description: "LoadBalancer ID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"responsetimeout": {
				Description: "Connection response timeout. `milliseconds`",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"name": {
				Description: "Backend name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"mode": {
				Description: "Backend mode. `TCP`, `HTTP`.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"stickysessions": {
				Description: "Enable backend sticky sessions. `true`, `false`, `yes`, `no`.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"status": {
				Description: "Backend status. `UP` when targets are reachable and `DOWN` when no targets are reachable.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"targets": {
				Description: "Backend targets. Computed by LoadBalancer Targets setting the `backend` parameter.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGlesysLoadBalancerBackendCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	_, err := client.LoadBalancers.AddBackend(ctx, loadbalancerID, params)
	if err != nil {
		return diag.Errorf("Error creating LoadBalancer Backend: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceGlesysLoadBalancerBackendRead(ctx, d, m)
}

func resourceGlesysLoadBalancerBackendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)
	lb, err := client.LoadBalancers.Details(ctx, loadbalancerid)
	if err != nil {
		diag.Errorf("loadbalancer not found: %s", err)
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

func resourceGlesysLoadBalancerBackendUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	_, err := client.LoadBalancers.EditBackend(ctx, loadbalancerid, params)
	if err != nil {
		return diag.Errorf("Error updating LoadBalancer Backend: %s", err)
	}

	return resourceGlesysLoadBalancerBackendRead(ctx, d, m)
}

func resourceGlesysLoadBalancerBackendDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.RemoveBackendParams{
		Name: d.Get("name").(string),
	}

	err := client.LoadBalancers.RemoveBackend(ctx, loadbalancerid, params)
	if err != nil {
		return diag.Errorf("Error deleting LoadBalancer Backend: %s", err)
	}

	d.SetId("")
	return nil
}
