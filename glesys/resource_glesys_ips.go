package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysIps() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysIpsCreate,
		Read:   resourceGlesysIpsRead,
		Update: resourceGlesysIpsUpdate,
		Delete: resourceGlesysIpsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ipaddress": {
				Type:     schema.TypeString,
				Required: true,
			},
			"broadcast": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"lockedtoaccount": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nameservers": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"netmask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"platforms": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ptr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reserved": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"serverid": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"isipversion4": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"isipversion6": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceGlesysIpsCreate(d *schema.ResourceData, m interface{}) error {

	if d.Get("reserved").(string) == "no" {

		// Setup client to the API
		client := m.(*glesys.Client)

		_, err := client.IPs.Reserve(context.Background(), d.Get("ipaddress").(string))

		if err != nil {
			return fmt.Errorf("Error reserve ip: %+v\n", err)
		}
		return nil
	}

	return resourceGlesysIpsUpdate(d, m)
}

func resourceGlesysIpsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	// fetch updates about the resource
	ips, err := client.IPs.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("IP not found: %s\n", err)
		d.SetId("")
		return nil
	}

	d.Set("ipaddress", ips.Address)
	d.Set("datacenter", ips.DataCenter)
	d.Set("platform", ips.Platform)
	d.Set("broadcast", ips.Broadcast)
	d.Set("cost", ips.Cost)
	d.Set("datacenter", ips.DataCenter)
	d.Set("gateway", ips.Gateway)
	d.Set("lockedtoaccount", ips.LockedToAccount)
	d.Set("nameservers", ips.NameServers)
	d.Set("netmask", ips.Netmask)
	d.Set("platforms", ips.Platforms)
	d.Set("platform", ips.Platform)
	d.Set("ptr", ips.PTR)
	d.Set("reserved", ips.Reserved)
	d.Set("serverid", ips.ServerID)
	d.Set("isipversion4", ips.IsIPv4())
	d.Set("isipversion6", ips.IsIPv6())

	// Set the resource Id to ipaddress
	d.SetId(ips.Address)

	return nil
}

func resourceGlesysIpsUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*glesys.Client)

	opts := glesys.IP{}

	if d.HasChange("ptr") {
		// TODO handle reset ?
		opts.PTR = d.Get("ptr").(string)
	}

	_, err := client.IPs.SetPTR(context.Background(), d.Get("ipaddress").(string), opts.PTR)
	if err != nil {
		return fmt.Errorf("Error updating instance: %s", err)
	}

	return nil
}

func resourceGlesysIpsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	err := client.IPs.Release(context.Background(), d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	return nil
}
