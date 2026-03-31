package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: vagrant2terraform <Vagrantfile> [output.tf]\n")
		fmt.Fprintf(os.Stderr, "\nConverts a Vagrantfile to Terraform HCL for the eran132/vbox provider.\n")
		os.Exit(1)
	}

	vagrantfile := os.Args[1]
	output := "main.tf"
	if len(os.Args) >= 3 {
		output = os.Args[2]
	}

	data, err := os.ReadFile(vagrantfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", vagrantfile, err)
		os.Exit(1)
	}

	hcl := convertVagrantfileToHCL(string(data))

	if output == "-" {
		fmt.Print(hcl)
	} else {
		if err := os.WriteFile(output, []byte(hcl), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", output, err)
			os.Exit(1)
		}
		fmt.Printf("Converted %s -> %s\n", vagrantfile, output)
	}
}

// convertVagrantfileToHCL parses a Vagrantfile and generates Terraform HCL.
// This handles the most common Vagrant patterns — it's not a full Ruby parser.
func convertVagrantfileToHCL(vagrantfile string) string {
	var sb strings.Builder

	// Header
	sb.WriteString(`terraform {
  required_providers {
    virtualbox = {
      source  = "eran132/vbox"
      version = "~> 2.0"
    }
  }
}

`)

	// Parse box name
	box := extractValue(vagrantfile, `config\.vm\.box\s*=\s*["']([^"']+)["']`)
	if box == "" {
		box = "ubuntu/focal64"
	}

	// Parse box URL or image
	boxURL := extractValue(vagrantfile, `config\.vm\.box_url\s*=\s*["']([^"']+)["']`)
	if boxURL == "" {
		// Construct Vagrant Cloud URL from box name
		parts := strings.SplitN(box, "/", 2)
		if len(parts) == 2 {
			boxURL = fmt.Sprintf("https://app.vagrantup.com/%s/boxes/%s/versions/latest/providers/virtualbox.box", parts[0], parts[1])
		}
	}

	// Parse VM name
	vmName := extractValue(vagrantfile, `config\.vm\.hostname\s*=\s*["']([^"']+)["']`)
	if vmName == "" {
		vmName = extractValue(vagrantfile, `vb\.name\s*=\s*["']([^"']+)["']`)
	}
	if vmName == "" {
		vmName = strings.ReplaceAll(box, "/", "-")
	}

	// Parse memory
	memory := extractValue(vagrantfile, `vb\.memory\s*=\s*["']?(\d+)["']?`)
	if memory == "" {
		memory = "1024"
	}

	// Parse CPUs
	cpus := extractValue(vagrantfile, `vb\.cpus\s*=\s*["']?(\d+)["']?`)
	if cpus == "" {
		cpus = "2"
	}

	// Parse GUI mode
	gui := extractValue(vagrantfile, `vb\.gui\s*=\s*(true|false)`)

	// Parse linked clone
	linkedClone := extractValue(vagrantfile, `vb\.linked_clone\s*=\s*(true|false)`)

	// Start VM resource
	resourceName := sanitizeResourceName(vmName)
	sb.WriteString(fmt.Sprintf("resource \"virtualbox_vm\" %q {\n", resourceName))
	sb.WriteString(fmt.Sprintf("  name   = %q\n", vmName))
	sb.WriteString(fmt.Sprintf("  image  = %q\n", boxURL))
	sb.WriteString(fmt.Sprintf("  cpus   = %s\n", cpus))
	sb.WriteString(fmt.Sprintf("  memory = %q\n", memory+"mib"))

	if gui == "true" {
		sb.WriteString("  gui    = true\n")
	}
	if linkedClone == "true" {
		sb.WriteString("  linked_clone = true\n")
	}

	// Parse forwarded ports
	portRegex := regexp.MustCompile(`config\.vm\.network\s*["']forwarded_port["'],\s*guest:\s*(\d+),\s*host:\s*(\d+)(?:,\s*protocol:\s*["'](\w+)["'])?(?:,\s*host_ip:\s*["']([^"']+)["'])?`)
	portMatches := portRegex.FindAllStringSubmatch(vagrantfile, -1)

	// Parse private networks
	privateNetRegex := regexp.MustCompile(`config\.vm\.network\s*["']private_network["'],\s*ip:\s*["']([^"']+)["']`)
	privateNetMatches := privateNetRegex.FindAllStringSubmatch(vagrantfile, -1)

	// Parse public networks
	publicNetRegex := regexp.MustCompile(`config\.vm\.network\s*["']public_network["'](?:,\s*bridge:\s*["']([^"']+)["'])?`)
	publicNetMatches := publicNetRegex.FindAllStringSubmatch(vagrantfile, -1)

	// Default NAT adapter (Vagrant always has one)
	sb.WriteString("\n  network_adapter {\n")
	sb.WriteString("    type = \"nat\"\n")

	if len(portMatches) > 0 {
		sb.WriteString("\n")
		for _, match := range portMatches {
			guestPort := match[1]
			hostPort := match[2]
			protocol := "tcp"
			if len(match) > 3 && match[3] != "" {
				protocol = match[3]
			}
			hostIP := ""
			if len(match) > 4 && match[4] != "" {
				hostIP = match[4]
			}

			// Generate a rule name
			ruleName := fmt.Sprintf("port_%s", guestPort)
			if guestPort == "22" {
				ruleName = "ssh"
			} else if guestPort == "80" {
				ruleName = "http"
			} else if guestPort == "443" {
				ruleName = "https"
			} else if guestPort == "3306" {
				ruleName = "mysql"
			} else if guestPort == "5432" {
				ruleName = "postgres"
			} else if guestPort == "8080" {
				ruleName = "http_alt"
			}

			sb.WriteString("    port_forwarding {\n")
			sb.WriteString(fmt.Sprintf("      name       = %q\n", ruleName))
			sb.WriteString(fmt.Sprintf("      protocol   = %q\n", protocol))
			if hostIP != "" {
				sb.WriteString(fmt.Sprintf("      host_ip    = %q\n", hostIP))
			}
			sb.WriteString(fmt.Sprintf("      host_port  = %s\n", hostPort))
			sb.WriteString(fmt.Sprintf("      guest_port = %s\n", guestPort))
			sb.WriteString("    }\n")
		}
	}
	sb.WriteString("  }\n")

	// Host-only adapters
	for _, match := range privateNetMatches {
		ip := match[1]
		sb.WriteString("\n  network_adapter {\n")
		sb.WriteString("    type           = \"hostonly\"\n")
		sb.WriteString(fmt.Sprintf("    # Static IP: %s (configure in guest)\n", ip))
		sb.WriteString("    host_interface = \"VirtualBox Host-Only Ethernet Adapter\"\n")
		sb.WriteString("  }\n")
	}

	// Bridged adapters
	for _, match := range publicNetMatches {
		bridge := ""
		if len(match) > 1 {
			bridge = match[1]
		}
		sb.WriteString("\n  network_adapter {\n")
		sb.WriteString("    type           = \"bridged\"\n")
		if bridge != "" {
			sb.WriteString(fmt.Sprintf("    host_interface = %q\n", bridge))
		} else {
			sb.WriteString("    host_interface = \"\" # TODO: set your bridge interface\n")
		}
		sb.WriteString("  }\n")
	}

	// Parse synced folders
	syncRegex := regexp.MustCompile(`config\.vm\.synced_folder\s*["']([^"']+)["'],\s*["']([^"']+)["'](?:,\s*disabled:\s*(true))?`)
	syncMatches := syncRegex.FindAllStringSubmatch(vagrantfile, -1)

	for _, match := range syncMatches {
		hostPath := match[1]
		// guestPath := match[2] // mount point
		disabled := len(match) > 3 && match[3] == "true"

		if disabled || hostPath == "." {
			continue // Skip disabled or default sync
		}

		folderName := sanitizeResourceName(strings.ReplaceAll(hostPath, "/", "_"))
		sb.WriteString("\n  shared_folder {\n")
		sb.WriteString(fmt.Sprintf("    name       = %q\n", folderName))
		sb.WriteString(fmt.Sprintf("    host_path  = %q\n", hostPath))
		sb.WriteString("    auto_mount = true\n")
		sb.WriteString("  }\n")
	}

	// Parse customize commands (VBoxManage)
	customizeRegex := regexp.MustCompile(`vb\.customize\s*\[([^\]]+)\]`)
	customizeMatches := customizeRegex.FindAllStringSubmatch(vagrantfile, -1)

	if len(customizeMatches) > 0 {
		sb.WriteString("\n  customize = [\n")
		for _, match := range customizeMatches {
			args := match[1]
			// Parse Ruby array into strings
			argRegex := regexp.MustCompile(`["']([^"']+)["']|:id`)
			argMatches := argRegex.FindAllStringSubmatch(args, -1)

			sb.WriteString("    [")
			for i, arg := range argMatches {
				if i > 0 {
					sb.WriteString(", ")
				}
				if strings.Contains(arg[0], ":id") {
					sb.WriteString(`":id"`)
				} else {
					sb.WriteString(fmt.Sprintf("%q", arg[1]))
				}
			}
			sb.WriteString("],\n")
		}
		sb.WriteString("  ]\n")
	}

	// Parse shell provisioner — convert to comment with remote-exec hint
	shellRegex := regexp.MustCompile(`config\.vm\.provision\s*["']shell["'],\s*inline:\s*<<-?(\w+)([\s\S]*?)\1`)
	shellMatches := shellRegex.FindAllStringSubmatch(vagrantfile, -1)

	if len(shellMatches) > 0 {
		sb.WriteString("\n  # TODO: Convert shell provisioners to Terraform provisioner blocks\n")
		sb.WriteString("  # Example:\n")
		sb.WriteString("  # connection {\n")
		sb.WriteString("  #   type        = \"ssh\"\n")
		sb.WriteString("  #   user        = \"vagrant\"\n")
		sb.WriteString("  #   private_key = file(\"~/.vagrant.d/insecure_private_key\")\n")
		sb.WriteString("  #   host        = \"127.0.0.1\"\n")
		sb.WriteString("  #   port        = 2222\n")
		sb.WriteString("  # }\n")
		sb.WriteString("  # provisioner \"remote-exec\" {\n")
		sb.WriteString("  #   inline = [\n")
		for _, match := range shellMatches {
			script := strings.TrimSpace(match[2])
			for _, line := range strings.Split(script, "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					sb.WriteString(fmt.Sprintf("  #     %q,\n", line))
				}
			}
		}
		sb.WriteString("  #   ]\n")
		sb.WriteString("  # }\n")
	}

	sb.WriteString("}\n")

	// Add outputs
	sb.WriteString(fmt.Sprintf("\noutput \"vm_id\" {\n  value = virtualbox_vm.%s.id\n}\n", resourceName))
	sb.WriteString(fmt.Sprintf("\noutput \"vm_status\" {\n  value = virtualbox_vm.%s.status\n}\n", resourceName))

	return sb.String()
}

// extractValue extracts the first match of a regex pattern.
func extractValue(text, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(text)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

// sanitizeResourceName converts a string to a valid Terraform resource name.
func sanitizeResourceName(name string) string {
	name = strings.ToLower(name)
	re := regexp.MustCompile(`[^a-z0-9_]`)
	name = re.ReplaceAllString(name, "_")
	// Remove leading numbers
	name = regexp.MustCompile(`^[0-9]+`).ReplaceAllString(name, "")
	if name == "" {
		name = "vm"
	}
	return name
}
