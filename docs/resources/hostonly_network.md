---
page_title: "virtualbox_hostonly_network Resource - VirtualBox Provider"
subcategory: ""
description: |-
  Manages a VirtualBox host-only network interface.
---

# virtualbox_hostonly_network (Resource)

Creates and manages a VirtualBox host-only network interface with optional DHCP server configuration.

## Example Usage

```hcl
resource "virtualbox_hostonly_network" "dev" {
  ipv4_address = "192.168.56.1"
  ipv4_netmask = "255.255.255.0"
  dhcp_enabled = true
  dhcp_lower_ip = "192.168.56.100"
  dhcp_upper_ip = "192.168.56.200"
}
```

## Argument Reference

- `ipv4_address` - (Optional) IPv4 address for the host-only adapter.
- `ipv4_netmask` - (Optional) IPv4 netmask. Defaults to `"255.255.255.0"`.
- `ipv6_address` - (Optional) IPv6 address for the host-only adapter.
- `ipv6_prefix` - (Optional) IPv6 prefix length.
- `dhcp_enabled` - (Optional, Bool) Enable the DHCP server on this network.
- `dhcp_lower_ip` - (Optional) Lower bound of the DHCP address range.
- `dhcp_upper_ip` - (Optional) Upper bound of the DHCP address range.

## Attribute Reference

- `name` - (Computed) The interface name assigned by VirtualBox (e.g. `vboxnet0`).
