package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysObjectStorageCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysObjectStorageCredentialCreate,
		ReadContext:   resourceGlesysObjectStorageCredentialRead,
		DeleteContext: resourceGlesysObjectStorageCredentialDelete,

		Schema: map[string]*schema.Schema{
			"instanceid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accesskey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secretkey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

	credential, err := client.ObjectStorages.CreateCredential(context.Background(), params)
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

	err := client.ObjectStorages.DeleteCredential(context.Background(), params)
	if err != nil {
		return diag.Errorf("Error deleting object storage credential: %s", err)
	}

	return nil
}
