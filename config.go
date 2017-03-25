package main

import (
	"github.com/elliottsam/winrm-dns-client/dns"
)

type config struct {
	ServerName string
	Username   string
	Password   string
}

// Client configures the WinRM endpoint for managing Microsoft DNS
func (c *config) Client() (*dns.Client, error) {
	client := dns.Client{
		ServerName: c.ServerName,
		Username:   c.Username,
		Password:   c.Password,
	}
	if err := client.ConfigureWinRMClient(); err != nil {
		return nil, err
	}

	return &client, nil
}
