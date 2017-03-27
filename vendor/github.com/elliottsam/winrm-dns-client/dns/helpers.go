package dns

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
)

func tmplExec(r Record, tp string) (string, error) {
	t := template.New("tmpl")
	t, err := t.Parse(tp)
	if err != nil {
		return "", fmt.Errorf("Parsing template: %v", err)
	}
	var result bytes.Buffer
	if err := t.Execute(&result, r); err != nil {
		return "", fmt.Errorf("Generating template: %v", err)
	}

	return result.String(), nil
}

func unmarshalResponse(resp string) ([]interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(resp), &data); err != nil {
		return nil, fmt.Errorf("Unmarshalling json: %v", err)
	}
	return data.([]interface{}), nil
}

func convertResponse(r []interface{}, origrec Record) *[]Record {
	records := []Record{}
	for i := range r {
		var rec Record
		switch r[i].(map[string]interface{})["RecordData"].(map[string]interface{})["CimInstanceProperties"].(type) {
		case []interface{}:
			rec = Record{
				Dnszone: origrec.Dnszone,
				Name:    r[i].(map[string]interface{})["HostName"].(string),
				Type:    r[i].(map[string]interface{})["RecordType"].(string),
				Value:   strings.Split(r[i].(map[string]interface{})["RecordData"].(map[string]interface{})["CimInstanceProperties"].([]interface{})[0].(string), "\"")[1],
				TTL:     r[i].(map[string]interface{})["TimeToLive"].(map[string]interface{})["TotalSeconds"].(float64),
			}
			rec.ID = fmt.Sprintf("%s|%s|%s", rec.Dnszone, rec.Name, rec.Value)
		case string:
			rec = Record{
				Dnszone: origrec.Dnszone,
				Name:    r[i].(map[string]interface{})["HostName"].(string),
				Type:    r[i].(map[string]interface{})["RecordType"].(string),
				Value:   strings.Split(r[i].(map[string]interface{})["RecordData"].(map[string]interface{})["CimInstanceProperties"].(string), "\"")[1],
				TTL:     r[i].(map[string]interface{})["TimeToLive"].(map[string]interface{})["TotalSeconds"].(float64),
			}
			rec.ID = fmt.Sprintf("%s|%s|%s", rec.Dnszone, rec.Name, rec.Value)
		}
		records = append(records, rec)

	}
	return &records
}

func makeResponseArray(r string) string {
	if rune(r[0]) != '[' && rune(r[(len(r)-1)]) != ']' {
		return fmt.Sprintf("[%s]", r)
	}
	return r
}

//OutputTable a table containing DNS entries
func OutputTable(rec []Record) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"DnsZone", "Name", "Type", "Value", "TTL", "ID"})
	for _, v := range rec {
		table.Append([]string{v.Dnszone, v.Name, v.Type, v.Value, fmt.Sprint(v.TTL), v.ID})
	}
	table.Render()
}

func powershell(psCmd string) string {
	wideCmd := ""
	for _, b := range []byte(psCmd) {
		wideCmd += string(b) + "\x00"
	}
	input := []uint8(wideCmd)
	encodedCmd := base64.StdEncoding.EncodeToString(input)
	return fmt.Sprintf("powershell.exe -EncodedCommand %s", encodedCmd)
}
