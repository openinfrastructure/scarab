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

package scarab

import (
	"context"
	"log"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
)

// NewService returns the GCP Compute API.
// See: https://godoc.org/google.golang.org/api/compute/v1
func NewService() *compute.Service {
	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return svc
}

// NewDNSService returns the GCP DNS API.  See:
// https://godoc.org/google.golang.org/api/dns/v1
func NewDNSService() *dns.Service {
	svc, err := dns.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return svc
}

// DoDNSChange executes a DNS Change requiest if and only if the change
// resource has at least one insertion or at least one deletion.
func DoDNSChange(svc *dns.Service, project string, dnszone string, change *dns.Change) (c *dns.Change, err error) {
	if (len(change.Additions) < 1) && (len(change.Deletions) < 1) {
		log.Println("No changes to make")
		return change, nil
	}
	c, err = svc.Changes.Create(project, dnszone, change).Do()
	log.Println("Change status:", c.Status)
	return c, err
}
