package main

import (
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
				Required:    false,
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

	rec := dns.Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
		ID:      d.Get("value").(string),
	}

	rec, err := client.ReadRecord(rec)
	if err != nil {
		return err
	}

	return nil
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
