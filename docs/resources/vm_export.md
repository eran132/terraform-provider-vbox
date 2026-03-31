---
page_title: "virtualbox_vm_export Resource - VirtualBox Provider"
subcategory: ""
description: |-
  Exports a VirtualBox virtual machine to OVA or OVF format.
---

# virtualbox_vm_export (Resource)

Exports an existing VirtualBox virtual machine to an OVA or OVF appliance file. The exported file is managed as a Terraform resource and will be deleted when the resource is destroyed.

## Example Usage

```hcl
resource "virtualbox_vm" "web" {
  name   = "web-server"
  image  = "https://app.vagrantup.com/ubuntu/boxes/bionic64/versions/20180903.0.0/providers/virtualbox.box"
  cpus   = 2
  memory = "1024mib"
}

resource "virtualbox_vm_export" "backup" {
  vm_id       = virtualbox_vm.web.id
  output_path = "/tmp/web-server-backup.ova"
  manifest    = true
}
```

### Export as OVF

```hcl
resource "virtualbox_vm_export" "ovf_export" {
  vm_id       = virtualbox_vm.web.id
  output_path = "/tmp/web-server.ovf"
  format      = "ovf"
  manifest    = true
}
```

## Argument Reference

### Required

- `vm_id` - (Required, ForceNew) UUID or name of the VM to export.
- `output_path` - (Required, ForceNew) Output file path for the exported appliance (`.ova` or `.ovf`).

### Optional

- `format` - (Optional, ForceNew) Export format: `ova` (default) or `ovf`. When set to `ovf`, the `--ovf20` flag is passed to VBoxManage.
- `manifest` - (Optional, ForceNew) Whether to include a manifest file. Defaults to `true`.
- `options` - (Optional, ForceNew) List of additional option strings passed to `VBoxManage export`.

## Attribute Reference

- `id` - The output file path (same as `output_path`).
- `file_path` - The actual output file path after export.

## Import

VM exports can be imported by file path:

```shell
terraform import virtualbox_vm_export.backup /tmp/web-server-backup.ova
```
