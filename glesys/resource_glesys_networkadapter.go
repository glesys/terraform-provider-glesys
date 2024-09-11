package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysNetworkAdapter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysNetworkAdapterCreate,
		ReadContext:   resourceGlesysNetworkAdapterRead,
		UpdateContext: resourceGlesysNetworkAdapterUpdate,
		DeleteContext: resourceGlesysNetworkAdapterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Create a networkadapter attached to a VMware server.",

		Schema: map[string]*schema.Schema{
			"adaptertype": {
				Description: "`VMXNET 3` (default) or `E1000`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"bandwidth": {
				Description: "adapter bandwidth",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "Network Adapter name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"networkid": {
				Description: "Network ID to connect to. Defaults to `internet`.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"serverid": {
				Description: "Server ID to connect the adapter to",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceGlesysNetworkAdapterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreateNetworkAdapterParams{
		AdapterType: d.Get("adaptertype").(string),
		Bandwidth:   d.Get("bandwidth").(int),
		NetworkID:   d.Get("networkid").(string),
		ServerID:    d.Get("serverid").(string),
	}

	networkadapter, err := client.NetworkAdapters.Create(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating adapter: %s", err)
	}

	d.SetId(networkadapter.ID)
	return resourceGlesysNetworkAdapterRead(ctx, d, m)
}

func resourceGlesysNetworkAdapterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	networkadapter, err := client.NetworkAdapters.Details(ctx, d.Id())
	if err != nil {
		diag.Errorf("adapter not found: %s", err)
		d.SetId("")
		return nil
	}

	d.Set("adaptertype", networkadapter.AdapterType)
	d.Set("bandwidth", networkadapter.Bandwidth)
	d.Set("name", networkadapter.Name)
	d.Set("networkid", networkadapter.NetworkID)
	d.Set("serverid", networkadapter.ServerID)

	return nil
}

func resourceGlesysNetworkAdapterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditNetworkAdapterParams{}

	// if d.HasChange("adaptertype") {
	//  	params.AdapterType = d.Get("adaptertype").(string)
	// }
	if d.HasChange("bandwidth") {
		params.Bandwidth = d.Get("bandwidth").(int)
	}
	if d.HasChange("networkid") {
		params.NetworkID = d.Get("networkid").(string)
	}

	_, err := client.NetworkAdapters.Edit(ctx, d.Id(), params)
	if err != nil {
		return diag.Errorf("Error updating adapter: %s", err)
	}
	return resourceGlesysNetworkAdapterRead(ctx, d, m)
}

func resourceGlesysNetworkAdapterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.NetworkAdapters.Destroy(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting adapter: %s", err)
	}
	d.SetId("")
	return nil
}
