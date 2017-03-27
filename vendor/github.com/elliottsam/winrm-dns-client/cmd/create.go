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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elliottsam/winrm-dns-client/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a DNS record on Windows server",
	Long: `Creates a DNS record on Windows server
	configured in the configuration file, all flags except
	TTL are required.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dnsZone == "" || name == "" || recordType == "" || value == "" {
			fmt.Println("Please provide all required parameters")
			os.Exit(1)
		}

		rec := dns.Record{
			Dnszone: dnsZone,
			Name:    name,
			Type:    strings.ToUpper(recordType),
			Value:   value,
			TTL:     ttl,
		}
		ClientConfig := dns.GenerateClient(viper.GetString("servername"), viper.GetString("username"), viper.GetString("password"))
		ClientConfig.ConfigureWinRMClient()

		record, err := ClientConfig.CreateRecord(rec)
		if err != nil {
			log.Fatalln("Error creating DNS record:", err)
		}
		dns.OutputTable(record)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&dnsZone, "DnsZone", "d", "", "DNS Zone to create record for, this is required")
	createCmd.PersistentFlags().StringVarP(&name, "Name", "n", "", "Name of record to create, this is required")
	createCmd.PersistentFlags().StringVarP(&recordType, "Type", "t", "", "Type of DNS record to create, this is required")
	createCmd.PersistentFlags().StringVarP(&value, "Value", "v", "", "Value of DNS record, this is required")
	createCmd.PersistentFlags().Float64VarP(&ttl, "TTL", "l", 900, "TTL for record in seconds")
	createCmd.MarkPersistentFlagRequired("DnsZone")
}
