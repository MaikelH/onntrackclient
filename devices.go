// Package onntrackclient provides a client for the Jimi tracking dashboard REST API.
package onntrackclient

import (
	"context"
	"fmt"
	"net/http"
)

// DevicesService handles communication with the device related
// methods of the Jimi API.
type DevicesService service

// Device represents a Jimi tracking device.
type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	IMEI        string `json:"imei"`
	Status      string `json:"status"`
	LastUpdated string `json:"last_updated"`
	// Add other device fields as needed
}

// DeviceListOptions specifies the optional parameters to the
// DevicesService.List method.
type DeviceListOptions struct {
	// Page number for pagination
	Page int `url:"page,omitempty"`

	// Number of results per page
	PerPage int `url:"per_page,omitempty"`

	// Filter by device status
	Status string `url:"status,omitempty"`
}

// List devices.
//
// Jimi API docs: [URL to API documentation]
func (s *DevicesService) List(ctx context.Context, opts *DeviceListOptions) ([]*Device, *http.Response, error) {
	u := "devices"
	// TODO: Add query parameters from opts

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var devices []*Device
	resp, err := s.client.Do(req, &devices)
	if err != nil {
		return nil, resp, err
	}

	return devices, resp, nil
}

// Get a single device.
//
// Jimi API docs: [URL to API documentation]
func (s *DevicesService) Get(ctx context.Context, deviceID string) (*Device, *http.Response, error) {
	u := fmt.Sprintf("devices/%s", deviceID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	device := new(Device)
	resp, err := s.client.Do(req, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, nil
}

// DeviceCreateRequest represents a request to create a device.
type DeviceCreateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	IMEI string `json:"imei"`
	// Add other required fields
}

// Create a new device.
//
// Jimi API docs: [URL to API documentation]
func (s *DevicesService) Create(ctx context.Context, deviceReq *DeviceCreateRequest) (*Device, *http.Response, error) {
	u := "devices"

	req, err := s.client.NewRequest(ctx, http.MethodPost, u, deviceReq)
	if err != nil {
		return nil, nil, err
	}

	device := new(Device)
	resp, err := s.client.Do(req, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, nil
}

// DeviceUpdateRequest represents a request to update a device.
type DeviceUpdateRequest struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	// Add other fields that can be updated
}

// Update a device.
//
// Jimi API docs: [URL to API documentation]
func (s *DevicesService) Update(ctx context.Context, deviceID string, deviceReq *DeviceUpdateRequest) (*Device, *http.Response, error) {
	u := fmt.Sprintf("devices/%s", deviceID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, u, deviceReq)
	if err != nil {
		return nil, nil, err
	}

	device := new(Device)
	resp, err := s.client.Do(req, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, nil
}

// Delete a device.
//
// Jimi API docs: [URL to API documentation]
func (s *DevicesService) Delete(ctx context.Context, deviceID string) (*http.Response, error) {
	u := fmt.Sprintf("devices/%s", deviceID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
