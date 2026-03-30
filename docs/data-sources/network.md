---
page_title: "virtualbox_network Data Source - VirtualBox Provider"
subcategory: ""
description: |-
  Lists VirtualBox network configurations.
---

# virtualbox_network (Data Source)

Retrieves all configured VirtualBox networks, including host-only, NAT, and bridged interfaces.

## Example Usage

```hcl
data "virtualbox_network" "all" {}

output "nat_networks" {
  value = data.virtualbox_network.all.nat_networks
}
```

## Attribute Reference

- `host_only_networks` - A list of host-only networks. Each entry contains:
  - `name` - Network interface name.
  - `ipv4_address` - IPv4 address.
  - `ipv4_netmask` - IPv4 netmask.
  - `dhcp_enabled` - Whether DHCP is enabled.
- `nat_networks` - A list of NAT networks. Each entry contains:
  - `name` - Network name.
  - `network` - Network CIDR.
  - `dhcp_enabled` - Whether DHCP is enabled.
  - `ipv6` - Whether IPv6 is enabled.
  - `enabled` - Whether the network is active.
- `bridged_interfaces` - A list of bridged interfaces. Each entry contains:
  - `name` - Interface name.
  - `status` - Interface status.
