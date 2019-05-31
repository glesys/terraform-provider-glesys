package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysNetworkAdapter() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysNetworkAdapterCreate,
		Read:   resourceGlesysNetworkAdapterRead,
		Update: resourceGlesysNetworkAdapterUpdate,
		Delete: resourceGlesysNetworkAdapterDelete,

		Schema: map[string]*schema.Schema{
			"adaptertype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"networkid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"serverid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGlesysNetworkAdapterCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.CreateNetworkAdapterParams{
		AdapterType: d.Get("adaptertype").(string),
		Bandwidth:   d.Get("bandwidth").(int),
		NetworkID:   d.Get("networkid").(string),
		ServerID:    d.Get("serverid").(string),
	}

	networkadapter, err := client.NetworkAdapters.Create(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating adapter: %s", err)
	}

	d.SetId(networkadapter.ID)
	return resourceGlesysNetworkAdapterRead(d, m)
}

func resourceGlesysNetworkAdapterRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	networkadapter, err := client.NetworkAdapters.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("Adapter not found: %s\n", err)
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

func resourceGlesysNetworkAdapterUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.EditNetworkAdapterParams{}

	//if d.HasChange("adaptertype") {
	//	params.AdapterType = d.Get("adaptertype").(string)
	//}
	if d.HasChange("bandwidth") {
		params.Bandwidth = d.Get("bandwidth").(int)
	}
	if d.HasChange("networkid") {
		params.NetworkID = d.Get("networkid").(string)
	}

	_, err := client.NetworkAdapters.Edit(context.Background(), d.Id(), params)
	if err != nil {
		return fmt.Errorf("Error updating adapter: %s", err)
	}
	return resourceGlesysNetworkAdapterRead(d, m)
}

func resourceGlesysNetworkAdapterDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	err := client.NetworkAdapters.Destroy(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting adapter: %s", err)
	}
	d.SetId("")
	return nil
}
