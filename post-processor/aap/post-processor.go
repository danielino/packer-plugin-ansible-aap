// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package aap

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/danielino/packer-plugin-ansible-aap/aap"
	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// AAP Connection settings
	AAPUrl      string `mapstructure:"aap_url" required:"true"`
	AAPUsername string `mapstructure:"aap_username" required:"true"`
	AAPPassword string `mapstructure:"aap_password" required:"true"`

	// Inventory settings
	CreateInventory    bool   `mapstructure:"create_inventory"`
	InventoryName      string `mapstructure:"inventory_name"`
	InventoryDesc      string `mapstructure:"inventory_description"`
	OrganizationID     int    `mapstructure:"organization_id"`
	ExistingInventory  int    `mapstructure:"existing_inventory_id"`

	// Host settings
	CreateHost      bool              `mapstructure:"create_host"`
	HostName        string            `mapstructure:"host_name"`
	HostDescription string            `mapstructure:"host_description"`
	HostVariables   map[string]string `mapstructure:"host_variables"`
	AddToInventory  bool              `mapstructure:"add_host_to_inventory"`

	// Job template settings
	RunJobTemplate bool              `mapstructure:"run_job_template"`
	JobTemplateID  int               `mapstructure:"job_template_id"`
	JobLimit       string            `mapstructure:"job_limit"`
	JobExtraVars   map[string]string `mapstructure:"job_extra_vars"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.post-processor.aap",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	// Validation
	if p.config.AAPUrl == "" {
		return fmt.Errorf("aap_url is required")
	}
	if p.config.AAPUsername == "" {
		return fmt.Errorf("aap_username is required")
	}
	if p.config.AAPPassword == "" {
		return fmt.Errorf("aap_password is required")
	}

	if p.config.CreateInventory {
		if p.config.InventoryName == "" {
			return fmt.Errorf("inventory_name is required when create_inventory is true")
		}
		if p.config.OrganizationID == 0 {
			return fmt.Errorf("organization_id is required when create_inventory is true")
		}
	}

	if p.config.AddToInventory {
		if !p.config.CreateInventory && p.config.ExistingInventory == 0 {
			return fmt.Errorf("either create_inventory must be true or existing_inventory_id must be set when add_host_to_inventory is true")
		}
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	client := aap.NewClient(p.config.AAPUrl, p.config.AAPUsername, p.config.AAPPassword)

	var inventoryID int

	// Step 1: Create inventory if requested
	if p.config.CreateInventory {
		ui.Say(fmt.Sprintf("Creating inventory '%s' in AAP...", p.config.InventoryName))
		inventory, err := client.CreateInventory(p.config.InventoryName, p.config.InventoryDesc, p.config.OrganizationID)
		if err != nil {
			return nil, false, false, fmt.Errorf("failed to create inventory: %w", err)
		}
		inventoryID = inventory.ID
		ui.Say(fmt.Sprintf("Successfully created inventory with ID: %d", inventoryID))
	} else if p.config.ExistingInventory > 0 {
		inventoryID = p.config.ExistingInventory
	}

	// Step 2: Create or add host
	if p.config.CreateHost {
		hostname := p.config.HostName
		if hostname == "" {
			// Use artifact ID as hostname if not specified
			hostname = source.Id()
		}

		// Convert string map to interface map for variables
		variables := make(map[string]interface{})
		for k, v := range p.config.HostVariables {
			variables[k] = v
		}
		// Add artifact information to variables
		variables["packer_artifact_id"] = source.Id()
		variables["packer_builder_type"] = source.BuilderId()

		if p.config.AddToInventory && inventoryID > 0 {
			ui.Say(fmt.Sprintf("Adding host '%s' to inventory ID %d...", hostname, inventoryID))
			host, err := client.AddHost(inventoryID, hostname, p.config.HostDescription, variables)
			if err != nil {
				return nil, false, false, fmt.Errorf("failed to add host to inventory: %w", err)
			}
			ui.Say(fmt.Sprintf("Successfully added host with ID: %d", host.ID))
		} else {
			ui.Say(fmt.Sprintf("Creating host '%s' in AAP...", hostname))
			host, err := client.CreateHost(hostname, p.config.HostDescription, variables)
			if err != nil {
				return nil, false, false, fmt.Errorf("failed to create host: %w", err)
			}
			ui.Say(fmt.Sprintf("Successfully created host with ID: %d", host.ID))
		}
	}

	// Step 3: Run job template if requested
	if p.config.RunJobTemplate {
		if p.config.JobTemplateID == 0 {
			return nil, false, false, fmt.Errorf("job_template_id is required when run_job_template is true")
		}

		ui.Say(fmt.Sprintf("Running job template ID %d...", p.config.JobTemplateID))

		// Convert string map to interface map for extra vars
		extraVars := make(map[string]interface{})
		for k, v := range p.config.JobExtraVars {
			extraVars[k] = v
		}
		// Add artifact information to extra vars
		extraVars["packer_artifact_id"] = source.Id()

		launch := &aap.JobLaunch{
			ExtraVars: extraVars,
		}

		if inventoryID > 0 {
			launch.Inventory = inventoryID
		}
		if p.config.JobLimit != "" {
			launch.Limit = p.config.JobLimit
		}

		result, err := client.RunJobTemplate(p.config.JobTemplateID, launch)
		if err != nil {
			return nil, false, false, fmt.Errorf("failed to run job template: %w", err)
		}
		ui.Say(fmt.Sprintf("Successfully launched job with ID: %d, status: %s", result.ID, result.Status))
	}

	// Return the original artifact unchanged
	return source, true, true, nil
}
