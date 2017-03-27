package dns

import (
	"fmt"
	"strings"
)

// Record containing information regarding DNS record
type Record struct {
	Dnszone  string
	Name     string
	Type     string
	Value    string
	TTL      float64
	ID       string
	NewValue string
	NewTTL   float64
}

// ReadRecords returns all DNS records matching query
func (c *Client) ReadRecords(rec Record) ([]Record, error) {
	// powershell script template to read record from DNS
	const tmplpscript = `
Get-DnsServerResourceRecord -ZoneName {{.Dnszone}}{{ if .Name }} -Name {{.Name}}{{end}} | ?{($_.RecordType -eq 'A' -or $_.RecordType -eq 'CNAME') -and $_.HostName -eq '{{ .Name }}'} | select DistinguishedName, HostName, RecordData, RecordType, TimeToLive | ConvertTo-Json
`

	pscript, err := tmplExec(rec, tmplpscript)
	if err != nil {
		return []Record{}, fmt.Errorf("Creating template: %v", err)
	}
	output, err := c.ExecutePowerShellScript(pscript)
	if err != nil {
		return []Record{}, fmt.Errorf("Running PowerShell script: %v", err)
	}
	if output.stdout == "" {
		return []Record{}, fmt.Errorf("No Record found: %v", rec.Name)
	}
	output.stdout = makeResponseArray(output.stdout)
	resp, err := unmarshalResponse(output.stdout)
	if err != nil {
		return []Record{}, fmt.Errorf("Unmarshalling response: %v", err)
	}
	return *convertResponse(resp, rec), nil
}

// ReadRecord performs DNS Record lookup from server
func (c *Client) ReadRecord(rec Record) (Record, error) {
	// powershell script template to read record from DNS
	const tmplpscript = `
Get-DnsServerResourceRecord -ZoneName {{.Dnszone}}{{ if .Name }} -Name {{.Name}}{{end}} | ?{($_.RecordType -eq 'A' -or $_.RecordType -eq 'CNAME') -and $_.HostName -eq '{{ .Name }}'} | select DistinguishedName, HostName, RecordData, RecordType, TimeToLive | ConvertTo-Json
`

	pscript, err := tmplExec(rec, tmplpscript)
	if err != nil {
		return Record{}, fmt.Errorf("Creating template: %v", err)
	}
	output, err := c.ExecutePowerShellScript(pscript)
	if err != nil {
		return Record{}, fmt.Errorf("Running PowerShell script: %v", err)
	}
	if output.stdout == "" {
		return Record{}, fmt.Errorf("No Record found: %v", rec.Name)
	}
	output.stdout = makeResponseArray(output.stdout)

	resp, err := unmarshalResponse(output.stdout)
	if err != nil {
		return Record{}, fmt.Errorf("Unmarshalling response: %v", err)
	}
	for _, v := range *convertResponse(resp, rec) {
		if v.Value == rec.Value {
			return v, nil
		}
	}
	return Record{}, fmt.Errorf("Record not found: %s", rec.Name)
}

// ReadRecordfromID retrieves specifc DNS record based on record ID
func (c *Client) ReadRecordfromID(recID string) (Record, error) {
	id := strings.Split(recID, "|")
	if len(id) != 3 {
		return Record{}, fmt.Errorf("ID is incorrect")
	}
	rec := Record{
		Dnszone: id[0],
		Name:    id[1],
		Value:   id[2],
	}
	result, err := c.ReadRecords(rec)
	if err != nil {
		return Record{}, fmt.Errorf("Reading record: %v", err)
	}
	for i, v := range result {
		if v.ID == recID {
			return result[i], nil
		}
	}
	return Record{}, fmt.Errorf("Record not found: %v", recID)
}

// CreateRecord creates new DNS records on server
func (c *Client) CreateRecord(rec Record) ([]Record, error) {
	const tmplscriptA = `
Add-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} -A -IPv4Address {{ .Value }} -TimeToLive (New-TimeSpan -Seconds {{ .TTL }})
`
	const tmplscriptCname = `
Add-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} -CName -HostNameAlias {{ .Value }} -TimeToLive (New-TimeSpan -Seconds {{ .TTL }})
`
	var (
		pscript string
		err     error
	)

	if c.RecordExist(rec) {
		return []Record{}, fmt.Errorf("Record already exists: %v", rec)
	}
	rec.ID = fmt.Sprintf("%s|%s|%s", rec.Dnszone, rec.Name, rec.Value)
	switch rec.Type {
	case "A":
		pscript, err = tmplExec(rec, tmplscriptA)
		if err != nil {
			return []Record{}, fmt.Errorf("Creating template: %v", err)
		}
	case "CNAME":
		pscript, err = tmplExec(rec, tmplscriptCname)
		if err != nil {
			return []Record{}, fmt.Errorf("Creating template: %v", err)
		}
	}
	_, err = c.ExecutePowerShellScript(pscript)
	if err != nil {
		return []Record{}, fmt.Errorf("Executing PowerShell script: %v", err)
	}
	record, err := c.ReadRecordfromID(rec.ID)
	if err != nil {
		return []Record{}, fmt.Errorf("Reading record: %v", err)
	}

	var result []Record
	result = append(result, record)

	return result, nil
}

// DeleteRecord deletes DNS record specified
func (c *Client) DeleteRecord(rec Record) error {
	const tmplscriptA string = `
(Get-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }}) | ?{$_.HostName -eq '{{ .Name }}' -and $_.RecordData.IPv4Address -match '{{ .Value }}'} | Remove-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Force
`
	const tmplscriptCname string = `
(Get-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }}) | ?{$_.HostName -eq '{{ .Name }}' -and $_.RecordData.HostNameAlias -match '{{ .Value }}'} | Remove-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Force
`
	var (
		pscript string
		err     error
	)

	if !c.RecordExist(rec) {
		return fmt.Errorf("Record not found: %v", rec)
	}

	switch rec.Type {
	case "A":
		pscript, err = tmplExec(rec, tmplscriptA)
		if err != nil {
			return fmt.Errorf("Creating template: %v", err)
		}
	case "CNAME":
		pscript, err = tmplExec(rec, tmplscriptCname)
		if err != nil {
			return fmt.Errorf("Creating template: %v", err)
		}
	}

	_, err = c.ExecutePowerShellScript(pscript)
	if err != nil {
		return fmt.Errorf("Executing PowerShell script: %v", err)
	}

	return nil
}

// UpdateRecord updates an existing DNS record
func (c *Client) UpdateRecord(rec Record, newValue string, newTTL float64) (Record, error) {
	const tmplscriptA string = `
$old = Get-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} | ?{$_.HostName -eq '{{ .Name }}' -and $_.RecordData.IPv4Address -eq '{{ .Value }}'}
$new = Get-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} | ?{$_.HostName -eq '{{ .Name }}' -and $_.RecordData.IPv4Address -eq '{{ .Value }}'}
{{ if .NewValue -}}
$new.RecordData.IPv4Address = [System.Net.IPAddress]::Parse('{{ .NewValue }}')
{{ end -}}
{{ if ne .NewTTL 0.0 -}}
$new.TimeToLive = New-Timespan -Seconds {{ .NewTTL }}
{{ end -}}
Set-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -NewInputObject $new -OldInputObject $old
`
	const tmplscriptCname string = `
$old = Get-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} | ?{$_.HostName -eq '{{ .Name }}' -and $_.RecordData.HostNameAlias -eq '{{ .Value }}'}
$new = Get-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -Name {{ .Name }} | ?{$_.HostName -eq '{{ .Name }}' -and $_.RecordData.HostNameAlias -eq '{{ .Value }}'}
{{ if .NewValue -}}
$new.RecordData.HostNameAlias = '{{ .NewValue }}'
{{ end -}}
{{ if ne .NewTTL 0.0 -}}
$new.TimeToLive = New-Timespan -Seconds {{ .NewTTL }}
{{ end -}}
Set-DnsServerResourceRecord -ZoneName {{ .Dnszone }} -NewInputObject $new -OldInputObject $old
`
	var (
		pscript string
		err     error
	)

	if !c.RecordExist(rec) {
		return Record{}, fmt.Errorf("Record not found: %v", rec.Name)
	}
	rec, err = c.ReadRecord(rec)
	if err != nil {
		return Record{}, fmt.Errorf("Reading record: %v", err)
	}
	rec.NewValue = newValue
	rec.NewTTL = newTTL
	switch rec.Type {
	case "A":
		pscript, err = tmplExec(rec, tmplscriptA)
		if err != nil {
			return Record{}, fmt.Errorf("Createing template: %v", err)
		}
	case "CNAME":
		pscript, err = tmplExec(rec, tmplscriptCname)
		if err != nil {
			return Record{}, fmt.Errorf("Createing template: %v", err)
		}
	}

	_, err = c.ExecutePowerShellScript(pscript)
	if err != nil {
		return Record{}, fmt.Errorf("Excuting PowerShell script: %v", err)
	}
	if rec.NewValue != "" {
		rec.Value = rec.NewValue
	}

	rec, err = c.ReadRecord(rec)
	if err != nil {
		return Record{}, fmt.Errorf("Reading updated record: %v", err)
	}

	return rec, nil
}

// RecordExist returns if record exists or not
func (c *Client) RecordExist(rec Record) bool {
	var records []Record

	if rec.ID != "" {
		resp, err := c.ReadRecordfromID(rec.ID)
		if err != nil {
			return false
		}
		records = append(records, resp)
	} else {
		records, _ = c.ReadRecords(rec)
	}

	if len(records) > 0 {
		for _, v := range records {
			if v.Value == rec.Value {
				return true
			}
		}
	}
	return false
}
