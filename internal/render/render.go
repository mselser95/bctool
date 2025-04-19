package render

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/leekchan/accounting"
	"io"
	"strconv"
)

// RenderPricesTable renders prices in a tabular format
func RenderPricesTable(results map[string]float64, out io.Writer) {
	t := table.NewWriter()
	t.SetOutputMirror(out)
	t.AppendHeader(table.Row{"Ticker", "Price USD"})

	ac := accounting.Accounting{Symbol: "$", Precision: 2, Thousand: ",", Decimal: "."}

	for ticker, price := range results {
		t.AppendRow(table.Row{ticker, formatPrice(price, &ac)})
	}

	t.Render()
}

// formatPrice formats a price based on its magnitude
func formatPrice(price float64, ac *accounting.Accounting) string {
	if price < 0.01 {
		return "$" + strconv.FormatFloat(price, 'f', 6, 64)
	}
	return ac.FormatMoney(price)
}
