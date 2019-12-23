Plan
===

NEXT: Pick back up on api.go:84, once the external gateway is created the
tunnel should be able to be created.

 * [x] Make `scarab version` read a common semantic version
 * [x] Add `scarab list` to list vpn tunnels
 * [x] Scarab logs via standard output.
 * [x] Read configuration from /config/scarab/scarab.yaml
 * [x] Install scarab into EdgeRouter-4 /config
 * [ ] Update creates a new ExternalVPNGateway resource.
 * [ ] Update waits for the ExternalVPNGateway insert via waitDone().
 * [ ] waitDone() checks the status of each Operation.
 * [ ] waitDone() logs to the user.
 * [ ] pollOpenItem() sends forever until aborted.
 * [ ] go routine pollOpenItem() is verified to not leak.
 * [ ] Error fixed: Error 400: Invalid value for field, Either
     peerExternalGateway or peerGcpGateway must be specified for HA VPN tunnel
     creation., invalid
 * [ ] Set the ER-4 VPN local-address to the new IP.  `set vpn ipsec
     site-to-site peer 35.242.52.86 local-address $PPP_LOCAL`
 * [ ] Validate vpn tunnel is created correctly when run manually.
 * [ ] Execute `scarab update --address=${PPP_LOCAL}` from
     `/config/scripts/ppp/ip-up.d/scarab`
 * [ ] Update CloudFlare A record with new PPP_LOCAL address.

Optimizations
===

 * [ ] Modify waitDone() to use [globalOperations.wait][wait]
 * [ ] Log via syslog to `/var/log/messages` on the EdgeRouter-4
 * [ ] Log via Stackdriver Structured logs
 * [ ] Investigate and decide on async logging design

Development
===

Prefix matching is used to create new resources then delete the old resources.
`gcloud` may be used to quickly clean up during development.

Delete stale resources:

```bash
gcloud compute external-vpn-gateways list --format='get(NAME)' --filter='name:tun1*' \
  | xargs gcloud compute external-vpn-gateways delete
```

Call from pppd
===

This program is intended to update the VPN tunnel when pppd obtains a new IP
address.  dhclient hooks should also work but are untested.

On the ER-4 the process context is:

```
/usr/sbin/pppd call pppoe0
 \_ /bin/sh /etc/ppp/ip-up pppoe0 eth0.201 0 97.120.177.238 207.225.84.48 pppoe0
     \_ run-parts --regex ^[a-zA-Z0-9._-]+$ /config/scripts/ppp/ip-up.d --arg=pppoe0 --arg=eth0.201 --arg=0 --arg=97.120.177.238 --arg=207.225.84.48 --arg=pppoe0
         \_ /bin/bash /config/scripts/ppp/ip-up.d/scarab pppoe0 eth0.201 0 97.120.177.238 207.225.84.48 pppoe0
```

Environment variables exported to the ip-up.d script are:

```bash
BYTES_RCVD=0
BYTES_SENT=0
CALL_FILE=pppoe0
CONNECT_TIME=0
DEVICE=eth0.201
DNS1=205.171.3.25
DNS2=205.171.2.25
IFNAME=pppoe0
IPLOCAL=97.120.177.238
IPREMOTE=207.225.84.48
MACREMOTE=88:E0:F3:CE:A6:73
ORIG_UID=0
PATH=/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin
PPPD_PID=3838
PPPLOGNAME=root
PPP_IFACE=pppoe0
PPP_IPPARAM=pppoe0
PPP_LOCAL=97.120.177.238
PPP_REMOTE=207.225.84.48
PPP_SPEED=0
PPP_TTY=eth0.201
PPP_TTYNAME=eth0.201
PWD=/
SHLVL=2
USEPEERDNS=1
_=/usr/bin/env
```

Lessons learned
===

Channel senders close the channel
---

Channel senders are responsible for closing channels, not receivers.  This
avoids panics sending to closed channels.

Closed channels return nil value
---

Closed channels return the nil value for the channel type when selected.  This
is useful for abort channels, receivers communicate to senders by closing an
abort channel.

[wait]: https://cloud.google.com/compute/docs/reference/rest/v1/globalOperations/wait
