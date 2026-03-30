---
page_title: "virtualbox_vm Data Source - VirtualBox Provider"
subcategory: ""
description: |-
  Reads information about an existing VirtualBox VM.
---

# virtualbox_vm (Data Source)

Looks up an existing VirtualBox virtual machine by name or UUID and exposes its configuration as read-only attributes.

## Example Usage

```hcl
data "virtualbox_vm" "existing" {
  name = "my-existing-vm"
}

output "vm_status" {
  value = data.virtualbox_vm.existing.status
}
```

## Argument Reference

- `name` - (Optional) Name of the VM to look up. Exactly one of `name` or `uuid` must be specified.
- `uuid` - (Optional) UUID of the VM to look up. Exactly one of `name` or `uuid` must be specified.

## Attribute Reference

- `uuid` - The UUID of the virtual machine.
- `name` - The name of the virtual machine.
- `status` - Current power state (e.g. `"running"`, `"poweroff"`).
- `cpus` - Number of virtual CPUs.
- `memory` - Memory allocation string.
- `os_type` - Guest OS type identifier.
- `vram` - Video memory in MB.
