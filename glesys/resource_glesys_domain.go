package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysDNSDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysDNSDomainCreate,
		ReadContext:   resourceGlesysDNSDomainRead,
		UpdateContext: resourceGlesysDNSDomainUpdate,
		DeleteContext: resourceGlesysDNSDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Domain name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"createrecords": {
				Description: "Create default set of records when creating the domain. `0/1, yes/no, true/false`",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"createtime": {
				Description: "Domain create time",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"displayname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expire": {
				Description: "Domain expire TTL",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"minimum": {
				Description: "Domain minimum TTL",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"refresh": {
				Description: "Domain refresh TTL",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"retry": {
				Description: "Domain retry TTL",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"ttl": {
				Description: "Domain default TTL",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"primarynameserver": {
				Description: "Domain primary nameserver",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},

			"recordcount": {
				Description: "Number of records for the domain",
				Type:        schema.TypeInt,
				Computed:    true,
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

func resourceGlesysDNSDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	domain, err := client.DNSDomains.AddDNSDomain(ctx, params)
	if err != nil {
		return diag.Errorf("Error adding domain %s: %v", params.Name, err)
	}

	// Set the Id to domain.ID
	d.SetId(domain.Name)

	return resourceGlesysDNSDomainRead(ctx, d, m)
}

func resourceGlesysDNSDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	domain, err := client.DNSDomains.Details(ctx, d.Id())

	if err != nil {
		diag.Errorf("domain not found: %v", err)
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

func resourceGlesysDNSDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	_, err := client.DNSDomains.Edit(ctx, params)
	if err != nil {
		return diag.Errorf("Error updating domain: %v", err)
	}

	return resourceGlesysDNSDomainRead(ctx, d, m)
}

func resourceGlesysDNSDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.DeleteDNSDomainParams{
		Name: d.Id(),
	}

	err := client.DNSDomains.Delete(ctx, params)
	if err != nil {
		return diag.Errorf("Error deleting domain: %v", err)
	}
	d.SetId("")
	return nil
}
