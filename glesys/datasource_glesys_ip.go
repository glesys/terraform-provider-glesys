package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGlesysIP() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about a reserved IP address in your Glesys Project.",

		ReadContext: dataSourceGlesysIPRead,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address",
			},
			"broadcast": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP Broadcast address.",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP Datacenter association.",
			},
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default gateway.",
			},
			"netmask": {
				Description: "IP Netmask, IPv4: NN.NN.NN.NN, IPv6: /nn",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"platform": {
				Description: "IP Associated platform.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ptr": {
				Description: "IP PTR.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceGlesysIPRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*glesys.Client)

	var ip *glesys.IP
	ip, err := client.IPs.Details(ctx, d.Get("address").(string))
	if err != nil {
		diag.Errorf("IP not found: %s", err)
		d.SetId("")
		return nil
	}

	d.Set("address", ip.Address)
	d.Set("broadcast", ip.Broadcast)
	d.Set("datacenter", ip.DataCenter)
	d.Set("gateway", ip.Gateway)
	d.Set("netmask", ip.Netmask)
	d.Set("platform", ip.Platform)
	d.Set("ptr", ip.PTR)
	d.SetId(ip.Address)

	return nil
}
