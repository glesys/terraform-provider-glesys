package glesys

import (
	"context"
	"log"
	"strings"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGlesysEmailAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysEmailAccountCreate,
		ReadContext:   resourceGlesysEmailAccountRead,
		UpdateContext: resourceGlesysEmailAccountUpdate,
		DeleteContext: resourceGlesysEmailAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Create a GleSYS Email Account.",

		Schema: map[string]*schema.Schema{
			"emailaccount": {
				Description: "Email account name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "Email Account password",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"antispamlevel": {
				Description: "Email Account antispam level. `0-5`",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"antivirus": {
				Description: "Email Account enable Antivirus. `yes/no`",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes", "no"}, false),
			},
			"autorespond": {
				Description: "Email Account Autoresponse. `yes/no`",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes", "no"}, false),
			},
			"autorespondmessage": {
				Description: "Email Account Autoresponse message.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"autorespondsaveemail": {
				Description: "Email Account Save emails on autorespond.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"rejectspam": {
				Description: "Email Account Reject spam setting. `yes/no`",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes", "no"}, false),
			},
			"quotaingib": {
				Description: "Email Account Quota (GiB)",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"created": {
				Description: "Email Account created date",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"displayname": {
				Description: "Email Account displayname",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified": {
				Description: "Email Account modification date",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGlesysEmailAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.CreateAccountParams{
		EmailAccount:       d.Get("emailaccount").(string),
		AntiVirus:          d.Get("antivirus").(string),
		AntiSpamLevel:      d.Get("antispamlevel").(int),
		AutoRespond:        d.Get("autorespond").(string),
		AutoRespondMessage: d.Get("autorespondmessage").(string),
		QuotaInGiB:         d.Get("quotaingib").(int),
		RejectSpam:         d.Get("rejectspam").(string),
	}

	account, err := client.EmailDomains.CreateAccount(ctx, params)
	if err != nil {
		return diag.Errorf("Error creating account: %s", err)
	}

	d.SetId(account.EmailAccount)
	return resourceGlesysEmailAccountRead(ctx, d, m)
}

func resourceGlesysEmailAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	components := strings.Split(d.Id(), "@")
	domain := components[1]
	listParams := glesys.ListEmailsParams{Filter: d.Id()}
	accounts, err := client.EmailDomains.List(ctx, domain, listParams)
	if err != nil {
		diag.Errorf("Error listing email domains for domain (%s): %s", domain, err)
		return nil
	}
	if len(accounts.EmailAccounts) == 1 {
		log.Printf("[INFO] Found account: %s", accounts.EmailAccounts[0].EmailAccount)
		d.Set("antispamlevel", accounts.EmailAccounts[0].AntiSpamLevel)
		d.Set("antivirus", accounts.EmailAccounts[0].AntiVirus)
		d.Set("autorespond", accounts.EmailAccounts[0].AutoRespond)
		d.Set("autorespondmessage", accounts.EmailAccounts[0].AutoRespondMessage)
		d.Set("autorespondsaveemail", accounts.EmailAccounts[0].AutoRespondSaveEmail)
		d.Set("created", accounts.EmailAccounts[0].Created)
		d.Set("displayname", accounts.EmailAccounts[0].DisplayName)
		d.Set("emailaccount", accounts.EmailAccounts[0].EmailAccount)
		d.Set("modified", accounts.EmailAccounts[0].Modified)
		d.Set("quotaingib", accounts.EmailAccounts[0].QuotaInGiB)
		d.Set("rejectspam", accounts.EmailAccounts[0].RejectSpam)

		d.SetId(accounts.EmailAccounts[0].EmailAccount)
	} else if len(accounts.EmailAccounts) < 1 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceGlesysEmailAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.EditAccountParams{}

	if d.HasChange("antispamlevel") {
		params.AntiSpamLevel = d.Get("antispamlevel").(int)
	}

	if d.HasChange("antivirus") {
		params.AntiVirus = d.Get("antivirus").(string)
	}

	if d.HasChange("autorespond") {
		params.AutoRespond = d.Get("autorespond").(string)
	}

	if d.HasChange("autorespondmessage") {
		params.AutoRespondMessage = d.Get("autorespondmessage").(string)
	}

	if d.HasChange("quotaingib") {
		params.QuotaInGiB = d.Get("quotaingib").(int)
	}

	if d.HasChange("rejectspam") {
		params.RejectSpam = d.Get("rejectspam").(string)
	}

	_, err := client.EmailDomains.EditAccount(ctx, d.Id(), params)
	if err != nil {
		return diag.Errorf("Error updating email account (%s): %s", d.Id(), err)
	}
	return resourceGlesysEmailAccountRead(ctx, d, m)
}

func resourceGlesysEmailAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	err := client.EmailDomains.Delete(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting email account: %s", err)
	}
	d.SetId("")
	return nil
}
