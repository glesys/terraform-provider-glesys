package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGlesysNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about a Network associated with your GleSYS Project.",

		ReadContext: dataSourceGlesysNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "network ID.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network description.",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network datacenter.",
			},
			"public": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network public, yes/no.",
			},
		},
	}
}

func dataSourceGlesysNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	var network *glesys.Network
	if networkid, ok := d.GetOk("id"); ok {
		net, err := client.Networks.Details(ctx, networkid.(string))
		if err != nil {
			return diag.Errorf("Error retrieving network: %s", err)
		}
		network = net
	}

	d.SetId(network.ID)
	d.Set("description", network.Description)
	d.Set("datacenter", network.DataCenter)
	d.Set("public", network.Public)

	return nil
}
