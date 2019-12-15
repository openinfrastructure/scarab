Scarab
===

Scarab updates a [Cloud VPN][1] peer gateway with the IP address provided.
The intended use case is to re-establish a site-to-site VPN when a dynamic IP
address changes.

The name comes from the Golden Scarab in the [story about syncronicity][2].

Use case

 * CenturyLink Gigabit internet with PPPoE dynamic IP.
 * Ubiquiti EdgeRouter-4 v2.0.8 firmware (mips64 GNU/Linux)
 * Google Cloud VPN HA - IKEv2 tunnel with BGP dynamic routing.

Roadmap
===

 * [ ] Make `scarab version` read a common semantic version
 * [ ] Add `scarab create` to create a tunnel
 * [ ] Add `scarab delete` to create a tunnel
 * [ ] Add `scarab update` to update an existing tunnel
 * [ ] Design how to handle the tunnel preshared key

[1]: https://cloud.google.com/vpn/docs/concepts/overview
[2]: https://en.wikipedia.org/wiki/Synchronicity
