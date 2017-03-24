package main

import (
	"fmt"
	"sync"

	"github.com/elliottsam/winrm-dns-client/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

var mutex = &sync.Mutex{}

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordUpdate,
		Delete: resourceDNSRecordDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "TTL in seconds",
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

	rec := dns.Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
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

	if !client.RecordExist(rec) {
		return fmt.Errorf("Record not found: %v", rec.Name)
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

	d.Set("domain", rec.Dnszone)
	d.Set("fqdn", fmt.Sprintf("%s.%s", rec.Name, rec.Dnszone))
	d.Set("name", rec.Name)
	d.Set("type", rec.Type)
	d.Set("value", rec.Value)
	d.Set("ttl", rec.TTL)

	return nil
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	mutex.Lock()
	defer mutex.Unlock()
	var (
		newValue string
		newTTL   int64
	)
	client := m.(*dns.Client)

	rec, err := client.ReadRecordfromID(d.Get("Id").(string))
	if err != nil {
		return fmt.Errorf("Error reading record: %v", err)
	}

	if d.Get("value").(string) != rec.Value {
		newValue = d.Get("value").(string)
	}
	if d.Get("ttl").(int64) != rec.TTL {
		newTTL = d.Get("value").(int64)
	}

	rec, err = client.UpdateRecord(rec, newValue, newTTL)
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

	if !client.RecordExist(rec) {
		return fmt.Errorf("Record not found: %s", rec.Name)
	}

	if err := client.DeleteRecord(rec); err != nil {
		return fmt.Errorf("Error deleting record: %v", err)
	}

	return nil
}
