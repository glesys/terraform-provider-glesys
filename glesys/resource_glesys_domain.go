package glesys

import (
	"context"
	"fmt"

	"github.com/glesys/glesys-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlesysDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlesysDNSDomainCreate,
		Read:   resourceGlesysDNSDomainRead,
		Update: resourceGlesysDNSDomainUpdate,
		Delete: resourceGlesysDNSDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"createrecords": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"createtime": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"displayname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expire": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"minimum": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"refresh": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"retry": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"primarynameserver": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"recordcount": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"registrarinfo_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registrarinfo_statedescr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registrarinfo_expire": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registrarinfo_autorenew": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registrarinfo_tld": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registrarinfo_invoicenumber": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"responsibleperson": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"usingglesysnameserver": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGlesysDNSDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	// Add a domain in the glesys platform. Do not register new domains right now.
	params := glesys.AddDNSDomainParams{
		Name:              d.Get("name").(string),
		CreateRecords:     d.Get("createrecords").(string),
		Expire:            d.Get("expire").(int),
		Minimum:           d.Get("minimum").(int),
		Refresh:           d.Get("refresh").(int),
		Retry:             d.Get("retry").(int),
		TTL:               d.Get("ttl").(int),
		PrimaryNameServer: d.Get("primarynameserver").(string),
		ResponsiblePerson: d.Get("responsibleperson").(string),
	}

	domain, err := client.DNSDomains.AddDNSDomain(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error adding domain %s: %v", params.Name, err)
	}

	// Set the Id to domain.ID
	d.SetId((*domain).Name)

	return resourceGlesysDNSDomainRead(d, m)
}

func resourceGlesysDNSDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	domain, err := client.DNSDomains.Details(context.Background(), d.Id())

	if err != nil {
		fmt.Errorf("domain not found: %v", err)
		d.SetId("")
		return nil
	}

	d.Set("createtime", domain.CreateTime)
	d.Set("displayname", domain.DisplayName)
	d.Set("expire", domain.Expire)
	d.Set("minimum", domain.Minimum)
	d.Set("name", domain.Name)
	d.Set("refresh", domain.Refresh)
	d.Set("retry", domain.Retry)
	d.Set("ttl", domain.TTL)
	d.Set("recordcount", domain.RecordCount)
	d.Set("responsibleperson", domain.ResponsiblePerson)
	d.Set("primarynameserver", domain.PrimaryNameServer)
	d.Set("usingglesysnameserver", domain.UsingGlesysNameserver)

	d.Set("registrarinfo_state", domain.RegistrarInfo.State)
	d.Set("registrarinfo_statedescr", domain.RegistrarInfo.StateDescription)
	d.Set("registrarinfo_expire", domain.RegistrarInfo.Expire)
	d.Set("registrarinfo_autorenew", domain.RegistrarInfo.AutoRenew)
	d.Set("registrarinfo_tld", domain.RegistrarInfo.TLD)
	d.Set("registrarinfo_invoicenumber", domain.RegistrarInfo.InvoiceNumber)

	return nil
}

func resourceGlesysDNSDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.EditDNSDomainParams{Name: d.Id()}

	if d.HasChange("expire") {
		params.Expire = d.Get("expire").(int)
	}

	if d.HasChange("minimum") {
		params.Minimum = d.Get("minimum").(int)
	}

	if d.HasChange("refresh") {
		params.Refresh = d.Get("refresh").(int)
	}

	if d.HasChange("retry") {
		params.Retry = d.Get("retry").(int)
	}

	if d.HasChange("ttl") {
		params.TTL = d.Get("ttl").(int)
	}

	if d.HasChange("primarynameserver") {
		params.PrimaryNameServer = d.Get("primarynameserver").(string)
	}

	if d.HasChange("responsibleperson") {
		params.ResponsiblePerson = d.Get("responsibleperson").(string)
	}

	_, err := client.DNSDomains.Edit(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error updating domain: %v", err)
	}

	return resourceGlesysDNSDomainRead(d, m)
}

func resourceGlesysDNSDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*glesys.Client)

	params := glesys.DeleteDNSDomainParams{
		Name: d.Id(),
	}

	err := client.DNSDomains.Delete(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error deleting domain: %v", err)
	}
	d.SetId("")
	return nil
}
