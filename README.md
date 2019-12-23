Scarab
===

Scarab updates a [Cloud VPN][1] peer gateway with the IP address provided.
The intended use case is to re-establish a site-to-site VPN when a dynamic IP
address changes.

The name comes from the Golden Scarab in the [story about syncronicity][2].

Use case
---

When a new public IP is obtained, the site-to-site VPN should be automatically
reconnected within 5 minutes.

 * CenturyLink Gigabit internet with PPPoE dynamic IP.
 * Ubiquiti EdgeRouter-4 v2.0.8 firmware (mips64 GNU/Linux).
 * Google Cloud HA VPN - IKEv2 tunnel with BGP dynamic routing.
 * Cloudflare DNS A record needing to be updated as well.

IAM Roles
===

The service account used by Scarab needs the following roles assigned:

Compute Network Admin
---

`roles/compute.networkAdmin` is required to grant the
`compute.vpnTunnels.create` permission to create VPN Tunnels.
`compute.routers.use` is required to create the VPN Tunnel association with the
cloud router.

Compute Network User
---

`roles/compute.networkUser` is required to grant the
`compute.externalVpnGateways.list` permission to list VPN Gateways.

[1]: https://cloud.google.com/vpn/docs/concepts/overview
[2]: https://en.wikipedia.org/wiki/Synchronicity
[3]: https://cloud.google.com/compute/docs/reference/rest/v1/vpnTunnels
