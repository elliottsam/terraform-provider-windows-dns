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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	dnsZone    string
	name       string
	recordType string
	value      string
	ttl        float64
	id         string
	newvalue   string
	newttl     float64
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "winrm-dns-client",
	Short: "CLI tool for interacting with Microsoft DNS Servers",
	Long: `winrm-dns-client allows retrieving, creating and updating
	of DNS records on Microsoft DNS servers

	At present this requires the server to have WinRM configured and
	also be running at least Windows Server 2012 and have the dnsserver
	PowerShell module installed`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.winrm-dns-client.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".winrm-dns-client") // name of config file (without extension)
	viper.AddConfigPath("$HOME")             // adding home directory as first search path
	viper.AutomaticEnv()                     // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

}
