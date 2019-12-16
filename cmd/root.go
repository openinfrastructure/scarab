/*
Copyright © 2019 Open Infrastructure Services, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/openinfrastructure/scarab/common/scarab"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scarab",
	Short: "Update a Cloud VPN tunnel",
	Long: `Update a vpn tunnel with the IP address provided.

Command line flags may be set via the environment, or config file.
Environment variables listed leftmost take precedence.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = scarab.Version.String()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scarab.yaml)")
	// Project ID flag
	rootCmd.PersistentFlags().String("project", "", "The GCP project id {SCARAB_PROJECT, CLOUDSDK_CORE_PROJECT}")
	err := viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	if err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".scarab" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".scarab")
	}

	viper.SetEnvPrefix("scarab")
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.BindEnv("project", "CLOUDSDK_CORE_PROJECT"); err != nil {
		log.Fatal(err)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// validateConfig validates required configuration values
func validateConfig() {
	if viper.Get("project") == "" {
		log.Fatal("Error: --project is required")
	}
}
