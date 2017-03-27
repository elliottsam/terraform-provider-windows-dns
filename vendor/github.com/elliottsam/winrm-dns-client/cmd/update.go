// Copyright © 2017 Sam Elliott <me@sam-e.co.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"strings"

	"github.com/elliottsam/winrm-dns-client/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update DNS record",
	Long:  `Update existing DNS record with new value`,
	Run: func(cmd *cobra.Command, args []string) {
		var rec dns.Record
		if id == "" && (dnsZone == "" || name == "" || value == "") {
			log.Fatal("Please specify ID or DnsZone, Name and Value as parameters")
		}

		if id != "" {
			rec.Dnszone = strings.Split(id, "|")[0]
			rec.Name = strings.Split(id, "|")[1]
			rec.Value = strings.Split(id, "|")[2]
		} else {
			rec.Dnszone = dnsZone
			rec.Name = name
			rec.Value = value
		}

		ClientConfig := dns.GenerateClient(viper.GetString("servername"), viper.GetString("username"), viper.GetString("password"))
		ClientConfig.ConfigureWinRMClient()

		rec, err := ClientConfig.UpdateRecord(rec, newvalue, newttl)
		if err != nil {
			log.Fatalf("Error updating record: %v", err)
		}

		var output []dns.Record
		output = append(output, rec)
		dns.OutputTable(output)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&dnsZone, "DnsZone", "d", "", "DNS Zone to create record for, this is required")
	updateCmd.PersistentFlags().StringVarP(&name, "Name", "n", "", "Name of record to create, this is required")
	updateCmd.PersistentFlags().StringVarP(&recordType, "Type", "t", "", "Type of DNS record to create, this is required")
	updateCmd.PersistentFlags().StringVarP(&value, "Value", "v", "", "Value of DNS record, this is required")
	updateCmd.PersistentFlags().StringVarP(&newvalue, "NewValue", "u", "", "Value to update DNS record")
	updateCmd.PersistentFlags().Float64VarP(&newttl, "NewTTL", "l", 0, "New TTL value to update DNS record")
	updateCmd.PersistentFlags().StringVarP(&id, "ID", "i", "", "ID of record to update")

}
