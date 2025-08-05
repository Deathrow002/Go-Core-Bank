package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// ValidateCustomer calls Customer Service to validate if customer exists
func (c *CustomerClient) ValidateCustomer(ctx context.Context, customerID uuid.UUID, authToken string) (bool, error) {
    // Prepare request payload
    reqBody := CustomerValidationRequest{
        CustomerID: customerID,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return false, fmt.Errorf("failed to marshal request: %w", err)
    }

    // Create HTTP request
    url := fmt.Sprintf("%s/api/v1/customers/validate", c.baseURL)
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return false, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", "Bearer "+ authToken)
    log.Println("Authorization Header in Client Call:", authToken) // Print to console


    // Make the request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return false, fmt.Errorf("failed to call customer service: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, fmt.Errorf("failed to read response: %w", err)
    }

    // Check status code
    if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("customer service returned status %d: %s", resp.StatusCode, string(body))
    }

    // Parse response
    var response CustomerValidationResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return false, fmt.Errorf("failed to decode response: %w", err)
    }

    return response.Exists, nil
}