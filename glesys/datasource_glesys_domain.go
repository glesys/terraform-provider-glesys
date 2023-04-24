package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGlesysDNSDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about a DNS Domain associated with your GleSYS Project.",

		ReadContext: dataSourceGlesysDomainRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the domain",
				ValidateFunc: validation.NoZeroValues,
			},

			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ttl of the domain.",
			},

			"expire": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "expire ttl of the domain.",
			},

			"retry": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "retry ttl of the domain.",
			},

			"refresh": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "refresh ttl of the domain.",
			},

			"minimum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "minimum ttl of the domain.",
			},
		},
	}
}

func dataSourceGlesysDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	name := d.Get("name").(string)

	domain, err := client.DNSDomains.Details(ctx, name)

	if err != nil {
		diag.Errorf("domain not found: %v", err)
		d.SetId("")
		return nil
	}

	d.SetId(domain.Name)
	d.Set("name", domain.Name)
	d.Set("ttl", domain.TTL)
	d.Set("expire", domain.Expire)
	d.Set("minimum", domain.Minimum)
	d.Set("refresh", domain.Refresh)
	d.Set("retry", domain.Retry)

	return nil
}
