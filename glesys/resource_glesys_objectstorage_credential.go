package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysObjectStorageCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysObjectStorageCredentialCreate,
		ReadContext:   resourceGlesysObjectStorageCredentialRead,
		DeleteContext: resourceGlesysObjectStorageCredentialDelete,

		Description: "ObjectStorage Credentials.",
		Schema: map[string]*schema.Schema{
			"instanceid": {
				Description: "Associated ObjectStorage instance.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"accesskey": {
				Description: "ObjectStorage credential access key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"secretkey": {
				Description: "ObjectStorage credential secret key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created": {
				Description: "ObjectStorage credential created timestamp.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "ObjectStorage credential description.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGlesysObjectStorageCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreateObjectStorageCredentialParams{
		InstanceID:  d.Get("instanceid").(string),
		Description: d.Get("description").(string),
	}

	credential, err := client.ObjectStorages.CreateCredential(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating object storage credential: %s", err)
	}

	d.SetId(credential.CredentialID)
	d.Set("accesskey", credential.AccessKey)
	d.Set("description", credential.Description)
	d.Set("secretkey", credential.SecretKey)
	d.Set("created", credential.Created)

	return nil
}

func resourceGlesysObjectStorageCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceGlesysObjectStorageCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.DeleteObjectStorageCredentialParams{
		InstanceID:   d.Get("instanceid").(string),
		CredentialID: d.Id(),
	}

	err := client.ObjectStorages.DeleteCredential(ctx, params)
	if err != nil {
		return diag.Errorf("Error deleting object storage credential: %s", err)
	}

	return nil
}
