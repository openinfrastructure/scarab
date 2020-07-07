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

var dnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS managed zones",
	Long:  "List DNS managed zones in the specified project",
	Run:   dnsListCmdMain,
}

func init() {
	dnsCmd.AddCommand(dnsListCmd)
}

func dnsListCmdMain(cmd *cobra.Command, args []string) {
	validateConfig()

	project := viper.GetString("project")
	svc := scarab.NewDNSService()

	l, err := svc.ManagedZones.List(project).Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range l.ManagedZones {
		log.Println(i.Name)
	}
}
