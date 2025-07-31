package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// CustomerClient handles communication with Customer Service
type CustomerClient struct {
    baseURL    string
    httpClient *http.Client
}

// CustomerValidationRequest matches the Customer Service request
type CustomerValidationRequest struct {
    CustomerID uuid.UUID `json:"customer_id"`
}

// CustomerValidationResponse matches the Customer Service response
type CustomerValidationResponse struct {
    CustomerID uuid.UUID `json:"customer_id"`
    Exists     bool      `json:"exists"`
}

// NewCustomerClient creates a new customer service client
func NewCustomerClient(baseURL string) *CustomerClient {
	return &CustomerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CustomerClient) ValidateCustomer(ctx context.Context, customerID uuid.UUID) (bool, error) {
	// Prepare request payload
	reqBody := CustomerValidationRequest{
		CustomerID: customerID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := c.baseURL + "/api/v1/customers/validate"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response CustomerValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Exists, nil
}