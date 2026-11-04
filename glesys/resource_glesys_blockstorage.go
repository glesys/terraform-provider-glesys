package glesys

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysBlockStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysBlockStorageCreate,
		ReadContext:   resourceGlesysBlockStorageRead,
		UpdateContext: resourceGlesysBlockStorageUpdate,
		DeleteContext: resourceGlesysBlockStorageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "BlockStorage volume.",

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Volume ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"datacenter": {
				Description: "Volume location.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Volume descriptive name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "Volume size in GiB.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"serverid": {
				Description: "Attached `glesys_server` id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"status": {
				Description: "Volume status.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"wwn": {
				Description: "Volume WWN id.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"tier": {
				Description: "Volume storage class tier [capacity|balanced|velocity].",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGlesysBlockStorageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Setup client to the API
	client := m.(*glesys.Client)

	// Setup blockstorage volume parameters
	params := glesys.CreateBlockStorageParams{
		Name:       d.Get("name").(string),
		Datacenter: d.Get("datacenter").(string),
		SizeInGIB:  d.Get("size").(int),
		Tier:       d.Get("tier").(string),
	}

	volume, err := client.BlockStorageVolumes.Create(ctx, params)
	if err != nil {
		return diag.Errorf("error creating volume: %+v", err)
	}

	// Set the resource Id to server ID
	d.SetId(volume.VolumeID)

	// Wait for volume to be ready
	if serverID, ok := d.GetOk("serverid"); ok {
		if _, err := waitForVolumeLocked(ctx, volume.VolumeID, "READY", []string{"true"}, "status", m); err != nil {
			return diag.Errorf("disk: error while waiting for Volume (%s) to be ready: %s", volume.VolumeID, err)
		}
		_, errAttach := client.BlockStorageVolumes.Attach(ctx, volume.VolumeID, serverID.(string))
		if errAttach != nil {
			return diag.Errorf("volume: error while attaching volume (%s) to server (%s): %v", volume.VolumeID, serverID.(string), errAttach)
		}
	}

	return resourceGlesysBlockStorageRead(ctx, d, m)
}

func waitForVolumeLocked(ctx context.Context, volumeID string, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:        pending,
		Target:         []string{target},
		Refresh:        blockstorageVolumeStateRefresh(ctx, volumeID, meta),
		Timeout:        10 * time.Minute,
		Delay:          5 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 60,
	}

	return stateConf.WaitForStateContext(ctx)
}

func blockstorageVolumeStateRefresh(ctx context.Context, volumeID string, meta interface{}) retry.StateRefreshFunc {
	client := meta.(*glesys.Client)

	return func() (interface{}, string, error) {
		volume, err := client.BlockStorageVolumes.Details(ctx, volumeID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving volume details for (%s) : %s", volumeID, err)
		}
		return &volume, volume.Status, nil
	}
}

func resourceGlesysBlockStorageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	vol, err := client.BlockStorageVolumes.Details(ctx, d.Id())
	if err != nil {
		diag.Errorf("blockstorage volume not found: %s", err)
		d.SetId("")
		return nil
	}
	log.Printf("[INFO]: Found volume: (%s - %s)", vol.VolumeID, vol.Name)

	if len(vol.Attachments) == 0 {
		log.Printf("[DEBUG]: volume not attached (%v)", d.Id())
		d.Set("serverid", "")
	} else {
		for _, att := range vol.Attachments {
			log.Printf("[DEBUG]: attachments: %v", att)
			d.Set("serverid", att.ServerID)
		}
	}

	d.Set("name", vol.Name)
	d.Set("datacenter", vol.Datacenterkey)
	d.Set("size", vol.SizeInGIB)
	d.Set("wwn", vol.WWN)
	d.Set("tier", vol.Tier.Slug)
	d.Set("status", vol.Status)

	return nil
}

func resourceGlesysBlockStorageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditBlockStorageParams{
		VolumeID: d.Id(),
	}

	// UpdateName has it's own function.
	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
		_, err := client.BlockStorageVolumes.UpdateName(ctx, params)
		if err != nil {
			return diag.Errorf("error updating volume Name: %s", err)
		}
	}

	if d.HasChange("size") {
		params.SizeInGIB = d.Get("size").(int)
		_, err := client.BlockStorageVolumes.Resize(ctx, params)
		if err != nil {
			return diag.Errorf("Error updating volume size: %s", err)
		}
	}

	if d.HasChange("serverid") {
		old, newSrv := d.GetChange("serverid")

		if old.(string) != "" {
			// If the volume is already attached to a server instance. Detach and attach.

			bsv, err := client.BlockStorageVolumes.Details(ctx, d.Id())
			if err != nil {
				return diag.Errorf("Error getting volume: %s", err)
			}

			if len(bsv.Attachments) > 0 {
				detachErr := client.BlockStorageVolumes.Detach(ctx, d.Id())
				if detachErr != nil {
					return diag.Errorf("Error when detaching volume (%s): %s", d.Id(), detachErr)
				}
			}

			_, attachErr := client.BlockStorageVolumes.Attach(ctx, d.Id(), newSrv.(string))
			if attachErr != nil {
				return diag.Errorf("Error when attaching volume (%s): %s", d.Id(), attachErr)
			}
		} else {
			_, attachErr := client.BlockStorageVolumes.Attach(ctx, d.Id(), newSrv.(string))
			if attachErr != nil {
				return diag.Errorf("Error when attaching volume (%s): %s", d.Id(), attachErr)
			}
		}
	}

	return resourceGlesysBlockStorageRead(ctx, d, m)
}

func resourceGlesysBlockStorageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	volID := d.Id()

	if d.Get("serverid") != "" {
		// if volume is attached to a server, detach first
		detachErr := client.BlockStorageVolumes.Detach(ctx, volID)

		if _, err := waitForVolumeLocked(ctx, volID, "READY", []string{"true"}, "status", m); err != nil {
			return diag.Errorf("disk: error while waiting for Volume (%s) to be ready: %s", volID, err)
		}
		if detachErr != nil {
			if !(strings.Contains(detachErr.Error(), "The volume is not attached to a server")) {
				return diag.Errorf("delete: Error when detaching volume (%s): %s", d.Id(), detachErr)
			}
		}
	}

	if _, err := waitForVolumeLocked(ctx, volID, "READY", []string{"true"}, "status", m); err != nil {
		return diag.Errorf("disk: error while waiting for Volume (%s) to be ready: %s", volID, err)
	}
	err := client.BlockStorageVolumes.Delete(ctx, volID)
	if err != nil {
		return diag.Errorf("Error deleting blockstorage volume (%s): %s", volID, err)
	}

	d.SetId("")
	return nil
}
