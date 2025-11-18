// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an Ansible Automation Platform API client
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewClient creates a new AAP client
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest executes an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// Inventory represents an AAP inventory
type Inventory struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Organization int    `json:"organization"`
}

// CreateInventory creates a new inventory in AAP
func (c *Client) CreateInventory(name, description string, organizationID int) (*Inventory, error) {
	inventory := &Inventory{
		Name:         name,
		Description:  description,
		Organization: organizationID,
	}

	resp, err := c.doRequest("POST", "/api/v2/inventories/", inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create inventory: status %d, body: %s", resp.StatusCode, string(body))
	}

	var createdInventory Inventory
	if err := json.NewDecoder(resp.Body).Decode(&createdInventory); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdInventory, nil
}

// Host represents an AAP host
type Host struct {
	ID          int                    `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Inventory   int                    `json:"inventory,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
}

// CreateHost creates a new host without associating it to an inventory
func (c *Client) CreateHost(hostname, description string, variables map[string]interface{}) (*Host, error) {
	host := &Host{
		Name:        hostname,
		Description: description,
		Variables:   variables,
	}

	resp, err := c.doRequest("POST", "/api/v2/hosts/", host)
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create host: status %d, body: %s", resp.StatusCode, string(body))
	}

	var createdHost Host
	if err := json.NewDecoder(resp.Body).Decode(&createdHost); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdHost, nil
}

// AddHost adds a host to an inventory
func (c *Client) AddHost(inventoryID int, hostname, description string, variables map[string]interface{}) (*Host, error) {
	host := &Host{
		Name:        hostname,
		Description: description,
		Inventory:   inventoryID,
		Variables:   variables,
	}

	resp, err := c.doRequest("POST", "/api/v2/hosts/", host)
	if err != nil {
		return nil, fmt.Errorf("failed to add host: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to add host: status %d, body: %s", resp.StatusCode, string(body))
	}

	var createdHost Host
	if err := json.NewDecoder(resp.Body).Decode(&createdHost); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdHost, nil
}

// JobLaunch represents parameters for launching a job
type JobLaunch struct {
	Inventory     int                    `json:"inventory,omitempty"`
	Limit         string                 `json:"limit,omitempty"`
	ExtraVars     map[string]interface{} `json:"extra_vars,omitempty"`
	SkipTags      string                 `json:"skip_tags,omitempty"`
	JobTags       string                 `json:"job_tags,omitempty"`
	Verbosity     int                    `json:"verbosity,omitempty"`
	DiffMode      bool                   `json:"diff_mode,omitempty"`
	ScmBranch     string                 `json:"scm_branch,omitempty"`
}

// JobResult represents the result of a job launch
type JobResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

// RunJobTemplate launches a job template
func (c *Client) RunJobTemplate(templateID int, launch *JobLaunch) (*JobResult, error) {
	path := fmt.Sprintf("/api/v2/job_templates/%d/launch/", templateID)

	resp, err := c.doRequest("POST", path, launch)
	if err != nil {
		return nil, fmt.Errorf("failed to run job template: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to run job template: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result JobResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
