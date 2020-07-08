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
	"google.golang.org/api/dns/v1"
)

var dnsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update DNS records",
	Long:  "Update DNS resource records in the specified zone atomically",
	Run:   dnsUpdateCmdMain,
}

func init() {
	dnsCmd.AddCommand(dnsUpdateCmd)

	// DNS Manged Zone to change
	dnsUpdateCmd.PersistentFlags().String("dnszone", "", "The managed zone {SCARAB_DNSZONE}")
	if err := viper.BindPFlag("dnszone", dnsUpdateCmd.PersistentFlags().Lookup("dnszone")); err != nil {
		log.Fatal(err)
	}

	// Resource Record names to update (for example, the fqdn)
	dnsUpdateCmd.PersistentFlags().StringSlice("rrnames", []string{}, "The record names to update (fqdns) {SCARAB_RRNAMES}")
	if err := viper.BindPFlag("rrnames", dnsUpdateCmd.PersistentFlags().Lookup("rrnames")); err != nil {
		log.Fatal(err)
	}

	// Resource Record type, for example "A" for an A record.
	dnsUpdateCmd.PersistentFlags().String("rrtype", "A", "The DNS record type {SCARAB_RRTYPE}")
	if err := viper.BindPFlag("rrtype", dnsUpdateCmd.PersistentFlags().Lookup("rrtype")); err != nil {
		log.Fatal(err)
	}

	// Resource Record TTL
	dnsUpdateCmd.PersistentFlags().Int64("rrttl", 60, "The DNS record TTL {SCARAB_RRTTL}")
	if err := viper.BindPFlag("rrttl", dnsUpdateCmd.PersistentFlags().Lookup("rrttl")); err != nil {
		log.Fatal(err)
	}

	// Resource Record data values to update to (for example, the IP address)
	dnsUpdateCmd.PersistentFlags().StringSlice("rrdatas", []string{}, "The new record values (addresses).  If not provided, remove the record. {SCARAB_RRDATAS}")
	if err := viper.BindPFlag("rrdatas", dnsUpdateCmd.PersistentFlags().Lookup("rrdatas")); err != nil {
		log.Fatal(err)
	}
}

/* dnsUpdateCmdMain Atomically updates the ResourceRecordSet collection using
the [DNS v1 Changes
API](https://cloud.google.com/dns/docs/reference/v1/changes).

TODO: Help the user if they forget a trailing dot (.) on the FQDN, otherwise they get:
Reason: invalid, Message: Invalid value for 'entity.change.additions[0].name': 'jeff.ois.run'
Reason: invalid, Message: Invalid value for 'entity.change.additions[1].name': 'foo.jeff.ois.run'
*/
func dnsUpdateCmdMain(cmd *cobra.Command, args []string) {
	validateConfig()

	project := viper.GetString("project")
	dnszone := viper.GetString("dnszone")
	rrnames := viper.GetStringSlice("rrnames")
	rrdatas := viper.GetStringSlice("rrdatas")
	rrtype := viper.GetString("rrtype")
	rrttl := viper.GetInt64("rrttl")

	svc := scarab.NewDNSService()

	if len(rrnames) < 1 {
		log.Fatalf("At least one value for rrnames is required.")
	}

	// See: https://cloud.google.com/dns/docs/reference/v1/changes#resource
	change := &dns.Change{}

	for _, name := range rrnames {
		if len(rrdatas) > 0 {
			a := &dns.ResourceRecordSet{
				Name:    name,
				Rrdatas: rrdatas,
				Ttl:     rrttl,
				Type:    rrtype,
			}
			change.Additions = append(change.Additions, a)
		}

		// If the name already exists it needs to be added to the deletions field
		// of the Change resource.
		rr, err := svc.ResourceRecordSets.List(project, dnszone).Name(name).Type(rrtype).Do()
		if err != nil {
			log.Fatal(err)
		}
		// Add existing records to the list of records to delete in the change transaction.
		for _, rrset := range rr.Rrsets {
			change.Deletions = append(change.Deletions, rrset)
		}
	}

	_, err := scarab.DoDNSChange(svc, project, dnszone, change)
	if err != nil {
		log.Fatal(err)
	}
}

func debug(c *dns.Change) bool {
	log.Println(c)
	return false
}
