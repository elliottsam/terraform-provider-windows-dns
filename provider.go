package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider allows making changes to Windows DNS server
// Utilises Powershell to connect to domain controller
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
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

			"https": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINRM_HTTPS", false),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("WINRM_SECURE", false),
			},
			"port": {
				Type:		schema.TypeInt,
				Optional:	true,
				DefaultFunc: schema.EnvDefaultFunc("WINRM_PORT", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"windowsdns_record": resourceDNSRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var port int
	if  d.Get("port") == nil || d.Get("port") == 0 {
		if d.Get("https").(bool) {
			port = 5986
		} else {
			port = 5985
		}
	} else {
		port = d.Get("port").(int)
	}

	config := config{
		ServerName: d.Get("server").(string),
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		HTTPS:		d.Get("https").(bool),
		Insecure:	d.Get("insecure").(bool),
		Port:		port,
	}

	return config.Client()
}
