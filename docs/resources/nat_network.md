---
page_title: "virtualbox_nat_network Resource - VirtualBox Provider"
subcategory: ""
description: |-
  Manages a VirtualBox NAT network.
---

# virtualbox_nat_network (Resource)

Creates and manages a VirtualBox NAT network with optional port forwarding rules.

## Example Usage

```hcl
resource "virtualbox_nat_network" "lab" {
  name         = "lab-network"
  network      = "10.0.2.0/24"
  dhcp_enabled = true
  ipv6         = false
  enabled      = true

  port_forwarding {
    name       = "ssh"
    protocol   = "tcp"
    host_ip    = ""
    host_port  = 2222
    guest_ip   = "10.0.2.15"
    guest_port = 22
  }
}
```

## Argument Reference

- `name` - (Required) Name of the NAT network.
- `network` - (Required) Network CIDR (e.g. `"10.0.2.0/24"`).
- `dhcp_enabled` - (Optional, Bool) Enable DHCP on the NAT network. Defaults to `true`.
- `ipv6` - (Optional, Bool) Enable IPv6 support. Defaults to `false`.
- `enabled` - (Optional, Bool) Whether the network is enabled. Defaults to `true`.

### `port_forwarding`

- `name` - (Required) Rule name.
- `protocol` - (Required) Protocol: `tcp` or `udp`.
- `host_ip` - (Optional) Host IP to bind.
- `host_port` - (Required) Host port number.
- `guest_ip` - (Optional) Guest IP address.
- `guest_port` - (Required) Guest port number.
