// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/danielino/packer-plugin-ansible-aap/builder/scaffolding"
	scaffoldingData "github.com/danielino/packer-plugin-ansible-aap/datasource/scaffolding"
	aapPP "github.com/danielino/packer-plugin-ansible-aap/post-processor/aap"
	scaffoldingPP "github.com/danielino/packer-plugin-ansible-aap/post-processor/scaffolding"
	aapProv "github.com/danielino/packer-plugin-ansible-aap/provisioner/aap"
	scaffoldingProv "github.com/danielino/packer-plugin-ansible-aap/provisioner/scaffolding"
	scaffoldingVersion "github.com/danielino/packer-plugin-ansible-aap/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("my-builder", new(scaffolding.Builder))
	pps.RegisterProvisioner("my-provisioner", new(scaffoldingProv.Provisioner))
	pps.RegisterProvisioner("aap", new(aapProv.Provisioner))
	pps.RegisterPostProcessor("my-post-processor", new(scaffoldingPP.PostProcessor))
	pps.RegisterPostProcessor("aap", new(aapPP.PostProcessor))
	pps.RegisterDatasource("my-datasource", new(scaffoldingData.Datasource))
	pps.SetVersion(scaffoldingVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
