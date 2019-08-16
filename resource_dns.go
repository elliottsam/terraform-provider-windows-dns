package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Bourne-ID/winrm-dns-client/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

var mutex = &sync.Mutex{}

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordUpdate,
		Delete: resourceDNSRecordDelete,
		Exists: resourceDNSRecordExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TTL as a duration",
				Default:     "15m0s",
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSRecordCreate(d *schema.ResourceData, m interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	client := m.(*dns.Client)

	ttl, err := time.ParseDuration(d.Get("ttl").(string))
	if err != nil {
		return fmt.Errorf("Invalid time duration: %v", err)
	}

	rec := dns.Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
		TTL:     ttl.Seconds(),
	}

	resp, err := client.CreateRecord(rec)
	if err != nil {
		return err
	}

	d.SetId(resp[0].ID)
	return nil
}

func resourceDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	client := m.(*dns.Client)

	var err error
	rec := dns.Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
		ID:      d.Id(),
	}

	if rec.ID != "" {
		rec, err = client.ReadRecordfromID(rec.ID)
		if err != nil {
			return err
		}
	} else {
		rec, err = client.ReadRecord(rec)
		if err != nil {
			return err
		}
	}

	ttl, err := time.ParseDuration(fmt.Sprintf("%vs", rec.TTL))
	if err != nil {
		fmt.Errorf("Invalid time duration: %v", err)
	}

	d.Set("domain", rec.Dnszone)
	d.Set("fqdn", fmt.Sprintf("%s.%s", rec.Name, rec.Dnszone))
	d.Set("name", rec.Name)
	d.Set("type", rec.Type)
	d.Set("value", rec.Value)
	d.Set("ttl", ttl.String())
	d.SetId(rec.ID)

	return nil
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	var (
		err      error
		newValue string
		newTTL   time.Duration
	)
	client := m.(*dns.Client)

	rec, err := client.ReadRecordfromID(d.Id())
	if err != nil {
		return fmt.Errorf("Error reading record: %v", err)
	}

	oldval, newval := d.GetChange("value")
	if oldval != newval {
		newValue = d.Get("value").(string)
	}

	oldval, newval = d.GetChange("ttl")
	if oldval != newval {
		newTTL, err = time.ParseDuration(d.Get("ttl").(string))
		if err != nil {
			return fmt.Errorf("Invalid time duration: %v", err)
		}
	}

	rec, err = client.UpdateRecord(rec, newValue, newTTL.Seconds())
	if err != nil {
		return fmt.Errorf("Error updating record: %v", err)
	}

	d.SetId(rec.ID)

	return nil
}

func resourceDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	client := m.(*dns.Client)

	rec := dns.Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
	}

	if err := client.DeleteRecord(rec); err != nil {
		return fmt.Errorf("Error deleting record: %v", err)
	}

	return nil
}

func resourceDNSRecordExists(d *schema.ResourceData, m interface{}) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()
	client := m.(*dns.Client)

	rec := dns.Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
		ID:      d.Id(),
	}

	if !client.RecordExist(rec) {
		return false, fmt.Errorf("Record not found: %v", rec)
	}

	return true, nil
}
