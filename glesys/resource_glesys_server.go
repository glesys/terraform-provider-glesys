package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysServerCreate,
		ReadContext:   resourceGlesysServerRead,
		UpdateContext: resourceGlesysServerUpdate,
		DeleteContext: resourceGlesysServerDelete,

		Description: "Create a new GleSYS virtual server.",

		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Description: "Server network adapter bandwidth",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"campaigncode": {
				Description: "Campaigncode used during creation for possible discount",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cpu": {
				Description: "Server CPU cores count",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"datacenter": {
				Description: "Server datacenter placement",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Server description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"hostname": {
				Description:      "Server hostname",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: IgnoreCase,
			},
			"ipv4_address": {
				Description: "Server IPv4 address, set `None` to disable IP allocation",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"ipv6_address": {
				Description: "Server IPv6 address, set `None` to disable IP allocation",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"memory": {
				Description: "Server RAM setting",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"password": {
				Description: "Server root password, VMware only",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"platform": {
				Description: "Server virtualisation platform, `KVM` or `VMware`",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"publickey": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"storage": {
				Description: "Server disk space",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"template": {
				Description: "Server OS template",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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
							Description: "User SSH key(s), as a list. '[\"ssh-rsa abc...\", \"ssh-rsa foo...\"]'",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
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

func resourceGlesysServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Setup client to the API
	client := m.(*glesys.Client)

	// Setup server parameters
	srv := buildServerParamStruct(d)

	if srv.Platform == "KVM" {
		usersList, err := expandUsers(d.Get("user").(*schema.Set).List())
		if err != nil {
			return diag.Errorf("Error when expanding users: %s", err)
		}
		srv.Users = usersList
	}

	host, err := client.Servers.Create(ctx, *srv)

	if err != nil {
		return diag.Errorf("error creating server: %+v", err)
	}

	// Set the resource Id to server ID
	d.SetId((*host).ID)

	return resourceGlesysServerRead(ctx, d, m)
}

func getTemplate(original string, srv *glesys.ServerDetails) string {
	for _, tag := range srv.InitialTemplate.CurrentTags {
		if tag == original {
			return original
		}
	}
	return srv.Template
}

func resourceGlesysServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// fetch updates about the resource
	srv, err := client.Servers.Details(context.Background(), d.Id())
	if err != nil {
		diag.Errorf("server not found: %s", err)
		d.SetId("")
		return nil
	}

	// Workaround for the API not returning the correct Bandwith value for KVM servers
	if srv.Platform != "KVM" {
		d.Set("bandwidth", srv.Bandwidth)
	}

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

func resourceGlesysServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error updating instance: %s", err)
	}

	return nil
}

func resourceGlesysServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// destroy the server, don't keep the ip.
	err := client.Servers.Destroy(ctx, d.Id(), glesys.DestroyServerParams{KeepIP: false})

	if err != nil {
		return diag.Errorf("Error deleting instance (%s): %s", d.Id(), err)
	}

	return nil
}
