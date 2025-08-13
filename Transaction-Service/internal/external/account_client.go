package external

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type AccountClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAccountClient(baseURL string) *AccountClient {
	return &AccountClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AccountClient) ValidateAccount(ctx context.Context, accountNumber string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/validate/%s", c.baseURL, accountNumber)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		log.Printf("Failed to validate account: %v", err)
		return false, fmt.Errorf("failed to validate account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to validate account: %s", resp.Status)
		return false, fmt.Errorf("failed to validate account: %s", resp.Status)
	}

	var result struct {
		Valid bool `json:"valid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode validation response: %v", err)
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Valid, nil
}

func (c *AccountClient) GetAccountBalance(ctx context.Context, accountNumber string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/balance/%s", c.baseURL, url.QueryEscape(accountNumber))
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get account balance: %s", resp.Status)
		return 0, fmt.Errorf("failed to get account balance: %s", resp.Status)
	}

	var balance float64
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		log.Printf("Failed to decode balance response: %v", err)
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return balance, nil
}