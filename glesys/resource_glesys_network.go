package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysNetworkCreate,
		ReadContext:   resourceGlesysNetworkRead,
		UpdateContext: resourceGlesysNetworkUpdate,
		DeleteContext: resourceGlesysNetworkDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
			},
		},
	}
}

func resourceGlesysNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreateNetworkParams{
		DataCenter:  d.Get("datacenter").(string),
		Description: d.Get("description").(string),
	}

	network, err := client.Networks.Create(context.Background(), params)
	if err != nil {
		return diag.Errorf("Error creating network: %s", err)
	}

	// Set the Id to network.ID
	d.SetId(network.ID)
	return resourceGlesysNetworkRead(ctx, d, m)
}

func resourceGlesysNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	network, err := client.Networks.Details(context.Background(), d.Id())
	if err != nil {
		diag.Errorf("network not found: %s", err)
		d.SetId("")
		return nil
	}

	d.Set("datacenter", network.DataCenter)
	d.Set("description", network.Description)
	d.Set("public", network.Public)

	return nil
}

func resourceGlesysNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditNetworkParams{}

	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
	}

	_, err := client.Networks.Edit(context.Background(), d.Id(), params)
	if err != nil {
		return diag.Errorf("Error updating network: %s", err)
	}
	return resourceGlesysNetworkRead(ctx, d, m)
}

func resourceGlesysNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// TODO: check if network is used before deletion.
	// remove networkadapter, then network
	err := client.Networks.Destroy(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting network: %s", err)
	}
	d.SetId("")
	return nil
}
