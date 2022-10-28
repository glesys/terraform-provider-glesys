package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysIPCreate,
		ReadContext:   resourceGlesysIPRead,
		UpdateContext: resourceGlesysIPUpdate,
		DeleteContext: resourceGlesysIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"broadcast": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amount": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"currency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"locked_to_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_servers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platforms": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ptr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGlesysIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Setup client to the API
	client := m.(*glesys.Client)

	address := d.Get("address").(string)
	if address == "" {
		// A reserved address has not been set
		params := glesys.AvailableIPsParams{
			DataCenter: d.Get("datacenter").(string),
			Platform:   d.Get("platform").(string),
			Version:    d.Get("version").(int),
		}

		// Get available ip addresses
		ips, err := client.IPs.Available(context.Background(), params)
		if err != nil {
			return diag.Errorf("Error listing available IPs: %v", err)
		}

		// Select the first available ip address for reservation
		address = (*ips)[0].Address
	}

	ip, err := client.IPs.Reserve(context.Background(), address)
	if err != nil {
		return diag.Errorf("Error reserving IP %s: %v", address, err)
	}

	ptr := d.Get("ptr").(string)
	if ptr != "" {
		_, err := client.IPs.SetPTR(context.Background(), address, ptr)

		if err != nil {
			return diag.Errorf("Error setting PTR %s: %v", ptr, err)
		}
	}

	// Set the resource Id to IP address
	d.SetId(ip.Address)
	return resourceGlesysIPRead(ctx, d, m)
}

func resourceGlesysIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// Fetch updates about the IP
	ip, err := client.IPs.Details(context.Background(), d.Id())
	if err != nil {
		diag.Errorf("IP not found: %s", err)
		d.SetId("")
		return nil
	}

	d.Set("address", ip.Address)
	d.Set("broadcast", ip.Broadcast)
	d.Set("datacenter", ip.DataCenter)
	d.Set("gateway", ip.Gateway)
	d.Set("locked_to_account", ip.LockedToAccount)
	d.Set("name_servers", ip.NameServers)
	d.Set("netmask", ip.Netmask)
	d.Set("platform", ip.Platform)
	d.Set("platforms", ip.Platforms)
	d.Set("reserved", ip.Reserved)
	d.Set("server_id", ip.ServerID)
	d.Set("ptr", ip.PTR)

	cost := map[string]interface{}{
		"amount":      ip.Cost.Amount,
		"currency":    ip.Cost.Currency,
		"time_period": ip.Cost.TimePeriod,
	}
	err = d.Set("cost", []map[string]interface{}{cost})
	if err != nil {
		return diag.Errorf("Error setting cost: %v", err)
	}

	return nil
}

func resourceGlesysIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	if d.HasChange("ptr") {
		// There should be support here for resetting pointer records when they are zeroed.
		// Because of how Optional+Computed attributes work it is not possible with the current SDK.
		// More info in upstream issue #282: https://github.com/hashicorp/terraform-plugin-sdk/issues/282

		ptr := d.Get("ptr").(string)
		_, err := client.IPs.SetPTR(context.Background(), d.Id(), ptr)
		if err != nil {
			return diag.Errorf("Error updating reverse pointer on IP: %s", err)
		}
	}

	return nil
}

func resourceGlesysIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.IPs.Release(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error releasing IP %s: %v", d.Id(), err)
	}

	d.SetId("")
	return nil
}
