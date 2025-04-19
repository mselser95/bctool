package prices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// HTTPClient defines the interface to mock HTTP calls
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client holds configuration and HTTP client
type Client struct {
	BaseURL string
	Client  HTTPClient
	Suffix  string
}

// NewClient constructs a new price-fetching client
func NewClient(httpClient HTTPClient) *Client {
	return &Client{
		BaseURL: "https://api.gateio.ws/api/v4/spot/tickers",
		Client:  httpClient,
		Suffix:  "_USDT",
	}
}

// TickerResponse represents an API response item
type TickerResponse struct {
	CurrencyPair string `json:"currency_pair"`
	Last         string `json:"last"` // Last traded price
}

// GetPrices fetches the price for each ticker concurrently.
// Returns a map with results and aborts early on critical HTTP/parse errors.
func (c *Client) GetPrices(currencyPairs []string) (map[string]float64, error) {
	results := make(map[string]float64)
	var mu sync.Mutex
	var wg sync.WaitGroup

	errs := make(chan error, len(currencyPairs))

	for _, ticker := range currencyPairs {
		tickerWithSuffix := ticker + c.Suffix
		wg.Add(1)

		go func(original, formatted string) {
			defer wg.Done()

			price, err := c.fetchPrice(formatted)
			if err != nil {
				errs <- fmt.Errorf("error fetching %s: %w", original, err)
				return
			}

			mu.Lock()
			results[original] = price
			mu.Unlock()
		}(ticker, tickerWithSuffix)
	}

	wg.Wait()
	close(errs)

	// Return first error if any occurred
	for err := range errs {
		return nil, err
	}

	return results, nil
}

// fetchPrice performs an HTTP request to get the price of a currency pair
func (c *Client) fetchPrice(ticker string) (float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?currency_pair=%s", c.BaseURL, ticker), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("API returned status code %d for %s", resp.StatusCode, ticker)
	}

	var data []TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("no data returned for %s", ticker)
	}

	price, err := strconv.ParseFloat(data[0].Last, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price format for %s: %w", ticker, err)
	}

	return price, nil
}
