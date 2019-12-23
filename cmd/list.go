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
	"log"

	"github.com/openinfrastructure/scarab/common/scarab"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Cloud VPN Tunnels and their Peer IP.",
	Long:  "List Cloud VPN Tunnels and their Peer IP address.",
	Run:   listCmdMain,
}

// listCmdMain is the main() function for the list command.  It obtains a list
// of Cloud VPN tunnels for a given project and region.
func listCmdMain(cmd *cobra.Command, args []string) {
	validateConfig()

	project := viper.GetString("project")
	region := viper.GetString("region")
	service := scarab.NewService()

	tuns, err := service.VpnTunnels.List(project, region).Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, tun := range tuns.Items {
		log.Println(tun.Name, tun.PeerIp, tun.DetailedStatus)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
