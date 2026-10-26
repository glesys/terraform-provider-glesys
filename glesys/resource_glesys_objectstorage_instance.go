package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
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

		Description: "ObjectStorage Instance Resource.",
		Schema: map[string]*schema.Schema{
			"datacenter": {
				Description: "ObjectStorage location.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "ObjectStorage Description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created": {
				Description: "ObjectStorage created timestamp.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"accesskey": {
				Description: "ObjectStorage default access key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"secretkey": {
				Description: "ObjectStorage default secret key.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
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

	if cred, ok := trackedObjectStorageCredential(instance.Credentials, d.Get("accesskey").(string)); ok {
		d.Set("accesskey", cred.AccessKey)
		_, secretkeyIsSet := d.GetOk("secretkey")
		if !secretkeyIsSet {
			d.Set("secretkey", cred.SecretKey)
		}
	}

	return nil
}

// trackedObjectStorageCredential returns the credential that the accesskey and
// secretkey attributes track. In observed responses the details call lists
// newer credentials first and omits secret keys, so the secret key in state is
// the one stored when the instance was created. A credential added later, for
// example by the object browser in Glesys Cloud, must not replace the access
// key in state: that would pair its access key with the original secret key.
// The credential whose access key is already in state is kept. Without one (a
// new or imported instance, or the tracked credential was deleted) the first
// credential is used, as before; after a deletion the stored secret key no
// longer matches, which this function cannot repair because no secret key is
// returned.
func trackedObjectStorageCredential(creds []glesys.ObjectStorageCredential, accessKey string) (glesys.ObjectStorageCredential, bool) {
	if len(creds) == 0 {
		return glesys.ObjectStorageCredential{}, false
	}
	if accessKey != "" {
		for _, c := range creds {
			if c.AccessKey == accessKey {
				return c, true
			}
		}
	}
	return creds[0], true
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
