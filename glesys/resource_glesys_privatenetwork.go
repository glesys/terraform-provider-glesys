package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysPrivateNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysPrivateNetworkCreate,
		ReadContext:   resourceGlesysPrivateNetworkRead,
		UpdateContext: resourceGlesysPrivateNetworkUpdate,
		DeleteContext: resourceGlesysPrivateNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Create a PrivateNetwork resource.",

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "PrivateNetwork name",
				Type:        schema.TypeString,
				Required:    true,
			},

			"ipv6aggregate": {
				Description: "IPv6Aggregate for the PrivateNetwork.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGlesysPrivateNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	name := d.Get("name").(string)
	network, err := client.PrivateNetworks.Create(ctx, name)
	if err != nil {
		return diag.Errorf("Error adding privatenetwork %s: %v", name, err)
	}

	// Set the Id to network.ID
	d.SetId(network.ID)

	return resourceGlesysPrivateNetworkRead(ctx, d, m)
}

func resourceGlesysPrivateNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	network, err := client.PrivateNetworks.Details(ctx, d.Id())

	if err != nil {
		diag.Errorf("privatenetwork not found: %v", err)
		d.SetId("")
		return nil
	}

	d.Set("ipv6aggregate", network.IPv6Aggregate)
	d.Set("name", network.Name)

	return nil
}

func resourceGlesysPrivateNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditPrivateNetworkParams{ID: d.Id()}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	_, err := client.PrivateNetworks.Edit(ctx, params)
	if err != nil {
		return diag.Errorf("Error updating privatenetwork: %v", err)
	}

	return resourceGlesysPrivateNetworkRead(ctx, d, m)
}

func resourceGlesysPrivateNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.PrivateNetworks.Destroy(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting privatenetwork: %v", err)
	}
	d.SetId("")
	return nil
}
