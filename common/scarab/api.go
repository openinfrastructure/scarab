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
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewService returns the compute service.
func NewService() *compute.Service {
	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return svc
}

// TunnelNeedsUpdate returns true if the tunnel needs to be updated with the
// correct IP address.
// Possible tun.Status values: "NO_INCOMING_PACKETS",
func TunnelNeedsUpdate(tun *compute.VpnTunnel, ip string) (update bool) {
	update = false
	if ip != tun.PeerIp {
		update = true
	}
	return
}

// CreateExternalGateway creates an external gateway.  An external gateway is
// required when creating a VPN Tunnel.  The external gateway must not be used
// by any other VPN tunnels.
//
// See: https://cloud.google.com/compute/docs/reference/rest/v1/externalVpnGateways/insert
// gcloud compute external-vpn-gateways create hq-b564 --interfaces=0=97.120.179.239
func CreateExternalGateway(svc *compute.Service, project, name, ip string) (*compute.ExternalVpnGateway, error) {
	gwIn := &compute.ExternalVpnGateway{
		Description:    "Created by scarab",
		Name:           name,
		RedundancyType: "SINGLE_IP_INTERNALLY_REDUNDANT",
		Interfaces: []*compute.ExternalVpnGatewayInterface{
			{
				Id:              0,
				IpAddress:       ip,
				ForceSendFields: []string{"Id", "IpAddress"},
			},
		},
	}

	op, err := svc.ExternalVpnGateways.Insert(project, gwIn).Do()
	if err != nil {
		return nil, fmt.Errorf("CreateExternalGateway %v: %w", name, err)
	}
	_, err = waitDone(svc, project, op)
	if err != nil {
		return nil, fmt.Errorf("CreateExternalGateway %v: %w", name, err)
	}
	gwOut, err := svc.ExternalVpnGateways.Get(project, name).Do()
	if err != nil {
		return nil, fmt.Errorf("CreateExternalGateway %v: %w", name, err)
	}
	return gwOut, nil
}

// opItem holds the multiple-values returned from Operations.Get() intended for
// channel communication
type opItem struct {
	op  *compute.Operation
	err error
}

// waitDone waits until an Operation.Status is "DONE", polling the operation
// periodically.
//
// See: https://telliott.io/2016/09/29/three-ish-ways-to-implement-timeouts-in-go.html
func waitDone(svc *compute.Service, project string, opIn *compute.Operation) (opOut *compute.Operation, err error) {
	timeout := time.NewTimer(600 * time.Second)
	ch := make(chan opItem, 1)
	quit := make(chan bool, 1)
	defer close(quit)

	go pollOpItem(svc, project, opIn, ch, quit)

	for {
		select {
		case it := <-ch:
			if it.err != nil {
				log.Println("Retrying:", opIn.Name, it.err)
				continue
			}
			log.Println(it.op.Status, it.op.Progress, it.op.OperationType, it.op.Name, it.op.TargetLink, it.op.Warnings)
			if it.op.Status == "DONE" {
				return it.op, nil
			}
		case <-timeout.C:
			return nil, fmt.Errorf("timeout: %v", opIn.Name)
		}
	}
}

// pollOpItem periodically gets an opItem and send it to ch until told to quit.
func pollOpItem(svc *compute.Service, project string, opIn *compute.Operation, ch chan<- opItem, quit <-chan bool) {
	defer close(ch)
	for {
		select {
		case <-quit:
			return
		case <-time.NewTimer(1 * time.Second).C:
			ch <- getOpItem(svc, project, opIn)
		}
	}
}

// getOpItem calls Operations.Get() once and returns an opItem
func getOpItem(svc *compute.Service, project string, opIn *compute.Operation) opItem {
	it := opItem{}

	if opIn.Zone != "" {
		sl := strings.Split(opIn.Zone, "/")
		zone := sl[len(sl)-1]
		it.op, it.err = svc.ZoneOperations.Get(project, zone, opIn.Name).Do()
	} else if opIn.Region != "" {
		sl := strings.Split(opIn.Region, "/")
		region := sl[len(sl)-1]
		it.op, it.err = svc.RegionOperations.Get(project, region, opIn.Name).Do()
	} else {
		it.op, it.err = svc.GlobalOperations.Get(project, opIn.Name).Do()
	}
	return it
}

// CreateTunnel creates a new VPN Tunnel given a reference VpnTunnel instance
// See https://cloud.google.com/compute/docs/reference/rest/v1/vpnTunnels/insert
func CreateTunnel(svc *compute.Service, ref *compute.VpnTunnel, project, region, name, ip, secret string) error {
	tun := new(compute.VpnTunnel)
	tun.Name = name
	tun.SharedSecret = secret
	tun.Description = ref.Description
	tun.IkeVersion = ref.IkeVersion
	tun.Router = ref.Router
	tun.VpnGateway = ref.VpnGateway
	tun.VpnGatewayInterface = ref.VpnGatewayInterface

	extGw, err := CreateExternalGateway(svc, project, name, ip)
	if err != nil {
		log.Fatal(fmt.Errorf("CreateExternalGateway %v: %w", name, err))
	}

	tun.PeerExternalGateway = extGw.SelfLink
	tun.PeerExternalGatewayInterface = 0

	// Fix for https://github.com/googleapis/google-api-go-client/issues/413
	tun.ForceSendFields = []string{"VpnGatewayInterface", "PeerExternalGatewayInterface"}

	op, err := svc.VpnTunnels.Insert(project, region, tun).Do()
	if err != nil {
		return fmt.Errorf("CreateTunnel %v: %w", name, err)
	}
	_, err = waitDone(svc, project, op)
	if err != nil {
		return fmt.Errorf("CreateTunnel %v: %w", name, err)
	}
	return nil
}

// RandStr returns a random string intended for GCP resource name suffixes
func RandStr(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// Don't allow digits for the first character
	sb.WriteByte(chars[rand.Int63()%int64(len(chars)-10)])
	// Allow digits for the remainder of the string
	for i := 1; i < n; i++ {
		sb.WriteByte(chars[rand.Int63()%int64(len(chars))])
	}
	return sb.String()
}

// VpnTunnel returns a pointer to the named VPN tunnel
// func VpnTunnel(name, region, project string, s *compute.Service) (*compute.VpnTunnel, error) {
// }
