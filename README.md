# Packer Plugin Ansible AAP

This plugin provides integration between HashiCorp Packer and Red Hat Ansible Automation Platform (AAP). It enables automatic registration of newly provisioned infrastructure with AAP, including inventory creation, host registration, and job template execution.

## Features

- **Create Inventory**: Automatically create inventories in Ansible Automation Platform
- **Add Host to Inventory**: Register newly created hosts with AAP inventories
- **Create Host**: Create standalone hosts in AAP without inventory association
- **Run Job Template**: Trigger job templates immediately after provisioning

## Components

This plugin contains:
- An AAP post-processor ([post-processor/aap](post-processor/aap)) - Main component for AAP integration
- A builder ([builder/scaffolding](builder/scaffolding))
- A provisioner ([provisioner/scaffolding](provisioner/scaffolding))
- A post-processor ([post-processor/scaffolding](post-processor/scaffolding))
- A data source ([datasource/scaffolding](datasource/scaffolding))
- Docs ([docs](docs))
- Working examples ([example](example))

## Usage

Add the plugin to your Packer template:

```hcl
packer {
  required_plugins {
    ansible-aap = {
      version = ">=v0.1.0"
      source  = "github.com/danielino/ansible-aap"
    }
  }
}

build {
  sources = ["source.amazon-ebs.example"]

  post-processor "aap" {
    aap_url      = "https://aap.example.com"
    aap_username = "admin"
    aap_password = "password"

    create_inventory     = true
    inventory_name       = "packer-builds"
    organization_id      = 1

    create_host          = true
    add_host_to_inventory = true
  }
}
```

See [example/aap.pkr.hcl](example/aap.pkr.hcl) for a complete example.

## Build from source

1. Clone this GitHub repository locally.

2. Run this command from the root directory: 
```shell 
go build -ldflags="-X github.com/danielino/packer-plugin-ansible-aap/version.VersionPrerelease=dev" -o packer-plugin-ansible-aap
```

3. After you successfully compile, the `packer-plugin-ansible-aap` plugin binary file is in the root directory. 

4. To install the compiled plugin, run the following command 
```shell
packer plugins install --path packer-plugin-ansible-aap github.com/danielino/ansible-aap
```

### Build on *nix systems
Unix like systems with the make, sed, and grep commands installed can use the `make dev` to execute the build from source steps. 

### Build on Windows Powershell
The preferred solution for building on Windows are steps 2-4 listed above.
If you would prefer to script the building process you can use the following as a guide

```powershell
$MODULE_NAME = (Get-Content go.mod | Where-Object { $_ -match "^module"  }) -replace 'module ',''
$FQN = $MODULE_NAME -replace 'packer-plugin-',''
go build -ldflags="-X $MODULE_NAME/version.VersionPrerelease=dev" -o packer-plugin-scaffolding.exe
packer plugins install --path packer-plugin-scaffolding.exe $FQN
```

## Running Acceptance Tests

Make sure to install the plugin locally using the steps in [Build from source](#build-from-source).

Once everything needed is set up, run:
```
PACKER_ACC=1 go test -count 1 -v ./... -timeout=120m
```

This will run the acceptance tests for all plugins in this set.

## Registering Plugin as Packer Integration

Partner and community plugins can be hard to find if a user doesn't know what 
they are looking for. To assist with plugin discovery Packer offers an integration
portal at https://developer.hashicorp.com/packer/integrations to list known integrations 
that work with the latest release of Packer. 

Registering a plugin as an integration requires [metadata configuration](./metadata.hcl) within the plugin
repository and approval by the Packer team. To initiate the process of registering your 
plugin as a Packer integration refer to the [Developing Plugins](https://developer.hashicorp.com/packer/docs/plugins/creation#registering-plugins) page.

# Requirements

-	[packer-plugin-sdk](https://github.com/hashicorp/packer-plugin-sdk) >= v0.5.2
-	[Go](https://golang.org/doc/install) >= 1.20

## Packer Compatibility
This scaffolding template is compatible with Packer >= v1.10.2
