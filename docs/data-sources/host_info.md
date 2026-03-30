---
page_title: "virtualbox_host_info Data Source - VirtualBox Provider"
subcategory: ""
description: |-
  Reads information about the VirtualBox host environment.
---

# virtualbox_host_info (Data Source)

Retrieves information about the VirtualBox installation and available host networking interfaces.

## Example Usage

```hcl
data "virtualbox_host_info" "current" {}

output "vbox_version" {
  value = data.virtualbox_host_info.current.virtualbox_version
}
```

## Attribute Reference

- `virtualbox_version` - The installed VirtualBox version string.
- `host_only_interfaces` - A list of host-only network interfaces. Each entry contains:
  - `name` - Interface name.
  - `ipv4_address` - IPv4 address.
  - `ipv4_netmask` - IPv4 netmask.
  - `status` - Interface status.
- `bridged_interfaces` - A list of bridged network interfaces. Each entry contains:
  - `name` - Interface name.
  - `status` - Interface status.
