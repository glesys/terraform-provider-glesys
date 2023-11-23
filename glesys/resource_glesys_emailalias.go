package glesys

import (
	"context"
	"log"
	"strings"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysEmailAlias() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysEmailAliasCreate,
		ReadContext:   resourceGlesysEmailAliasRead,
		UpdateContext: resourceGlesysEmailAliasUpdate,
		DeleteContext: resourceGlesysEmailAliasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Create a GleSYS Email alias.",

		Schema: map[string]*schema.Schema{
			"emailalias": {
				Description: "Email alias name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"goto": {
				Description: "Email alias goto. Comma separated list of email destinations.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceGlesysEmailAliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EmailAliasParams{
		EmailAlias: d.Get("emailalias").(string),
		GoTo:       d.Get("goto").(string),
	}

	alias, err := client.EmailDomains.CreateAlias(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating alias: %s", err)
	}

	d.SetId(alias.EmailAlias)
	return resourceGlesysEmailAliasRead(ctx, d, m)
}

func resourceGlesysEmailAliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	components := strings.Split(d.Id(), "@")
	domain := components[1]
	listParams := glesys.ListEmailsParams{Filter: d.Id()}
	accounts, err := client.EmailDomains.List(ctx, domain, listParams)
	if err != nil {
		diag.Errorf("Error listing email domains for domain (%s): %s", domain, err)
		return nil
	}
	if len(accounts.EmailAliases) == 1 {
		log.Printf("[INFO] Found alias: %s", accounts.EmailAliases[0].EmailAlias)
		d.Set("emailalias", accounts.EmailAliases[0].EmailAlias)
		d.Set("goto", accounts.EmailAliases[0].GoTo)

		d.SetId(accounts.EmailAliases[0].EmailAlias)
	} else if len(accounts.EmailAliases) < 1 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceGlesysEmailAliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EmailAliasParams{
		EmailAlias: d.Id(),
	}

	if d.HasChange("goto") {
		params.GoTo = d.Get("goto").(string)
	}

	_, err := client.EmailDomains.EditAlias(ctx, params)
	if err != nil {
		return diag.Errorf("Error updating email alias (%s): %s", d.Id(), err)
	}
	return resourceGlesysEmailAliasRead(ctx, d, m)
}

func resourceGlesysEmailAliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.EmailDomains.Delete(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting email alias: %s", err)
	}
	d.SetId("")
	return nil
}
