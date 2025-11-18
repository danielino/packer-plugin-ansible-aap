// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aap

import (
	"testing"
)

func TestProvisioner_Prepare(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid minimal config",
			config: map[string]interface{}{
				"aap_url":      "https://aap.example.com",
				"aap_username": "admin",
				"aap_password": "password",
			},
			wantErr: false,
		},
		{
			name: "missing aap_url",
			config: map[string]interface{}{
				"aap_username": "admin",
				"aap_password": "password",
			},
			wantErr: true,
			errMsg:  "aap_url is required",
		},
		{
			name: "missing aap_username",
			config: map[string]interface{}{
				"aap_url":      "https://aap.example.com",
				"aap_password": "password",
			},
			wantErr: true,
			errMsg:  "aap_username is required",
		},
		{
			name: "missing aap_password",
			config: map[string]interface{}{
				"aap_url":      "https://aap.example.com",
				"aap_username": "admin",
			},
			wantErr: true,
			errMsg:  "aap_password is required",
		},
		{
			name: "create inventory without name",
			config: map[string]interface{}{
				"aap_url":          "https://aap.example.com",
				"aap_username":     "admin",
				"aap_password":     "password",
				"create_inventory": true,
			},
			wantErr: true,
			errMsg:  "inventory_name is required when create_inventory is true",
		},
		{
			name: "create inventory without organization",
			config: map[string]interface{}{
				"aap_url":          "https://aap.example.com",
				"aap_username":     "admin",
				"aap_password":     "password",
				"create_inventory": true,
				"inventory_name":   "test-inventory",
			},
			wantErr: true,
			errMsg:  "organization_id is required when create_inventory is true",
		},
		{
			name: "valid create inventory config",
			config: map[string]interface{}{
				"aap_url":          "https://aap.example.com",
				"aap_username":     "admin",
				"aap_password":     "password",
				"create_inventory": true,
				"inventory_name":   "test-inventory",
				"organization_id":  1,
			},
			wantErr: false,
		},
		{
			name: "add host without inventory",
			config: map[string]interface{}{
				"aap_url":               "https://aap.example.com",
				"aap_username":          "admin",
				"aap_password":          "password",
				"add_host_to_inventory": true,
			},
			wantErr: true,
			errMsg:  "either create_inventory must be true or existing_inventory_id must be set when add_host_to_inventory is true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provisioner{}
			err := p.Prepare(tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Prepare() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Prepare() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Prepare() unexpected error = %v", err)
				}
			}
		})
	}
}
