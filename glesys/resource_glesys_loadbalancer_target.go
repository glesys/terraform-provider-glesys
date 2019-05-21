package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysLoadBalancerTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysLoadBalancerTargetCreate,
		Read:   resourceGlesysLoadBalancerTargetRead,
		Update: resourceGlesysLoadBalancerTargetUpdate,
		Delete: resourceGlesysLoadBalancerTargetDelete,

		Schema: map[string]*schema.Schema{
			"backend": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},

			"loadbalancerid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"targetip": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"weight": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceGlesysLoadBalancerTargetCreate(d *schema.ResourceData, m interface{}) error {
	// Add target to glesys_loadbalancer_backend resource
	client := m.(*glesys.Client)

	params := glesys.AddTargetParams{
		Backend:  d.Get("backend").(string),
		Name:     d.Get("name").(string),
		Port:     d.Get("port").(int),
		TargetIP: d.Get("targetip").(string),
		Weight:   d.Get("weight").(int),
	}

	loadbalancerID := d.Get("loadbalancerid").(string)

	_, err := client.LoadBalancers.AddTarget(context.Background(), loadbalancerID, params)
	if err != nil {
		return fmt.Errorf("Error creating LoadBalancer Target: %s", err)
	}

	if d.Get("enabled").(bool) == false {
		// Disable the target after creation
		targetParams := glesys.ToggleTargetParams{Name: d.Get("name").(string), Backend: d.Get("backend").(string)}
		_, err := client.LoadBalancers.DisableTarget(context.Background(), loadbalancerID, targetParams)

		if err != nil {
			return fmt.Errorf("Could not disable Target during creation: %s\n", err)
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceGlesysLoadBalancerTargetRead(d, m)
}

func resourceGlesysLoadBalancerTargetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)
	lb, err := client.LoadBalancers.Details(context.Background(), loadbalancerid)
	if err != nil {
		fmt.Errorf("Loadbalancer not found: %s\n", err)
		d.SetId("")
		return nil
	}

	// iterate over all backends && targets of the loadbalancer
	for _, n := range lb.BackendsList {
		if n.Name == d.Get("backend").(string) {
			for _, t := range n.Targets {
				if t.Name == d.Get("name").(string) {
					d.Set("enabled", t.Enabled)
					d.Set("port", t.Port)
					d.Set("status", t.Status)
					d.Set("targetip", t.TargetIP)
					d.Set("weight", t.Weight)
				}
			}
		} else {
			fmt.Errorf("LoadBalancer Target not found: %s\n", d.Get("name").(string))
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceGlesysLoadBalancerTargetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.EditTargetParams{
		Backend: d.Get("backend").(string),
		Name:    d.Get("name").(string),
	}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	// If changed, toggle the enabled state of the target.
	if d.HasChange("enabled") {
		currentState, _ := d.GetChange("enabled")

		toggleParams := glesys.ToggleTargetParams{
			Backend: d.Get("backend").(string),
			Name:    d.Get("name").(string),
		}

		if currentState == true {
			_, err := client.LoadBalancers.DisableTarget(context.Background(), loadbalancerid, toggleParams)
			if err != nil {
				return fmt.Errorf("Error toggling LoadBalancer Target from: %s - %s\n", currentState, err)
			}
		} else {
			_, err := client.LoadBalancers.EnableTarget(context.Background(), loadbalancerid, toggleParams)
			if err != nil {
				return fmt.Errorf("Error toggling LoadBalancer Target from: %s - %s\n", currentState, err)
			}
		}
	}

	if d.HasChange("port") {
		params.Port = d.Get("port").(int)
	}

	if d.HasChange("targetip") {
		params.TargetIP = d.Get("targetip").(string)
	}

	if d.HasChange("weight") {
		params.Weight = d.Get("weight").(int)
	}

	_, err := client.LoadBalancers.EditTarget(context.Background(), loadbalancerid, params)

	if err != nil {
		return fmt.Errorf("Error updating LoadBalancer Target: %s", err)
	}

	return resourceGlesysLoadBalancerTargetRead(d, m)
}

func resourceGlesysLoadBalancerTargetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.RemoveTargetParams{
		Backend: d.Get("backend").(string),
		Name:    d.Get("name").(string),
	}

	err := client.LoadBalancers.RemoveTarget(context.Background(), loadbalancerid, params)
	if err != nil {
		return fmt.Errorf("Error deleting LoadBalancer Target: %s", err)
	}

	d.SetId("")
	return nil
}
