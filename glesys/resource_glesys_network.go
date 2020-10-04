package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysNetworkCreate,
		Read:   resourceGlesysNetworkRead,
		Update: resourceGlesysNetworkUpdate,
		Delete: resourceGlesysNetworkDelete,

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

func resourceGlesysNetworkCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.CreateNetworkParams{
		DataCenter:  d.Get("datacenter").(string),
		Description: d.Get("description").(string),
	}

	network, err := client.Networks.Create(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating network: %s", err)
	}

	// Set the Id to network.ID
	d.SetId((*network).ID)
	return resourceGlesysNetworkRead(d, m)
}

func resourceGlesysNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	network, err := client.Networks.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("network not found: %s", err)
		d.SetId("")
		return nil
	}

	d.Set("datacenter", network.DataCenter)
	d.Set("description", network.Description)
	d.Set("public", network.Public)

	return nil
}

func resourceGlesysNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.EditNetworkParams{}

	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
	}

	_, err := client.Networks.Edit(context.Background(), d.Id(), params)
	if err != nil {
		return fmt.Errorf("Error updating network: %s", err)
	}
	return resourceGlesysNetworkRead(d, m)
}

func resourceGlesysNetworkDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	// TODO: check if network is used before deletion.
	// remove networkadapter, then network
	err := client.Networks.Destroy(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting network: %s", err)
	}
	d.SetId("")
	return nil
}
