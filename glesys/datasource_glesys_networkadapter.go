package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGlesysNetworkAdapter() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about a NetworkAdapter associated with ServerID.",

		ReadContext: dataSourceGlesysNetworkAdapterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "networkadapter ID.",
			},
			"adaptertype": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "networkadapter adaptertype. (VMware)",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "networkadapter bandwidth.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "networkadapter name.",
			},
		},
	}
}

func dataSourceGlesysNetworkAdapterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	var na *glesys.NetworkAdapter
	if adapterID, ok := d.GetOk("id"); ok {
		nic, err := client.NetworkAdapters.Details(ctx, adapterID.(string))
		if err != nil {
			return diag.Errorf("Error retrieving networkadapter: %s", err)
		}
		na = nic
	}

	d.SetId(na.ID)
	d.Set("adaptertype", na.AdapterType)
	d.Set("bandwidth", na.Bandwidth)
	d.Set("name", na.Name)

	return nil
}
