// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateInventory(t *testing.T) {
	// Create a mock AAP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify path
		if r.URL.Path != "/api/v2/inventories/" {
			t.Errorf("Expected /api/v2/inventories/ path, got %s", r.URL.Path)
		}

		// Verify authentication
		username, password, ok := r.BasicAuth()
		if !ok || username != "testuser" || password != "testpass" {
			t.Error("Basic auth failed")
		}

		// Decode request body
		var inventory Inventory
		if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify request data
		if inventory.Name != "test-inventory" {
			t.Errorf("Expected name 'test-inventory', got %s", inventory.Name)
		}
		if inventory.Organization != 1 {
			t.Errorf("Expected organization 1, got %d", inventory.Organization)
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		response := Inventory{
			ID:           123,
			Name:         inventory.Name,
			Description:  inventory.Description,
			Organization: inventory.Organization,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL, "testuser", "testpass")

	// Test CreateInventory
	inventory, err := client.CreateInventory("test-inventory", "Test Description", 1)
	if err != nil {
		t.Fatalf("CreateInventory failed: %v", err)
	}

	if inventory.ID != 123 {
		t.Errorf("Expected ID 123, got %d", inventory.ID)
	}
	if inventory.Name != "test-inventory" {
		t.Errorf("Expected name 'test-inventory', got %s", inventory.Name)
	}
}

func TestCreateHost(t *testing.T) {
	// Create a mock AAP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v2/hosts/" {
			t.Errorf("Expected /api/v2/hosts/ path, got %s", r.URL.Path)
		}

		// Decode request body
		var host Host
		if err := json.NewDecoder(r.Body).Decode(&host); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify request data
		if host.Name != "test-host" {
			t.Errorf("Expected name 'test-host', got %s", host.Name)
		}
		if host.Inventory != 0 {
			t.Errorf("Expected inventory to be 0 (not set), got %d", host.Inventory)
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		response := Host{
			ID:          456,
			Name:        host.Name,
			Description: host.Description,
			Variables:   host.Variables,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL, "testuser", "testpass")

	// Test CreateHost
	variables := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	host, err := client.CreateHost("test-host", "Test Host Description", variables)
	if err != nil {
		t.Fatalf("CreateHost failed: %v", err)
	}

	if host.ID != 456 {
		t.Errorf("Expected ID 456, got %d", host.ID)
	}
	if host.Name != "test-host" {
		t.Errorf("Expected name 'test-host', got %s", host.Name)
	}
}

func TestAddHost(t *testing.T) {
	// Create a mock AAP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v2/hosts/" {
			t.Errorf("Expected /api/v2/hosts/ path, got %s", r.URL.Path)
		}

		// Decode request body
		var host Host
		if err := json.NewDecoder(r.Body).Decode(&host); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify request data
		if host.Name != "test-host" {
			t.Errorf("Expected name 'test-host', got %s", host.Name)
		}
		if host.Inventory != 123 {
			t.Errorf("Expected inventory 123, got %d", host.Inventory)
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		response := Host{
			ID:          789,
			Name:        host.Name,
			Description: host.Description,
			Inventory:   host.Inventory,
			Variables:   host.Variables,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL, "testuser", "testpass")

	// Test AddHost
	variables := map[string]interface{}{
		"ansible_host": "10.0.0.1",
	}
	host, err := client.AddHost(123, "test-host", "Test Host", variables)
	if err != nil {
		t.Fatalf("AddHost failed: %v", err)
	}

	if host.ID != 789 {
		t.Errorf("Expected ID 789, got %d", host.ID)
	}
	if host.Inventory != 123 {
		t.Errorf("Expected inventory 123, got %d", host.Inventory)
	}
}

func TestRunJobTemplate(t *testing.T) {
	// Create a mock AAP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedPath := "/api/v2/job_templates/42/launch/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected %s path, got %s", expectedPath, r.URL.Path)
		}

		// Decode request body
		var launch JobLaunch
		if err := json.NewDecoder(r.Body).Decode(&launch); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify request data
		if launch.Inventory != 123 {
			t.Errorf("Expected inventory 123, got %d", launch.Inventory)
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		response := JobResult{
			ID:     999,
			Status: "pending",
			URL:    "/api/v2/jobs/999/",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL, "testuser", "testpass")

	// Test RunJobTemplate
	launch := &JobLaunch{
		Inventory: 123,
		ExtraVars: map[string]interface{}{
			"var1": "value1",
		},
	}
	result, err := client.RunJobTemplate(42, launch)
	if err != nil {
		t.Fatalf("RunJobTemplate failed: %v", err)
	}

	if result.ID != 999 {
		t.Errorf("Expected job ID 999, got %d", result.ID)
	}
	if result.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", result.Status)
	}
}
