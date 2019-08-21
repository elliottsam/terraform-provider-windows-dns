package main

import (
	"github.com/Bourne-ID/winrm-dns-client/dns"
)

type config struct {
	ServerName string
	Username   string
	Password   string
	Port       int
	HTTPS      bool
	Insecure   bool
}

// Client configures the WinRM endpoint for managing Microsoft DNS
func (c *config) Client() (*dns.Client, error) {
	client := dns.Client{
		ServerName: c.ServerName,
		Username:   c.Username,
		Password:   c.Password,
		Port:		c.Port,
		HTTPS:		c.HTTPS,
		Insecure:	c.Insecure,
	}

	if err := client.ConfigureWinRMClient(); err != nil {
		return nil, err
	}

	return &client, nil
}
