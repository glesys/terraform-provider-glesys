package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysPrivateNetworkSegment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysPrivateNetworkSegmentCreate,
		ReadContext:   resourceGlesysPrivateNetworkSegmentRead,
		UpdateContext: resourceGlesysPrivateNetworkSegmentUpdate,
		DeleteContext: resourceGlesysPrivateNetworkSegmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Create a PrivateNetwork Segment to connect VM NetworkAdapters.",

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "PrivateNetworkSegment name",
				Type:        schema.TypeString,
				Required:    true,
			},

			"privatenetworkid": {
				Description: "PrivateNetwork ID for the Segment.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"platform": {
				Description: "PrivateNetworkSegment Platform.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"datacenter": {
				Description: "PrivateNetworkSegment Datacenter.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"ipv4subnet": {
				Description: "PrivateNetworkSegment IPv4 Subnet.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"ipv6subnet": {
				Description: "PrivateNetworkSegment IPv6 Subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGlesysPrivateNetworkSegmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreatePrivateNetworkSegmentParams{
		PrivateNetworkID: d.Get("privatenetworkid").(string),
		Name:             d.Get("name").(string),
		Datacenter:       d.Get("datacenter").(string),
		Platform:         d.Get("platform").(string),
		IPv4Subnet:       d.Get("ipv4subnet").(string),
	}

	name := d.Get("name").(string)
	segment, err := client.PrivateNetworks.CreateSegment(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating segment %s: %v", name, err)
	}

	// Set the Id to network.ID
	d.SetId(segment.ID)

	return resourceGlesysPrivateNetworkSegmentRead(ctx, d, m)
}

func resourceGlesysPrivateNetworkSegmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	// List segments for 'privatenetworkid'
	segments, err := client.PrivateNetworks.ListSegments(ctx, d.Get("privatenetworkid").(string))

	if err != nil {
		diag.Errorf("privatenetwork not found: %v", err)
		d.SetId("")
		return nil
	}

	for _, seg := range *segments {
		if seg.ID == d.Id() {
			d.Set("datacenter", seg.Datacenter)
			d.Set("ipv4subnet", seg.IPv4Subnet)
			d.Set("ipv6subnet", seg.IPv6Subnet)
			d.Set("name", seg.Name)
			d.Set("platform", seg.Platform)
		}
	}

	return nil
}

func resourceGlesysPrivateNetworkSegmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditPrivateNetworkSegmentParams{ID: d.Id()}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	_, err := client.PrivateNetworks.EditSegment(ctx, params)
	if err != nil {
		return diag.Errorf("Error updating segment: %v", err)
	}

	return resourceGlesysPrivateNetworkSegmentRead(ctx, d, m)
}

func resourceGlesysPrivateNetworkSegmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.PrivateNetworks.DestroySegment(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting segment: %v", err)
	}
	d.SetId("")
	return nil
}
