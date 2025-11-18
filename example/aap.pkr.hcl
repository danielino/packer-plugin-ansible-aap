# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Example Packer template for using the AAP provisioner
# This example demonstrates the AAP provisioner capabilities:
# 1. Creating an inventory in Ansible Automation Platform
# 2. Adding a host to the inventory
# 3. Running a job template
#
# The provisioner runs BEFORE the VM is powered off, during the provisioning phase

packer {
  required_plugins {
    ansible-aap = {
      version = ">=v0.1.0"
      source  = "github.com/danielino/ansible-aap"
    }
  }
}

# Variables for AAP connection
variable "aap_url" {
  type        = string
  description = "URL of your Ansible Automation Platform"
  default     = "https://aap.example.com"
}

variable "aap_username" {
  type        = string
  description = "AAP username"
  default     = "admin"
}

variable "aap_password" {
  type        = string
  description = "AAP password"
  sensitive   = true
}

variable "organization_id" {
  type        = number
  description = "AAP Organization ID"
  default     = 1
}

# Example using a null builder for demonstration
# In practice, you would use a real builder like amazon-ebs, virtualbox-iso, etc.
source "null" "example" {
  communicator = "none"
}

build {
  sources = ["source.null.example"]

  # Example 1: Create inventory and add host during provisioning
  # This runs BEFORE the VM is powered off
  provisioner "aap" {
    aap_url      = var.aap_url
    aap_username = var.aap_username
    aap_password = var.aap_password

    # Create a new inventory
    create_inventory      = true
    inventory_name        = "packer-example-inventory"
    inventory_description = "Example inventory created by Packer"
    organization_id       = var.organization_id

    # Create a host and add it to the inventory
    create_host           = true
    add_host_to_inventory = true
    host_name             = "example-host-${timestamp()}"
    host_description      = "Example host created by Packer provisioner"
    host_variables = {
      environment  = "development"
      ansible_host = "10.0.1.100"
      ansible_port = "22"
      ansible_user = "ubuntu"
    }

    # Run a job template
    run_job_template = true
    job_template_id  = 5
    job_extra_vars = {
      setup_type  = "initial"
      app_version = "1.0.0"
    }
  }
}

