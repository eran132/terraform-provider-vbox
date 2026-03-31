# Changelog

## v1.1.0 (Unreleased)

### Features
- Add `terraform import` support for all resources (vm, disk, snapshot, hostonly_network, nat_network)
- Add linked clone support (`linked_clone` and `source_vm` attributes on virtualbox_vm)
- Add cloud-init user_data support on virtualbox_vm
- Improved box format handling for newer Vagrant boxes (Ubuntu 22.04+, Rocky 9, etc.)
- Recursive disk discovery for nested box layouts

### Bug Fixes
- Pin GoReleaser to v2.14.0 to avoid SHA256SUMS.sig bug in v2.15.1

## v1.0.0

### Features
- Complete rewrite from terra-farm/terraform-provider-virtualbox
- 5 resources: virtualbox_vm (30+ attributes), virtualbox_disk, virtualbox_snapshot, virtualbox_hostonly_network, virtualbox_nat_network
- 3 data sources: virtualbox_host_info, virtualbox_vm, virtualbox_network
- VBoxManage command layer (internal/vboxmanage/) with mockable Driver interface
- terraform-plugin-framework + mux for future resource migration
- NAT port forwarding, host-only networking, shared folders, snapshots
- OS type, firmware (EFI/BIOS), graphics controller, chipset selection
- USB controllers, serial ports, clipboard, drag-and-drop
- CPU execution cap, hardware virtualization flags
- VBoxManage customize escape hatch for arbitrary commands
- Acceptance tests verified against real VirtualBox
- Full Terraform Registry documentation

### Improvements over terra-farm/terraform-provider-virtualbox
- Upgraded from Go 1.20 to Go 1.25
- Upgraded from terraform-plugin-sdk v2.21.0 to latest + plugin-framework
- Fixed all typos and bugs in error messages
- Refactored 895-line monolith into focused modules
- Updated CI/CD (golangci-lint v7, goreleaser v2, GitHub Actions v4+)
- Added .gitignore, golangci-lint config, registry manifest
