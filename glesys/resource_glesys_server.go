package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: IgnoreCase,
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

			"user": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"publickeys": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
					},
				},
			},
		},
	}
}

func expandUsers(config []interface{}) ([]glesys.User, error) {
	users := make([]glesys.User, 0, len(config))

	for _, rawUser := range config {
		user := rawUser.(map[string]interface{})
		var pks []string

		u := glesys.User{
			Username: user["username"].(string),
			Password: user["password"].(string),
		}

		// Append publickeys to the slice for PublicKeys
		for _, pk := range user["publickeys"].([]interface{}) {
			pks = append(pks, pk.(string))
		}
		u.PublicKeys = pks

		users = append(users, u)
	}

	return users, nil
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

	if srv.Platform == "KVM" {
		usersList, err := expandUsers(d.Get("user").(*schema.Set).List())
		if err != nil {
			return err
		}
		srv.Users = usersList
	}

	host, err := client.Servers.Create(context.Background(), *srv)

	if err != nil {
		return fmt.Errorf("error creating server: %+v", err)
	}

	// Set the resource Id to server ID
	d.SetId((*host).ID)
	return resourceGlesysServerRead(d, m)
}

func getTemplate(original string, srv *glesys.ServerDetails) string {
	for _, tag := range srv.InitialTemplate.CurrentTags {
		if tag == original {
			return original
		}
	}
	return srv.Template
}

func resourceGlesysServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	// fetch updates about the resource
	srv, err := client.Servers.Details(context.Background(), d.Id())
	if err != nil {
		fmt.Errorf("server not found: %s", err)
		d.SetId("")
		return nil
	}

	d.Set("bandwidth", srv.Bandwidth)
	d.Set("cpu", srv.CPU)
	d.Set("datacenter", srv.DataCenter)
	d.Set("description", srv.Description)
	d.Set("hostname", srv.Hostname)
	for i := range srv.IPList {
		if srv.IPList[i].Version == 4 {
			d.Set("ipv4_address", srv.IPList[i].Address)
		}
		if srv.IPList[i].Version == 6 {
			d.Set("ipv6_address", srv.IPList[i].Address)
		}
	}
	d.Set("memory", srv.Memory)
	d.Set("platform", srv.Platform)
	d.Set("storage", srv.Storage)
	d.Set("template", getTemplate(d.Get("template").(string), srv))

	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": d.Get("ipv4_address").(string),
	})

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
