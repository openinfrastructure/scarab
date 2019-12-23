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
	"log"

	"github.com/openinfrastructure/scarab/common/scarab"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VPN and DNS with new IP",
	Long: `Given a new IP address, update the Cloud VPN Tunnel and Cloud DNS A
record with the new value of the IP address.  The existing tunnel resource is
deleted because the API does not allow patching of an existing tunnel.`,
	Run: updateCmdMain,
}

func updateCmdMain(cmd *cobra.Command, args []string) {
	validateConfig()

	project := viper.GetString("project")
	region := viper.GetString("region")
	tunnel := viper.GetString("tunnel")
	if tunnel == "" {
		log.Fatal("tunnel is required")
	}
	address := viper.GetString("address")
	if address == "" {
		log.Fatal("address is required")
	}

	s := scarab.NewService()

	tun, err := s.VpnTunnels.Get(project, region, tunnel).Do()
	if err != nil {
		log.Fatal(err)
	}

	if scarab.TunnelNeedsUpdate(tun, address) {
		log.Println("Tunnel needs to be udpated:", tun.Status, tun.PeerIp)
		secret := viper.GetString("secret")

		randomSuffix := scarab.RandStr(4)
		name := fmt.Sprintf("%s-%s", tunnel, randomSuffix)

		if err := scarab.CreateTunnel(s, tun, project, region, name, address, secret); err != nil {
			log.Fatal(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Dynamic IP address of the VPN Peer
	updateCmd.Flags().String("address", "", "IP address new value {SCARAB_ADDRESS}")
	if err := viper.BindPFlag("address", updateCmd.Flags().Lookup("address")); err != nil {
		log.Fatal(err)
	}
	// The tunnel to check
	updateCmd.Flags().String("tunnel", "", "VPN Tunnel to update with new IP {SCARAB_TUNNEL}")
	if err := viper.BindPFlag("tunnel", updateCmd.Flags().Lookup("tunnel")); err != nil {
		log.Fatal(err)
	}
	// Shared Secret
	updateCmd.Flags().String("secret", "", "VPN Tunnel shared secret {SCARAB_SECRET}")
	if err := viper.BindPFlag("secret", updateCmd.Flags().Lookup("secret")); err != nil {
		log.Fatal(err)
	}
	// Shared Secret
	updateCmd.Flags().String("secrethash", "", "VPN Tunnel shared secret hash {SCARAB_SECRETHASH}")
	if err := viper.BindPFlag("secrethash", updateCmd.Flags().Lookup("secrethash")); err != nil {
		log.Fatal(err)
	}
}
