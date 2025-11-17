package glesys

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGlesysServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysServerCreate,
		ReadContext:   resourceGlesysServerRead,
		UpdateContext: resourceGlesysServerUpdate,
		DeleteContext: resourceGlesysServerDelete,

		Description: "Create a new Glesys virtual server.",

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
			"cloudconfig": {
				Description: "Cloudconfig used to provision server using a provided cloud-config mustache template.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloudconfigparams": {
				Description: "Cloudconfigparams is used to provide additional parameters to the template in `cloudconfig` using a map. This can be set using a Terraform Local Value.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Description: "Server IPv4 address, set `none` to disable IP allocation",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"ipv6_address": {
				Description: "Server IPv6 address, set `none` to disable IP allocation",
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
			"islocked": {
				Description: "Server locked state",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"isrunning": {
				Description: "Server running state",
				Type:        schema.TypeBool,
				Computed:    true,
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

			"extra_disks": {
				Description: "Disks associated with the server. Use `glesys_server_disk` resource to manage these.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"primary_networkadapter_network": {
				Description: "(VMware) Set the network for the primary network adapter.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"network_adapters": {
				Description: "Network adapters associated with the server. `glesys_networkadapter`",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Networkadapter ID.",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"adaptertype": {
							Description: "`VMXNET 3` (default) or `E1000`",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"bandwidth": {
							Description: "adapter bandwidth",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Network Adapter name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"networkid": {
							Description: "Network ID to connect to. Defaults to `internet`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
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
							Optional: true,
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

			"backups_schedule": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "KVM Server backup schedule definition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"daily", "weekly"}, false),
						},
						"retention": {
							Type:     schema.TypeInt,
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

func expandBackupSchedules(config []interface{}) ([]glesys.ServerBackupSchedule, error) {
	schedules := make([]glesys.ServerBackupSchedule, 0, len(config))

	for _, rawSchedule := range config {
		schedule := rawSchedule.(map[string]interface{})

		s := glesys.ServerBackupSchedule{
			Frequency:            schedule["frequency"].(string),
			Numberofimagestokeep: schedule["retention"].(int),
		}

		schedules = append(schedules, s)
	}
	return schedules, nil
}

func buildServerParamStruct(d *schema.ResourceData) *glesys.CreateServerParams {
	opts := glesys.CreateServerParams{
		Bandwidth:         d.Get("bandwidth").(int),
		CampaignCode:      d.Get("campaigncode").(string),
		CloudConfig:       d.Get("cloudconfig").(string),
		CloudConfigParams: d.Get("cloudconfigparams").(map[string]any),
		CPU:               d.Get("cpu").(int),
		DataCenter:        d.Get("datacenter").(string),
		Description:       d.Get("description").(string),
		Hostname:          d.Get("hostname").(string),
		IPv4:              d.Get("ipv4_address").(string),
		IPv6:              d.Get("ipv6_address").(string),
		Memory:            d.Get("memory").(int),
		Password:          d.Get("password").(string),
		Platform:          d.Get("platform").(string),
		PublicKey:         d.Get("publickey").(string),
		Storage:           d.Get("storage").(int),
		Template:          d.Get("template").(string),
	}.WithDefaults()

	return &opts
}

func resourceGlesysServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Setup client to the API
	client := m.(*glesys.Client)

	// Setup server parameters
	srv := buildServerParamStruct(d)

	backupsList, err := expandBackupSchedules(d.Get("backups_schedule").(*schema.Set).List())
	if err != nil {
		return diag.Errorf("Error when expanding backup schedules: %v", err)
	}
	srv.Backup = backupsList

	// Setup users for server creation
	usersList, err := expandUsers(d.Get("user").(*schema.Set).List())
	if err != nil {
		return diag.Errorf("Error when expanding users: %s", err)
	}
	srv.Users = usersList

	host, err := client.Servers.Create(ctx, *srv)

	if err != nil {
		return diag.Errorf("error creating server: %+v", err)
	}

	// Set the resource Id to server ID
	d.SetId(host.ID)

	if _, err = waitForServerAttribute(ctx, d, "true", []string{"false"}, "isrunning", m); err != nil {
		return diag.Errorf("error while waiting for Server (%s) to be started: %s", d.Id(), err)
	}
	if _, err = waitForServerAttribute(ctx, d, "false", []string{"true"}, "islocked", m); err != nil {
		return diag.Errorf("error while waiting for Server (%s) to be completed: %s", d.Id(), err)
	}

	// After the server is unlocked && running. Check if primary_networkadapter_network is set and update the adapter.
	_, nic1netOK := d.GetOk("primary_networkadapter_network")
	if nic1netOK {
		err := setServerNetworkAdapter(ctx, d, client)
		if err != nil {
			return diag.Errorf("error while setting default adapter network: %s", err)
		}
	}

	return resourceGlesysServerRead(ctx, d, m)
}

func getTemplate(original string, srv *glesys.ServerDetails) string {
	for _, tag := range srv.InitialTemplate.CurrentTags {
		if tag == original {
			return original
		}
	}
	if original == srv.InitialTemplate.ID {
		return original
	}
	return srv.Template
}

func resourceGlesysServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// fetch updates about the resource
	srv, err := client.Servers.Details(ctx, d.Id())
	if err != nil {
		diag.Errorf("server not found: %s", err)
		d.SetId("")
		return nil
	}

	// Workaround for the API not returning the correct Bandwidth value for KVM servers
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
	d.Set("islocked", srv.IsLocked)
	d.Set("isrunning", srv.IsRunning)
	d.Set("storage", srv.Storage)
	d.Set("template", getTemplate(d.Get("template").(string), srv))
	var diskIDs []string
	for _, d := range srv.AdditionalDisks {
		diskIDs = append(diskIDs, d.ID)
	}
	d.Set("extra_disks", diskIDs)

	var backupSchedules []map[string]interface{}
	for _, bs := range srv.Backup.Schedules {
		schedule := map[string]interface{}{
			"frequency": bs.Frequency,
			"retention": bs.Numberofimagestokeep,
		}
		backupSchedules = append(backupSchedules, schedule)
	}

	if err := d.Set("backups_schedule", backupSchedules); err != nil {
		return diag.Errorf("unable to set backups_schedule, read value %v", err)
	}

	var adapters []map[string]interface{}
	netAdapters, _ := client.Servers.NetworkAdapters(ctx, d.Id())
	for _, v := range *netAdapters {
		n := map[string]interface{}{
			"id":          v.ID,
			"adaptertype": v.AdapterType,
			"bandwidth":   v.Bandwidth,
			"name":        v.Name,
			"networkid":   v.NetworkID,
		}
		if v.Name == "Network adapter 1" || v.IsPrimary {
			d.Set("primary_networkadapter_network", v.NetworkID)
		}
		adapters = append(adapters, n)
	}

	if err := d.Set("network_adapters", adapters); err != nil {
		return diag.Errorf("unable to set network_adapters, read value %v", err)
	}

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
	if d.HasChange("bandwidth") {
		opts.Bandwidth = d.Get("bandwidth").(int)
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
	if d.HasChange("backups_schedule") {
		backupsList, err := expandBackupSchedules(d.Get("backups_schedule").(*schema.Set).List())
		if err != nil {
			diag.Errorf("Error updating backups_schedule: %s", err)
		}
		opts.Backup = backupsList
	}
	_, err := client.Servers.Edit(ctx, d.Id(), opts)
	if err != nil {
		return diag.Errorf("Error updating instance: %s", err)
	}

	// Check if the setting for the primary network adapter has changed
	if d.HasChange("primary_networkadapter_network") {
		if err := setServerNetworkAdapter(ctx, d, client); err != nil {
			return diag.Errorf("error while updating primary networkadapter network: %s", err)
		}
	}

	return resourceGlesysServerRead(ctx, d, m)
}

func setServerNetworkAdapter(ctx context.Context, d *schema.ResourceData, client *glesys.Client) error {
	netadapterparams := glesys.EditNetworkAdapterParams{}
	var netadapterID string
	// fetch current networkadapters
	netAdapters, _ := client.Servers.NetworkAdapters(ctx, d.Id())
	for _, v := range *netAdapters {
		if v.Name == "Network adapter 1" || v.IsPrimary {
			netadapterID = v.ID
		}
	}

	log.Printf("[INFO]: setServerNetworkAdapter (%s) networkadapter found %s", d.Id(), netadapterID)
	netadapterparams.NetworkID = d.Get("primary_networkadapter_network").(string)

	_, err := client.NetworkAdapters.Edit(ctx, netadapterID, netadapterparams)
	if err != nil {
		return err
	}

	return nil
}

func resourceGlesysServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// Call waitForServerAttribute to make sure the server isn't locked before deleting it.
	_, err := waitForServerAttribute(ctx, d, "false", []string{"true"}, "islocked", m)

	if err != nil {
		return diag.Errorf("Error waiting for server to be unlocked for destroy (%s): %s", d.Id(), err)
	}
	// destroy the server, don't keep the ip.
	err = client.Servers.Destroy(ctx, d.Id(), glesys.DestroyServerParams{KeepIP: false})

	if err != nil {
		return diag.Errorf("Error deleting instance (%s): %s", d.Id(), err)
	}

	return resourceGlesysServerRead(ctx, d, m)
}

// waitForServerAttribute
func waitForServerAttribute(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, m interface{}) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    serverStateRefresh(ctx, d, m, attribute),
		Timeout:    20 * time.Minute,
		Delay:      6 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func serverStateRefresh(ctx context.Context, d *schema.ResourceData, m interface{}, attr string) retry.StateRefreshFunc {
	client := m.(*glesys.Client)
	return func() (interface{}, string, error) {
		// check state of server
		server, err := client.Servers.Details(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Server (%s): %s", d.Id(), err)
		}

		// depending on attribute, check if locked or running
		switch attr {
		case "islocked":
			log.Printf("[INFO] Still locked %s", d.Id())
			return server, strconv.FormatBool(server.IsLocked), nil
		case "isrunning":
			running := strconv.FormatBool(server.IsRunning)
			log.Printf("[INFO] Server (%s) started: %s", d.Id(), running)
			return server, running, nil
		default:
			return nil, "", nil
		}
	}
}
