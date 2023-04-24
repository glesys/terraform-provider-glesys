package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysObjectStorageInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysObjectStorageInstanceCreate,
		ReadContext:   resourceGlesysObjectStorageInstanceRead,
		UpdateContext: resourceGlesysObjectStorageInstanceUpdate,
		DeleteContext: resourceGlesysObjectStorageInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accesskey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secretkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGlesysObjectStorageInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreateObjectStorageInstanceParams{
		DataCenter:  d.Get("datacenter").(string),
		Description: d.Get("description").(string),
	}

	instance, err := client.ObjectStorages.CreateInstance(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating object storage: %s", err)
	}

	d.SetId(instance.InstanceID)
	d.Set("secretkey", instance.Credentials[0].SecretKey)

	return resourceGlesysObjectStorageInstanceRead(ctx, d, m)
}

func resourceGlesysObjectStorageInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	instance, err := client.ObjectStorages.InstanceDetails(ctx, d.Id())
	if err != nil {
		d.SetId("")
		return diag.Errorf("object storage not found: %s", err)
	}

	d.Set("datacenter", instance.DataCenter)
	d.Set("description", instance.Description)
	d.Set("created", instance.Created)

	creds := instance.Credentials
	if len(creds) > 0 {
		d.Set("accesskey", creds[0].AccessKey)
		_, secretkeyIsSet := d.GetOk("secretkey")
		if !secretkeyIsSet {
			d.Set("secretkey", creds[0].SecretKey)
		}
	}

	return nil
}

func resourceGlesysObjectStorageInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditObjectStorageInstanceParams{
		InstanceID: d.Id(),
	}

	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
	}

	_, err := client.ObjectStorages.EditInstance(ctx, params)
	if err != nil {
		return diag.Errorf("Error updating object storage: %s", err)
	}
	return resourceGlesysObjectStorageInstanceRead(ctx, d, m)
}

func resourceGlesysObjectStorageInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.ObjectStorages.DeleteInstance(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting object storage: %s", err)
	}

	return nil
}
