package cmd

import (
	"net/http"

	"github.com/mselser95/bctool/internal/prices"
	"github.com/mselser95/bctool/internal/render"
	"github.com/spf13/cobra"
)

// pricesCmd represents the prices command
var pricesCmd = &cobra.Command{
	Use:   "prices [TICKERS...]",
	Short: "Fetch cryptocurrency prices from Gate.io",
	Long: `Retrieve the current market prices for one or more cryptocurrencies
against USD from Gate.io's API.

Example:
  bctool prices TRC BTC ETH LINK UNI

+--------+------------+
| TICKER | PRICE USD  |
+--------+------------+
| ETH    | $1,589.34  |
| LINK   | $12.57     |
| BTC    | $84,483.00 |
| TRC    | $0.001580  |
| UNI    | $5.18      |
+--------+------------+

`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Instantiate the price client with the default HTTP client
		client := prices.NewClient(http.DefaultClient)

		// Fetch prices
		results, err := client.GetPrices(args)
		if err != nil {
			cmd.PrintErrf("Error fetching prices: %v\n", err)
		}

		// Render output
		render.RenderPricesTable(results, cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(pricesCmd)
}
