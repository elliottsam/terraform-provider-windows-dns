package main

import (
	"github.com/elliottsam/winrm-dns-client/dns"
)

type Config struct {
	ServerName string
	Username   string
	Password   string
}

type Record struct {
	Dnszone string
	Name    string
	Type    string
	Value   string
	TTL     int
}

// Configures the WinRM endpoint for managing Microsoft DNS
func (c *Config) Client() (*dns.Client, error) {
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
