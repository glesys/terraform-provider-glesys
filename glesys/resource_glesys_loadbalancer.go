package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysLoadBalancerCreate,
		ReadContext:   resourceGlesysLoadBalancerRead,
		UpdateContext: resourceGlesysLoadBalancerUpdate,
		DeleteContext: resourceGlesysLoadBalancerDelete,

		Description: "Create a LoadBalancer",

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Description: "LoadBalancer datacenter. `Falkenberg`, `Stockholm`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "LoadBalancer name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"iplist": {
				Description: "IPs set on the LoadBalancer.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"blacklist": {
				Description: "**DEPRECATED** Use blocklist instead.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Deprecated:  "use blocklist instead",
			},
			"blocklist": {
				Description: "LoadBalancer blocklist. List of IPs: `[\"a.b.c.d\",\"x.y.z.w\"]`",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGlesysLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreateLoadBalancerParams{
		DataCenter: d.Get("datacenter").(string),
		Name:       d.Get("name").(string),
	}

	loadbalancer, err := client.LoadBalancers.Create(context.Background(), params)
	if err != nil {
		return diag.Errorf("Error creating loadbalancer: %s", err)
	}

	// Set the Id to loadbalancer.ID
	d.SetId(loadbalancer.ID)

	return resourceGlesysLoadBalancerRead(ctx, d, m)
}

func resourceGlesysLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancer, err := client.LoadBalancers.Details(context.Background(), d.Id())
	if err != nil {
		diag.Errorf("loadbalancer not found: %s", err)
		d.SetId("")
		return nil
	}

	var ipAddresses []string
	var blocklistIps []string
	for i := range loadbalancer.IPList {
		ipAddresses = append(ipAddresses, loadbalancer.IPList[i].Address)
	}

	for i := range loadbalancer.Blocklists {
		blocklistIps = append(blocklistIps, loadbalancer.Blocklists[i])
	}

	d.Set("datacenter", loadbalancer.DataCenter)
	d.Set("name", loadbalancer.Name)
	d.Set("iplist", ipAddresses)
	d.Set("blacklist", blocklistIps)
	d.Set("blocklist", blocklistIps)

	return nil
}

func resourceGlesysLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditLoadBalancerParams{}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	_, err := client.LoadBalancers.Edit(context.Background(), d.Id(), params)
	if err != nil {
		return diag.Errorf("Error updating loadbalancer: %s", err)
	}

	return resourceGlesysLoadBalancerRead(ctx, d, m)
}

func resourceGlesysLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.LoadBalancers.Destroy(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting loadbalancer: %s", err)
	}
	d.SetId("")
	return nil
}
