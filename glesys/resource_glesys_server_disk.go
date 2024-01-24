package glesys

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysServerDisk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysServerDiskCreate,
		ReadContext:   resourceGlesysServerDiskRead,
		UpdateContext: resourceGlesysServerDiskUpdate,
		DeleteContext: resourceGlesysServerDiskDelete,

		Description: "An additional disk associated with a `glesys_server`",

		Importer: &schema.ResourceImporter{
			StateContext: resourceGleSYSServerDiskImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Disk ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Disk descriptive name.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"size": {
				Description: "Disk size in GIB",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"serverid": {
				Description: "Associated `glesys_server` id.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"scsiid": {
				Description: "Disk unit number.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

// resourceGleSYSServerDiskImport - import additional disks "wps12345,000000-1111-222-3333333"
func resourceGleSYSServerDiskImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		if len(s) < 2 {
			return nil, fmt.Errorf("not enough parameters ( <serverid>,<diskid> ) : %s", s)
		}

		d.SetId(s[1])
		d.Set("serverid", s[0])
	}

	return []*schema.ResourceData{d}, nil
}
func resourceGlesysServerDiskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Setup client to the API
	client := m.(*glesys.Client)

	// Setup server parameters
	params := glesys.CreateServerDiskParams{
		Name:      d.Get("name").(string),
		SizeInGIB: d.Get("size").(int),
		ServerID:  d.Get("serverid").(string),
	}

	// Wait for server to be running && !islocked
	if _, err := waitForServerLocked(ctx, params.ServerID, "false", []string{"true"}, "islocked", m); err != nil {
		return diag.Errorf("disk: error while waiting for Server (%s) to be completed: %s", params.ServerID, err)
	}
	disk, err := client.ServerDisks.Create(ctx, params)

	if err != nil {
		return diag.Errorf("error creating disk: %+v", err)
	}

	// Set the resource Id to server ID
	d.SetId(disk.ID)

	return resourceGlesysServerDiskRead(ctx, d, m)
}

func waitForServerLocked(ctx context.Context, serverID string, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        serverdiskStateRefresh(ctx, serverID, meta),
		Timeout:        10 * time.Minute,
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}
func serverdiskStateRefresh(ctx context.Context, serverID string, meta interface{}) retry.StateRefreshFunc {
	client := meta.(*glesys.Client)

	return func() (interface{}, string, error) {
		server, err := client.Servers.Details(ctx, serverID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving server details for %s : %s", serverID, err)
		}
		return &server, strconv.FormatBool(server.IsLocked), nil
	}
}

func resourceGlesysServerDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	serverid := d.Get("serverid").(string)
	server, err := client.Servers.Details(ctx, serverid)
	if err != nil {
		diag.Errorf("server not found: %s", err)
		d.SetId("")
		return nil
	}

	for _, n := range server.AdditionalDisks {
		if n.ID == d.Get("id").(string) {
			d.Set("name", n.Name)
			d.Set("size", n.SizeInGIB)
			d.Set("scsiid", n.SCSIID)
		}
	}

	return nil
}

func resourceGlesysServerDiskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditServerDiskParams{
		ID: d.Get("id").(string),
	}

	// UpdateName has it's own function.
	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
		_, err := client.ServerDisks.UpdateName(ctx, params)
		if err != nil {
			return diag.Errorf("Error updating ServerDisk Name: %s", err)
		}
	}

	if d.HasChange("size") {
		params.SizeInGIB = d.Get("size").(int)
		_, err := client.ServerDisks.Reconfigure(ctx, params)
		if err != nil {
			return diag.Errorf("Error updating ServerDisk Name: %s", err)
		}
	}
	// If further attributes can be changed in the future, add them here.

	return resourceGlesysServerDiskRead(ctx, d, m)
}

func resourceGlesysServerDiskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	diskid := d.Get("id").(string)

	err := client.ServerDisks.Delete(ctx, diskid)
	if err != nil {
		return diag.Errorf("Error deleting ServerDisk (%s): %s", diskid, err)
	}

	d.SetId("")
	return nil
}
