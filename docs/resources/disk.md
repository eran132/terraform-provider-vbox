---
page_title: "virtualbox_disk Resource - VirtualBox Provider"
subcategory: ""
description: |-
  Manages a VirtualBox virtual disk image.
---

# virtualbox_disk (Resource)

Creates and manages a VirtualBox virtual disk image (VDI, VMDK, VHD).

## Example Usage

```hcl
resource "virtualbox_disk" "data" {
  file_path = "/tmp/vbox-disks/data.vdi"
  size      = 20480
  format    = "VDI"
}
```

## Argument Reference

- `file_path` - (Required) Absolute path where the disk image file will be created.
- `size` - (Required) Disk size in megabytes.
- `format` - (Optional) Disk format: `VDI`, `VMDK`, or `VHD`. Defaults to `VDI`.
- `variant` - (Optional) Disk variant: `Standard` (dynamically allocated) or `Fixed`. Defaults to `Standard`.

## Attribute Reference

- `uuid` - (Computed) The UUID assigned to the disk by VirtualBox.
