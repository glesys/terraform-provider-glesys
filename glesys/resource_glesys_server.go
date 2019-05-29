package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysServerCreate,
		Read:   resourceGlesysServerRead,
		Update: resourceGlesysServerUpdate,
		Delete: resourceGlesysServerDelete,

		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"campaigncode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ipv4_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"publickey": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func buildServerParamStruct(d *schema.ResourceData) *glesys.CreateServerParams {
	opts := glesys.CreateServerParams{
		Bandwidth:    d.Get("bandwidth").(int),
		CampaignCode: d.Get("campaigncode").(string),
		CPU:          d.Get("cpu").(int),
		DataCenter:   d.Get("datacenter").(string),
		Description:  d.Get("description").(string),
		Hostname:     d.Get("hostname").(string),
		IPv4:         d.Get("ipv4_address").(string),
		IPv6:         d.Get("ipv6_address").(string),
		Memory:       d.Get("memory").(int),
		Password:     d.Get("password").(string),
		Platform:     d.Get("platform").(string),
		PublicKey:    d.Get("publickey").(string),
		Storage:      d.Get("storage").(int),
		Template:     d.Get("template").(string),
	}.WithDefaults()

	return &opts
}

func resourceGlesysServerCreate(d *schema.ResourceData, m interface{}) error {
	// Setup client to the API
	client := m.(*glesys.Client)

	// Setup server parameters
	srv := buildServerParamStruct(d)

	host, err := client.Servers.Create(context.Background(), *srv)

	if err != nil {
		return fmt.Errorf("Error creating server: %+v\n", err)
	}

	// Set the resource Id to server ID
	d.SetId((*host).ID)
	return resourceGlesysServerRead(d, m)
}

func resourceGlesysServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	// fetch updates about the resource
	srv, err := client.Servers.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("Server not found: %s\n", err)
		d.SetId("")
		return nil
	}

	d.Set("bandwidth", srv.Bandwidth)
	d.Set("cpu", srv.CPU)
	d.Set("datacenter", srv.DataCenter)
	d.Set("description", srv.Description)
	d.Set("hostname", srv.Hostname)
	for i := range srv.IPList {
		if srv.IPList[i].IsIPv4() {
			d.Set("ipv4_address", srv.IPList[i].Address)
		}
		if srv.IPList[i].IsIPv6() {
			d.Set("ipv6_address", srv.IPList[i].Address)
		}
	}
	d.Set("memory", srv.Memory)
	d.Set("platform", srv.Platform)
	d.Set("storage", srv.Storage)
	d.Set("template", srv.Template)

	return nil
}

func resourceGlesysServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	opts := glesys.EditServerParams{}

	if d.HasChange("cpu") {
		opts.CPU = d.Get("cpu").(int)
	}
	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}
	if d.HasChange("hostname") {
		opts.Hostname = d.Get("hostname").(string)
	}
	if d.HasChange("memory") {
		opts.Memory = d.Get("memory").(int)
	}
	if d.HasChange("storage") {
		opts.Storage = d.Get("storage").(int)
	}
	_, err := client.Servers.Edit(context.Background(), d.Id(), opts)
	if err != nil {
		return fmt.Errorf("Error updating instance: %s", err)
	}

	return nil
}

func resourceGlesysServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	// destroy the server, don't keep the ip.
	err := client.Servers.Destroy(context.Background(), d.Id(), glesys.DestroyServerParams{KeepIP: false})

	if err != nil {
		return fmt.Errorf("Error deleting instance: %s", err)
	}

	return nil
}
