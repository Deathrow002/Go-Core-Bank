package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// AccountClient handles communication with Account Service
type AccountClient struct {
	baseURL    string
	httpClient *http.Client
}

// AccountResponse matches the Account Service response
type AccountResponse struct {
	ID            string  `json:"id"`
	CustomerID    string  `json:"customer_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Status        string  `json:"status"`
}

// NewAccountClient creates a new account service client
func NewAccountClient(baseURL string) *AccountClient {
	return &AccountClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAccountBalance calls Account Service to get the account balance
func (c *AccountClient) GetAccountBalance(ctx context.Context, accountID string, authToken string) (float64, error) {
	// Log the account ID being requested
	log.Printf("Getting balance for account ID: %s", accountID)

	// Create HTTP request
	url := fmt.Sprintf("%s/api/v1/accounts/%s", c.baseURL, accountID)
	log.Printf("Making request to: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken) // Add this line
	log.Printf("Using authorization token: %s", authToken)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Error calling account service: %v", err)
		return 0, fmt.Errorf("failed to call account service: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Account service returned status %d: %s", resp.StatusCode, body)
		return 0, fmt.Errorf("account service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var account AccountResponse
	if err := json.Unmarshal(body, &account); err != nil {
		log.Printf("Error decoding response: %v", err)
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Account balance retrieved: %.2f", account.Balance)
	return account.Balance, nil
}