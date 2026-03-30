---
page_title: "virtualbox_snapshot Resource - VirtualBox Provider"
subcategory: ""
description: |-
  Manages a VirtualBox VM snapshot.
---

# virtualbox_snapshot (Resource)

Creates and manages a snapshot of a VirtualBox virtual machine.

## Example Usage

```hcl
resource "virtualbox_snapshot" "baseline" {
  vm_id       = virtualbox_vm.web.id
  name        = "baseline"
  description = "Clean state before deployment"
}
```

## Argument Reference

- `vm_id` - (Required) UUID of the virtual machine to snapshot.
- `name` - (Required) Name of the snapshot.
- `description` - (Optional) Description of the snapshot.
- `live` - (Optional, Bool) Take a live snapshot while the VM is running. Defaults to `false`.

## Attribute Reference

- `uuid` - (Computed) The UUID assigned to the snapshot by VirtualBox.
