package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/masterzen/winrm"
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
				Description: "TTL in minutes",
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func recordBuilder(d *schema.ResourceData) Record {
	return Record{
		Dnszone: d.Get("domain").(string),
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Value:   d.Get("value").(string),
		TTL:     d.Get("ttl").(int),
	}
}

func execTemplate(tmpl string, data interface{}) (string, error) {
	t := template.New("tmplate")
	_, err := t.Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v", err)
	}
	var result bytes.Buffer
	t.Execute(&result, data)

	return result.String(), nil
}

func resourceDNSRecordCreate(d *schema.ResourceData, m interface{}) error {
	const tmplpscript string = `
If (!(Get-DnsServerResourceRecord -ZoneName {{.Dnszone}} -Name {{.Name}} -RRType:{{.Rectype}} -ErrorAction SilentlyContinue | ?{$_.HostName -eq '{{.Name}}'})) {
	switch ('{{.Rectype}}') {
	'A' {
	    Add-DnsServerResourceRecordA -ZoneName {{.Dnszone}} -Name {{.Name}} -IPv4Address {{.Value}} -TimeToLive [System.TimeSpan]::FromMinutes({{.TTL}})
	}
	'CName' {
	    Add-DnsServerResourceRecordCName -ZoneName {{.Dnszone}} -Name {{.Name}} -HostNameAlias {{.Value}} -TimeToLive [System.TimeSpan]::FromMinutes({{.Ttl}})
	}
}
Get-DnsServerResourceRecord test.local test123 | select HostName, RecordType, TimeToLive, RecordData | ConvertTo-Json
`
	mutex.Lock()

	client, err := providerConfigure(d)
	if err != nil {
		return fmt.Errorf("Error creating WinRM client: %v", err)
	}
	winrmClient := client.(winrm.Client)

	rec := recordBuilder(d)
	pscript, err := execTemplate(tmplpscript, rec)
	if err != nil {
		return err
	}

	command := winrm.Powershell(pscript)
	out, outerr, exitcode, err := winrmClient.RunWithString(command, "")

	var response recordResponse
	json.Unmarshal([]byte(out), &response)

	if strings.Contains(outerr, "Error") || exitcode != 0 {
		return fmt.Errorf("Error adding new DNS record: %s\nExitcode: %v\n\n", outerr, exitcode)
	}

	d.SetId(response.RecordData.CimInstanceProperties[0])

	mutex.Unlock()
	return nil
}

func resourceDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
