package main

import "github.com/masterzen/winrm"

type Config struct {
	ServerName string
	Username string
	Password string
}

type Record struct {
	Dnszone string
	Name    string
	Type    string
	Value   string
	TTL     int
}

type recordResponse struct {
	HostName   string `json:"HostName"`
	RecordData struct {
		CimClass struct {
			CimClassMethods     string `json:"CimClassMethods"`
			CimClassProperties  string `json:"CimClassProperties"`
			CimClassQualifiers  string `json:"CimClassQualifiers"`
			CimSuperClass       string `json:"CimSuperClass"`
			CimSuperClassName   string `json:"CimSuperClassName"`
			CimSystemProperties string `json:"CimSystemProperties"`
		} `json:"CimClass"`
		CimInstanceProperties []string `json:"CimInstanceProperties"`
		CimSystemProperties   struct {
			ClassName  string      `json:"ClassName"`
			Namespace  string      `json:"Namespace"`
			Path       interface{} `json:"Path"`
			ServerName string      `json:"ServerName"`
		} `json:"CimSystemProperties"`
	} `json:"RecordData"`
	RecordType string `json:"RecordType"`
	TimeToLive struct {
		Days              int64   `json:"Days"`
		Hours             int64   `json:"Hours"`
		Milliseconds      int64   `json:"Milliseconds"`
		Minutes           int64   `json:"Minutes"`
		Seconds           int64   `json:"Seconds"`
		Ticks             int64   `json:"Ticks"`
		TotalDays         float64 `json:"TotalDays"`
		TotalHours        int64   `json:"TotalHours"`
		TotalMilliseconds int64   `json:"TotalMilliseconds"`
		TotalMinutes      int64   `json:"TotalMinutes"`
		TotalSeconds      int64   `json:"TotalSeconds"`
	} `json:"TimeToLive"`
}

// Configures the WinRM endpoint for managing Microsoft DNS
func (c * Config) Client() (*winrm.Client, error) {
	endpoint := winrm.NewEndpoint(c.ServerName, 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, c.Username, c.Password)
	if err != nil {
		return nil, err
	}

	return client, nil
}
