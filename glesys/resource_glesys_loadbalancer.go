package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysLoadBalancerCreate,
		Read:   resourceGlesysLoadBalancerRead,
		Update: resourceGlesysLoadBalancerUpdate,
		Delete: resourceGlesysLoadBalancerDelete,

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

func resourceGlesysLoadBalancerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.CreateLoadBalancerParams{
		DataCenter: d.Get("datacenter").(string),
		Name:       d.Get("name").(string),
	}

	loadbalancer, err := client.LoadBalancers.Create(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating loadbalancer: %s", err)
	}

	// Set the Id to loadbalancer.ID
	d.SetId((*loadbalancer).ID)

	return resourceGlesysLoadBalancerRead(d, m)
}

func resourceGlesysLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancer, err := client.LoadBalancers.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("loadbalancer not found: %s", err)
		d.SetId("")
		return nil
	}

	var ipAddresses []string
	var blacklistIps []string
	for i := range (*loadbalancer).IPList {
		ipAddresses = append(ipAddresses, (*loadbalancer).IPList[i].Address)

	}

	for i := range (*loadbalancer).Blacklists {
		blacklistIps = append(blacklistIps, (*loadbalancer).Blacklists[i])
	}

	d.Set("datacenter", loadbalancer.DataCenter)
	d.Set("name", loadbalancer.Name)
	d.Set("iplist", ipAddresses)
	d.Set("blacklist", blacklistIps)

	return nil
}

func resourceGlesysLoadBalancerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.EditLoadBalancerParams{}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	_, err := client.LoadBalancers.Edit(context.Background(), d.Id(), params)
	if err != nil {
		return fmt.Errorf("Error updating loadbalancer: %s", err)
	}

	return resourceGlesysLoadBalancerRead(d, m)
}

func resourceGlesysLoadBalancerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	err := client.LoadBalancers.Destroy(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting loadbalancer: %s", err)
	}
	d.SetId("")
	return nil
}
