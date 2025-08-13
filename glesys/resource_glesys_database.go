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

func resourceGlesysDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysDatabaseCreate,
		UpdateContext: resourceGlesysDatabaseUpdate,
		ReadContext:   resourceGlesysDatabaseRead,
		DeleteContext: resourceGlesysDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Database ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Database name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"engine": {
				Description: "Database engine name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"engineversion": {
				Description: "Database engine version",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"datacenterkey": {
				Description: "Datacenter key",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"fqdn": {
				Description: "Database FQDN",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "Database status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"plankey": {
				Description: "Database plan key",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"connectionstring": {
				Description: "Connectionstring to access database",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"allowlist": {
				Description: "Update the allow list for a database instance list to be either a single IP address or a CIDR range.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGlesysDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// Create a database in the Glesys platform
	rawAllowlist := d.Get("allowlist").([]interface{})
	allowlist, err := convertResourceDataToListOfStrings(rawAllowlist)
	params := glesys.CreateDatabaseParams{
		Name:          d.Get("name").(string),
		DataCenterKey: d.Get("datacenterkey").(string),
		Engine:        d.Get("engine").(string),
		EngineVersion: d.Get("engineversion").(string),
		PlanKey:       d.Get("plankey").(string),
		AllowList:     allowlist,
	}

	database, err := client.Databases.Create(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating database %s: %v", params.Name, err)
	}

	// Set the Id to domain.ID
	d.SetId(database.ID)
	if _, err = waitForDatabaseAttribute(ctx, d, "true", []string{"false"}, "RUNNING", m); err != nil {
		return diag.Errorf("error while waiting for Server (%s) to be started: %s", d.Id(), err)
	}

	return resourceGlesysDatabaseRead(ctx, d, m)
}

func convertResourceDataToListOfStrings(raw []interface{}) ([]string, error) {
	strs := make([]string, len(raw))
	for i, v := range raw {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}
		strs[i] = str
	}
	return strs, nil
}

func updateAllowlist(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	rawAllowlist := d.Get("allowlist").([]interface{})
	allowlist, err := convertResourceDataToListOfStrings(rawAllowlist)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to convert allowlist: %w", err))
	}

	_, err = client.Databases.UpdateAllowlist(ctx, glesys.UpdateAllowlistParams{
		ID:        d.Get("id").(string),
		AllowList: allowlist,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update allowlist: %w", err))
	}

	if _, err := waitForDatabaseAttribute(ctx, d, "true", []string{"false"}, "RUNNING", m); err != nil {
		return diag.Errorf("error while waiting for Server (%s) to be started: %s", d.Id(), err)
	}
	return nil
}

func resourceGlesysDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	database, err := client.Databases.Details(ctx, d.Id())
	if err != nil {
		diag.Errorf("database not found: %v", err)
		d.SetId("")
		return nil
	}
	connectionstring, err := client.Databases.ConnectionString(ctx, d.Id())
	fmt.Println(connectionstring)
	if err != nil {
		diag.Errorf("database not found: %v", err)
		d.SetId("")
		return nil
	}

	d.Set("id", database.ID)
	d.Set("name", database.Name)
	d.Set("engine", strings.ToLower(database.Engine))
	d.Set("engineversion", database.EngineVersion)
	d.Set("datacenterkey", database.DataCenterKey)
	d.Set("fqdn", database.Fqdn)
	d.Set("status", database.Status)
	d.Set("allowlist", database.Allowlist)
	d.Set("connectionstring", connectionstring.ConnectionString)
	d.Set("plan_key", database.Plan.Key)
	d.Set("plan_cpucores", database.Plan.CpuCores)
	d.Set("plan_memoryingb", database.Plan.MemoryInGib)
	d.Set("plan_storageingb", database.Plan.StorageInGib)
	d.Set("maintenancewindow_durationinminutes", database.MaintenanceWindow.DurationInMinutes)
	d.Set("maintenancewindow_starttime", database.MaintenanceWindow.StartTime)
	d.Set("maintenancewindow_weekday", database.MaintenanceWindow.WeekDay)

	return nil
}

func resourceGlesysDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	updateAllowlist(ctx, d, m)

	return nil
}

func resourceGlesysDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.Databases.Delete(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting database: %v", err)
	}
	d.SetId("")
	return nil
}

// waitForServerAttribute
func waitForDatabaseAttribute(
	ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, m interface{}) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    databaseStateRefresh(ctx, d, m, attribute),
		Timeout:    20 * time.Minute,
		Delay:      6 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func databaseStateRefresh(ctx context.Context, d *schema.ResourceData, m interface{}, attr string) retry.StateRefreshFunc {
	client := m.(*glesys.Client)
	return func() (interface{}, string, error) {
		// check state of server
		database, err := client.Databases.Details(ctx, d.Id())
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving Database (%s): %s", d.Id(), err)
		}

		return database, strconv.FormatBool(database.Status == attr), nil

	}
}
