package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ContactRequest represents the contact form data
type ContactRequest struct {
	Message  string `json:"message"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Discord  string `json:"discord,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Facebook string `json:"facebook,omitempty"`
	Source   string `json:"source"`
}

// ContactResponse represents the API response
type ContactResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Client handles API requests
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SubmitContact sends a contact form submission to the API
func (c *Client) SubmitContact(req ContactRequest) (*ContactResponse, error) {
	// Set source to SSH
	req.Source = "ssh"

	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/contact", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var contactResp ContactResponse
	if err := json.NewDecoder(resp.Body).Decode(&contactResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		if contactResp.Error != "" {
			return &contactResp, fmt.Errorf("API error: %s", contactResp.Error)
		}
		return &contactResp, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	return &contactResp, nil
}
