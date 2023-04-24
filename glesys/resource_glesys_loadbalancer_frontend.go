package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysLoadBalancerFrontend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysLoadBalancerFrontendCreate,
		ReadContext:   resourceGlesysLoadBalancerFrontendRead,
		UpdateContext: resourceGlesysLoadBalancerFrontendUpdate,
		DeleteContext: resourceGlesysLoadBalancerFrontendDelete,

		Description: "Create a LoadBalancer Frontend for a `glesys_loadbalancer`.",

		Schema: map[string]*schema.Schema{
			"backend": {
				Description: "LoadBalancer Backend name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"clienttimeout": {
				Description: "Client connection timeout. `milliseconds`",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"loadbalancerid": {
				Description: "LoadBalancer to associate the Frontend to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"maxconnections": {
				Description: "Maximum number of connections allowed.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"name": {
				Description: "Frontend name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"port": {
				Description: "Listen port.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			"sslcertificate": {
				Description: "Certificate bundle to use for terminating TLS connections.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"status": {
				Description: "Frontend status.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGlesysLoadBalancerFrontendCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	_, err := client.LoadBalancers.AddFrontend(ctx, loadbalancerID, params)
	if err != nil {
		return diag.Errorf("Error creating LoadBalancer Frontend: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceGlesysLoadBalancerFrontendRead(ctx, d, m)
}

func resourceGlesysLoadBalancerFrontendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)
	lb, err := client.LoadBalancers.Details(ctx, loadbalancerid)
	if err != nil {
		diag.Errorf("loadbalancer not found: %s", err)
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

func resourceGlesysLoadBalancerFrontendUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	_, err := client.LoadBalancers.EditFrontend(ctx, loadbalancerid, params)
	if err != nil {
		return diag.Errorf("Error updating LoadBalancer Frontend: %s", err)
	}

	return resourceGlesysLoadBalancerFrontendRead(ctx, d, m)
}

func resourceGlesysLoadBalancerFrontendDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.RemoveFrontendParams{
		Name: d.Get("name").(string),
	}

	err := client.LoadBalancers.RemoveFrontend(ctx, loadbalancerid, params)
	if err != nil {
		return diag.Errorf("Error deleting LoadBalancer Frontend: %s", err)
	}

	d.SetId("")
	return nil
}
