package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysObjectStorageCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysObjectStorageCredentialCreate,
		Read:   resourceGlesysObjectStorageCredentialRead,
		Delete: resourceGlesysObjectStorageCredentialDelete,

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

func resourceGlesysObjectStorageCredentialCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.CreateObjectStorageCredentialParams{
		InstanceID:  d.Get("instanceid").(string),
		Description: d.Get("description").(string),
	}

	credential, err := client.ObjectStorages.CreateCredential(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating object storage credential: %s", err)
	}

	d.SetId(credential.CredentialID)
	d.Set("accesskey", credential.AccessKey)
	d.Set("description", credential.Description)
	d.Set("secretkey", credential.SecretKey)
	d.Set("created", credential.Created)

	return nil
}

func resourceGlesysObjectStorageCredentialRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGlesysObjectStorageCredentialDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.DeleteObjectStorageCredentialParams{
		InstanceID:   d.Get("instanceid").(string),
		CredentialID: d.Id(),
	}

	err := client.ObjectStorages.DeleteCredential(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error deleting object storage credential: %s", err)
	}

	return nil
}
