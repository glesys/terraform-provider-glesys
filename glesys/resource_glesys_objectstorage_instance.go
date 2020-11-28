package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysObjectStorageInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysObjectStorageInstanceCreate,
		Read:   resourceGlesysObjectStorageInstanceRead,
		Update: resourceGlesysObjectStorageInstanceUpdate,
		Delete: resourceGlesysObjectStorageInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceGlesysObjectStorageInstanceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.CreateObjectStorageInstanceParams{
		DataCenter:  d.Get("datacenter").(string),
		Description: d.Get("description").(string),
	}

	instance, err := client.ObjectStorages.CreateInstance(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating object storage: %s", err)
	}

	d.SetId(instance.InstanceID)
	d.Set("secretkey", instance.Credentials[0].SecretKey)

	return resourceGlesysObjectStorageInstanceRead(d, m)
}

func resourceGlesysObjectStorageInstanceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	instance, err := client.ObjectStorages.InstanceDetails(context.Background(), d.Id())
	if err != nil {
		d.SetId("")
		return fmt.Errorf("object storage not found: %s", err)
	}

	d.Set("datacenter", instance.DataCenter)
	d.Set("description", instance.Description)
	d.Set("created", instance.Created)
	d.Set("accesskey", instance.Credentials[0].AccessKey)

	_, secretkeyIsSet := d.GetOk("secretkey")
	if !secretkeyIsSet {
		d.Set("secretkey", instance.Credentials[0].SecretKey)
	}

	return nil
}

func resourceGlesysObjectStorageInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.EditObjectStorageInstanceParams{
		InstanceID: d.Id(),
	}

	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
	}

	_, err := client.ObjectStorages.EditInstance(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error updating object storage: %s", err)
	}
	return resourceGlesysObjectStorageInstanceRead(d, m)
}

func resourceGlesysObjectStorageInstanceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	err := client.ObjectStorages.DeleteInstance(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting object storage: %s", err)
	}

	return nil
}
