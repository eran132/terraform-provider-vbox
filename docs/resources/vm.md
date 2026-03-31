---
page_title: "virtualbox_vm Resource - VirtualBox Provider"
subcategory: ""
description: |-
  Manages a VirtualBox virtual machine.
---

# virtualbox_vm (Resource)

Manages the full lifecycle of a VirtualBox virtual machine, including CPU, memory, networking, storage, and advanced hardware settings.

## Example Usage

```hcl
resource "virtualbox_vm" "web" {
  name   = "web-server"
  image  = "https://app.vagrantup.com/ubuntu/boxes/bionic64/versions/20180903.0.0/providers/virtualbox.box"
  cpus   = 2
  memory = "1024mib"
  status = "running"

  network_adapter {
    type   = "nat"
    device = "IntelPro1000MTDesktop"

    port_forwarding {
      name      = "ssh"
      protocol  = "tcp"
      host_ip   = "127.0.0.1"
      host_port = 2222
      guest_ip  = ""
      guest_port = 22
    }

    port_forwarding {
      name       = "http"
      protocol   = "tcp"
      host_ip    = "127.0.0.1"
      host_port  = 8080
      guest_ip   = ""
      guest_port = 80
    }
  }

  shared_folder {
    name      = "project"
    host_path = "/home/user/project"
    auto_mount = true
    writable   = true
  }
}
```

### Example: Import from OVA

```hcl
resource "virtualbox_vm" "imported" {
  name       = "imported-server"
  ova_source = "/path/to/appliance.ova"
  cpus       = 2
  memory     = "1024mib"
  status     = "running"

  network_adapter {
    type = "nat"
  }
}
```

## Argument Reference

### Required

- `name` - (Required, ForceNew) The name of the virtual machine.

### Image Source (one of the following)

- `image` - (Optional, ForceNew) URL or local path to a Vagrant `.box` file used to create the VM. Required unless `ova_source` or `linked_clone` is used.
- `ova_source` - (Optional, ForceNew) Path to an OVA/OVF file to import as this VM. Mutually exclusive with `image`. The OVA/OVF appliance is imported via `VBoxManage import` and the VM is then configured with the specified Terraform attributes.

### Optional - General

- `cpus` - (Optional) Number of virtual CPUs. Defaults to `2`.
- `memory` - (Optional) Amount of RAM (e.g. `"512mib"`, `"1gib"`). Defaults to `"512mib"`.
- `status` - (Optional) Desired power state: `"running"` or `"poweroff"`. Defaults to `"running"`.
- `os_type` - (Optional) Guest OS type identifier. Defaults to `"Linux_64"`.
- `gui` - (Optional) Whether to display the VM GUI window. Defaults to `false`.
- `user_data` - (Optional) Cloud-init user data (cloud-config YAML or script). Passed to the VM via VirtualBox guest properties.

### Optional - Display

- `vram` - (Optional) Video memory in MB. Defaults to `20`.
- `graphics_controller` - (Optional) Graphics controller type: `none`, `vboxvga`, `vmsvga`, or `vboxsvga`.

### Optional - Hardware

- `firmware` - (Optional) Firmware type: `bios`, `efi`, `efi32`, or `efi64`.
- `chipset` - (Optional) Chipset type: `piix3` or `ich9`.
- `cpu_execution_cap` - (Optional) CPU execution cap percentage, `1`-`100`. Defaults to `100`.
- `ioapic` - (Optional, Bool) Enable I/O APIC.
- `pae` - (Optional, Bool) Enable Physical Address Extension.
- `nested_hw_virt` - (Optional, Bool) Enable nested hardware virtualization.
- `largepages` - (Optional, Bool) Enable large pages.
- `vtx_vpid` - (Optional, Bool) Enable VT-x VPID.

### Optional - Interaction

- `clipboard_mode` - (Optional) Shared clipboard mode: `disabled`, `hosttoguest`, `guesttohost`, or `bidirectional`.
- `drag_and_drop` - (Optional) Drag-and-drop mode: `disabled`, `hosttoguest`, `guesttohost`, or `bidirectional`.

### Optional - USB

- `usb_controller` - (Optional) USB controller type: `ohci`, `ehci`, or `xhci`.

### Optional - Cloning

- `linked_clone` - (Optional, Bool) Create a linked clone instead of a full clone.
- `source_vm` - (Optional) Name or UUID of an existing VM to clone from.

### Optional - Boot & Media

- `optical_disks` - (Optional, List of String) Paths to ISO images to attach.
- `boot_order` - (Optional, List of String) Boot device order (e.g. `["disk", "dvd", "net"]`).

### Optional - Blocks

#### `network_adapter`

Configures a network adapter on the VM.

- `type` - (Required) Adapter type: `nat`, `bridged`, `hostonly`, `intnet`, `natnetwork`, or `generic`.
- `device` - (Optional) Adapter hardware model (e.g. `"IntelPro1000MTDesktop"`).
- `host_interface` - (Optional) Host interface name for `bridged` or `hostonly` types.
- `mac_address` - (Optional) MAC address. Auto-generated if omitted.
- `promiscuous_mode` - (Optional) Promiscuous mode: `deny`, `allow-vms`, or `allow-all`.
- `cable_connected` - (Optional, Bool) Whether the virtual cable is connected. Defaults to `true`.
- `nat_dns_host_resolver` - (Optional, Bool) Use the host's DNS resolver for NAT.
- `nat_dns_proxy` - (Optional, Bool) Use the NAT DNS proxy.

##### `port_forwarding` (nested inside `network_adapter`)

- `name` - (Required) Rule name.
- `protocol` - (Required) Protocol: `tcp` or `udp`.
- `host_ip` - (Optional) Host IP to bind.
- `host_port` - (Required) Host port number.
- `guest_ip` - (Optional) Guest IP.
- `guest_port` - (Required) Guest port number.

#### `storage_controller`

Configures a storage controller.

- `name` - (Required) Controller name.
- `type` - (Required) Controller type (e.g. `"ide"`, `"sata"`, `"scsi"`, `"sas"`, `"nvme"`).

#### `disk_attachment`

Attaches a disk to a storage controller.

- `controller` - (Required) Name of the storage controller.
- `port` - (Required) Port number.
- `device` - (Required) Device number.
- `disk_id` - (Required) ID of a `virtualbox_disk` resource.

#### `shared_folder`

Shares a host directory with the guest.

- `name` - (Required) Share name visible inside the guest.
- `host_path` - (Required) Absolute path on the host.
- `auto_mount` - (Optional, Bool) Auto-mount in the guest.
- `writable` - (Optional, Bool) Allow write access.

#### `serial_port`

Configures a serial port.

- `slot` - (Required) Serial port slot number (0-3).
- `mode` - (Required) Port mode (e.g. `"disconnected"`, `"hostpipe"`, `"file"`).
- `path` - (Optional) Path for the serial port output.

#### `customize`

Runs arbitrary `VBoxManage` commands after VM creation.

- Each entry is a list of strings representing a single `VBoxManage` invocation.

```hcl
customize = [
  ["modifyvm", "{{.Name}}", "--uartmode1", "file", "/tmp/console.log"],
]
```

## Attribute Reference

- `id` - The UUID of the virtual machine.
