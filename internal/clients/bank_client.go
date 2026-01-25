package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gin-quickstart/internal/models"
)

const (
	// BankAPIURL is the external API endpoint for fetching banks.
	BankAPIURL = "https://wema-alatdev-apimgt.developer.azure-api.net/api/GetBanks"
	// BankAPITimeout is the timeout for API requests.
	BankAPITimeout = 10 * time.Second
)

// BankClientInterface defines the contract for bank API client.
type BankClientInterface interface {
	GetBanks() ([]models.Bank, error)
}

// BankClient handles external bank API calls.
type BankClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewBankClient creates a new BankClient.
func NewBankClient() *BankClient {
	return &BankClient{
		httpClient: &http.Client{
			Timeout: BankAPITimeout,
		},
		baseURL: BankAPIURL,
	}
}

// NewBankClientWithURL creates a new BankClient with a custom URL (for testing).
func NewBankClientWithURL(url string) *BankClient {
	return &BankClient{
		httpClient: &http.Client{
			Timeout: BankAPITimeout,
		},
		baseURL: url,
	}
}

// GetBanks fetches the list of banks from the external API.
func (c *BankClient) GetBanks() ([]models.Bank, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch banks: status %d", resp.StatusCode)
	}

	var banks []models.Bank
	if err := json.NewDecoder(resp.Body).Decode(&banks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return banks, nil
}

// GetBanks is a package-level function for backward compatibility.
func GetBanks() ([]models.Bank, error) {
	client := NewBankClient()
	return client.GetBanks()
}
