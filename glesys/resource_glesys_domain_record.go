package glesys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/glesys/glesys-go/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysDNSDomainRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysDNSDomainRecordCreate,
		ReadContext:   resourceGlesysDNSDomainRecordRead,
		UpdateContext: resourceGlesysDNSDomainRecordUpdate,
		DeleteContext: resourceGlesysDNSDomainRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGlesysRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"data": {
				Description: "Record data field. Ex. `127.0.0.1`",
				Type:        schema.TypeString,
				Required:    true,
			},

			"domain": {
				Description: "Domain name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"host": {
				Description: "Record host field. Ex. `www`",
				Type:        schema.TypeString,
				Required:    true,
			},

			"recordid": {
				Description: "Record internal id",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"type": {
				Description: "Record type. Must be one of `SOA`, `A`, `AAAA`, `CAA`, `CNAME`, `MX`, `NS`, `TXT`, `SRV`, `URL` or `PTR`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"ttl": {
				Description: "Record TTL field",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// resourceGlesysRecordImport - import records "domain.tld,123456"
func resourceGlesysRecordImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		_, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid recordid: %v", err)
		}

		d.SetId(s[1])
		d.Set("domain", s[0])
	}

	return []*schema.ResourceData{d}, nil
}

func resourceGlesysDNSDomainRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	params := glesys.AddRecordParams{
		Data:       d.Get("data").(string),
		DomainName: d.Get("domain").(string),
		Host:       d.Get("host").(string),
		Type:       d.Get("type").(string),
		TTL:        d.Get("ttl").(int),
	}

	record, err := client.DNSDomains.AddRecord(ctx, params)
	if err != nil {
		return diag.Errorf("Error adding record \"%s\": %v", params.Data, err)
	}

	// Set the Id to domain.ID
	id := strconv.Itoa(record.RecordID)
	d.SetId(id)

	return resourceGlesysDNSDomainRecordRead(ctx, d, m)
}

func resourceGlesysDNSDomainRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	domain := d.Get("domain").(string)
	myID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid record id: %v", err)
	}

	record, err := findRecordByID(client, domain, myID)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("domain", record.DomainName)
	d.Set("data", record.Data)
	d.Set("host", record.Host)
	d.Set("recordid", record.RecordID)
	d.Set("ttl", record.TTL)
	d.Set("type", record.Type)

	return nil
}

func resourceGlesysDNSDomainRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	myID := d.Id()
	recordID, errid := strconv.Atoi(myID)
	if errid != nil {
		return diag.Errorf("Id must be converted to integer: %v", errid)
	}
	params := glesys.UpdateRecordParams{RecordID: recordID}

	if d.HasChange("data") {
		params.Data = d.Get("data").(string)
	}

	if d.HasChange("host") {
		params.Host = d.Get("host").(string)
	}

	if d.HasChange("ttl") {
		params.TTL = d.Get("ttl").(int)
	}

	if d.HasChange("type") {
		params.Type = d.Get("type").(string)
	}

	_, err := client.DNSDomains.UpdateRecord(ctx, params)
	if err != nil {
		return diag.Errorf("Error updating record: %v", err)
	}

	return resourceGlesysDNSDomainRecordRead(ctx, d, m)
}

func resourceGlesysDNSDomainRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	recordID, errid := strconv.Atoi(d.Id())
	if errid != nil {
		return diag.Errorf("Id must be converted to integer: %v", errid)
	}

	err := client.DNSDomains.DeleteRecord(ctx, recordID)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP error: 404") {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error deleting domain record: %v", err)
		}
	}
	d.SetId("")
	return nil
}

func findRecordByID(client *glesys.Client, domain string, id int) (*glesys.DNSDomainRecord, error) {
	records, err := client.DNSDomains.ListRecords(context.Background(), domain)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %s", err)
	}

	for _, rec := range *records {
		if rec.RecordID == id {
			return &rec, nil
		}
	}

	return nil, fmt.Errorf("no record found for ID %d", id)
}
