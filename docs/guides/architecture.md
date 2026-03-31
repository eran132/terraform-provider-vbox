---
page_title: "Provider Architecture"
subcategory: "Guides"
description: |-
  How the VirtualBox provider communicates with VirtualBox.
---

# Provider Architecture

## How It Works

The provider communicates with Oracle VirtualBox through the `VBoxManage` command-line interface. This is the same approach used by Vagrant and Packer — it's the most portable and well-tested method for VirtualBox automation.

```
┌─────────────┐     ┌──────────────────┐     ┌──────────────┐     ┌────────────┐
│  Terraform  │────▶│  eran132/vbox    │────▶│  VBoxManage  │────▶│ VirtualBox │
│    CLI      │     │  (Go provider)   │     │   (CLI)      │     │  Engine    │
└─────────────┘     └──────────────────┘     └──────────────┘     └────────────┘
```

## VBoxManage Layer

The provider includes a custom VBoxManage command layer (`internal/vboxmanage/`) with a mockable `Driver` interface for testing. This layer wraps common operations:

- **VM lifecycle**: createvm, startvm, controlvm, unregistervm, clonevm
- **Configuration**: modifyvm, showvminfo
- **Storage**: storagectl, storageattach, createmedium, modifymedium
- **Networking**: hostonlyif, natnetwork, dhcpserver
- **Snapshots**: snapshot take/restore/delete/list
- **Guest interaction**: guestproperty get/set

## Why Not COM/XPCOM?

VirtualBox offers COM (Windows) and XPCOM (Linux/macOS) APIs for direct programmatic access. We chose VBoxManage instead because:

1. **Cross-platform**: VBoxManage works identically on all platforms. COM/XPCOM requires platform-specific C bindings via CGO, which prevents Go cross-compilation.
2. **Version-independent**: VBoxManage maintains backward compatibility across VirtualBox versions. COM/XPCOM APIs change between versions.
3. **Debuggable**: Every provider operation can be replicated by running the equivalent `VBoxManage` command manually.
4. **Proven**: Vagrant (Ruby) and Packer (Go) both use VBoxManage. This approach has been validated across millions of deployments over a decade.

## Plugin Framework

The provider uses both Terraform Plugin SDK v2 and Plugin Framework v1.19 via `terraform-plugin-mux`. This allows incremental migration of resources from the legacy SDK to the modern Framework while maintaining backward compatibility.

## Requirements

- [Oracle VirtualBox](https://www.virtualbox.org/wiki/Downloads) must be installed
- The `VBoxManage` binary must be on the system `PATH` (or installed in a standard location — the provider auto-detects common install paths)
- No VirtualBox Extension Pack required for core features (see [Extension Pack guide](extension-pack.md))
