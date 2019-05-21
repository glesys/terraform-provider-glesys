package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysLoadbalancerCreate,
		Read:   resourceGlesysLoadbalancerRead,
		Update: resourceGlesysLoadbalancerUpdate,
		Delete: resourceGlesysLoadbalancerDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"iplist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"blacklist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGlesysLoadbalancerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.CreateLoadbalancerParams{
		DataCenter: d.Get("datacenter").(string),
		Name:       d.Get("name").(string),
	}

	loadbalancer, err := client.Loadbalancers.Create(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating loadbalancer: %s", err)
	}

	// Set the Id to loadbalancer.ID
	d.SetId((*loadbalancer).ID)

	return resourceGlesysLoadbalancerRead(d, m)
}

func resourceGlesysLoadbalancerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancer, err := client.Loadbalancers.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("Loadbalancer not found: %s\n", err)
		d.SetId("")
		return nil
	}

	var ip_addresses []string
	var blacklist_ips []string
	for i := range (*loadbalancer).IPList {
		ip_addresses = append(ip_addresses, (*loadbalancer).IPList[i].Address)

	}

	for i := range (*loadbalancer).Blacklists {
		blacklist_ips = append(blacklist_ips, (*loadbalancer).Blacklists[i])
	}

	d.Set("datacenter", loadbalancer.DataCenter)
	d.Set("name", loadbalancer.Name)
	d.Set("iplist", ip_addresses)
	d.Set("blacklist", blacklist_ips)

	return nil
}

func resourceGlesysLoadbalancerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.EditLoadbalancerParams{}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	_, err := client.Loadbalancers.Edit(context.Background(), d.Id(), params)
	if err != nil {
		return fmt.Errorf("Error updating loadbalancer: %s", err)
	}

	return resourceGlesysLoadbalancerRead(d, m)
}

func resourceGlesysLoadbalancerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	err := client.Loadbalancers.Destroy(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting loadbalancer: %s", err)
	}
	d.SetId("")
	return nil
}
