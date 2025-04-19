package prices_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/mselser95/bctool/internal/prices"
	"github.com/stretchr/testify/require"
)

// mockHTTPClient implements prices.HTTPClient
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetPrices(t *testing.T) {
	t.Run("single ticker success", func(t *testing.T) {
		mockResp := `[{"currency_pair":"BTC_USDT","last":"42000.5"}]`
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeMockResponse(http.StatusOK, mockResp), nil
			},
		})

		result, err := client.GetPrices([]string{"BTC"})
		require.NoError(t, err)
		require.Equal(t, 42000.5, result["BTC"])
	})

	t.Run("multiple tickers success", func(t *testing.T) {
		responses := map[string]string{
			"BTC_USDT": `[{"currency_pair":"BTC_USDT","last":"42000.5"}]`,
			"ETH_USDT": `[{"currency_pair":"ETH_USDT","last":"3000.99"}]`,
		}
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeMockResponse(http.StatusOK, responses[req.URL.Query().Get("currency_pair")]), nil
			},
		})

		result, err := client.GetPrices([]string{"BTC", "ETH"})
		require.NoError(t, err)
		require.Equal(t, 42000.5, result["BTC"])
		require.Equal(t, 3000.99, result["ETH"])
	})

	t.Run("malformed JSON", func(t *testing.T) {
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeMockResponse(http.StatusOK, `{"bad":`), nil
			},
		})

		_, err := client.GetPrices([]string{"BTC"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "decode")
	})

	t.Run("missing price field", func(t *testing.T) {
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeMockResponse(http.StatusOK, `[{"currency_pair":"BTC_USDT"}]`), nil
			},
		})

		_, err := client.GetPrices([]string{"BTC"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid price")
	})

	t.Run("http 500 error", func(t *testing.T) {
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeMockResponse(http.StatusInternalServerError, ""), nil
			},
		})

		_, err := client.GetPrices([]string{"BTC"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "API returned status code 500")
	})

	t.Run("http 429 error", func(t *testing.T) {
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return makeMockResponse(http.StatusTooManyRequests, ""), nil
			},
		})

		_, err := client.GetPrices([]string{"BTC"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "API returned status code 429")
	})

	t.Run("network failure", func(t *testing.T) {
		client := prices.NewClient(&mockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network is down")
			},
		})

		_, err := client.GetPrices([]string{"BTC"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "network is down")
	})
}

// makeMockResponse creates a mock HTTP response from body string
func makeMockResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
