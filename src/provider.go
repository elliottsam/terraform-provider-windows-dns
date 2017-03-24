package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider allows making changes to Windows DNS server
// Utilises Powershell to connect to domain controller
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINRM_SERVER", nil),
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINRM_USERNAME", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINRM_PASSWORD", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"windows-dns_record": resourceDNSRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := config{
		ServerName: d.Get("server_name").(string),
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
	}

	return config.Client()
}
