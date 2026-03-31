---
page_title: "VirtualBox Extension Pack"
subcategory: "Guides"
description: |-
  Which features require the Oracle VirtualBox Extension Pack.
---

# VirtualBox Extension Pack

Some provider features require the [Oracle VirtualBox Extension Pack](https://www.virtualbox.org/wiki/Downloads) to be installed on the host.

## Features Requiring Extension Pack

| Feature | Attribute | Without Extension Pack | With Extension Pack |
|---------|-----------|----------------------|-------------------|
| USB 2.0 (EHCI) | `usb_controller = "ehci"` | Error | Works |
| USB 3.0 (xHCI) | `usb_controller = "xhci"` | Error | Works |
| Remote Display (VRDE) | Via `customize` | Error | Works |
| Disk Encryption | Via `customize` | Error | Works |
| PXE Boot (Intel NICs) | Via `customize` | Error | Works |

## Features That Work Without Extension Pack

All core functionality works without the Extension Pack:

- VM creation, modification, and deletion
- CPU, memory, firmware, OS type configuration
- NAT, bridged, host-only, and internal networking
- Port forwarding
- Shared folders
- Snapshots and linked clones
- Storage controllers (SATA, IDE, SCSI, NVMe)
- Disk management
- USB 1.1 (OHCI)
- Serial ports
- Clipboard and drag-and-drop
- Cloud-init user_data
- All data sources

## Installing the Extension Pack

```bash
# Download (match your VirtualBox version)
# Visit: https://www.virtualbox.org/wiki/Downloads

# Install via VBoxManage
VBoxManage extpack install Oracle_VirtualBox_Extension_Pack-<version>.vbox-extpack

# Verify
VBoxManage list extpacks
```

~> **Important:** The Extension Pack version must match your VirtualBox version exactly. A version mismatch can cause VM startup failures.

## Handling Missing Extension Pack in Terraform

If you use a feature that requires the Extension Pack and it's not installed, VBoxManage will return an error during `terraform apply`. The error message will indicate which feature requires the Extension Pack.

To avoid this, check the Extension Pack status using a data source:

```hcl
data "virtualbox_host_info" "current" {}

output "vbox_version" {
  value = data.virtualbox_host_info.current.virtualbox_version
}
```
