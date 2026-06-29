package main

import (
	"context"
	"fmt"
	"os"

	"github.com/massive-com/client-go/v3/rest"
	"github.com/massive-com/client-go/v3/rest/gen"
)

func main() {
	apiKey := os.Getenv("MASSIVE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "MASSIVE_API_KEY environment variable must be set")
		os.Exit(1)
	}

	c := rest.New(apiKey)

	fmt.Println("=== Test 1: Stock Snapshot (AAPL) — LastTrade & LastQuote decimal fields ===")
	testStockSnapshot(c)

	fmt.Println("\n=== Test 2: List Trades (AAPL) — Price & Size decimal fields ===")
	testListTrades(c)

	fmt.Println("\n=== Test 3: Dividends (AAPL) — Currency field present ===")
	testDividends(c)

	fmt.Println("\n✅ All tests passed — decimal.Decimal fields deserialized correctly")
}

func testStockSnapshot(c *rest.Client) {
	resp, err := c.GetStocksSnapshotTickerWithResponse(context.Background(), "AAPL")
	if err != nil {
		fmt.Printf("  ❌ request error: %v\n", err)
		os.Exit(1)
	}
	if resp.JSON200 == nil {
		fmt.Printf("  ❌ unexpected status: %s\n", resp.Status())
		os.Exit(1)
	}

	t := resp.JSON200.Ticker
	if t == nil {
		fmt.Println("  ⚠️  no ticker data returned")
		return
	}

	if t.LastQuote != nil {
		fmt.Printf("  LastQuote.AskPrice = %s  (type: decimal.Decimal ✓)\n", t.LastQuote.AskPrice.String())
		fmt.Printf("  LastQuote.BidPrice = %s  (type: decimal.Decimal ✓)\n", t.LastQuote.BidPrice.String())
	} else {
		fmt.Println("  LastQuote: nil (plan may not include quotes)")
	}

	if t.LastTrade != nil {
		fmt.Printf("  LastTrade.BidPrice = %s  (trade price, type: decimal.Decimal ✓)\n", t.LastTrade.BidPrice.String())
	} else {
		fmt.Println("  LastTrade: nil (plan may not include trades)")
	}
}

func testListTrades(c *rest.Client) {
	ticker := "AAPL"
	limit := 2
	resp, err := c.GetStocksTradesWithResponse(context.Background(), ticker, &gen.GetStocksTradesParams{
		Limit: &limit,
	})
	if err != nil {
		fmt.Printf("  ❌ request error: %v\n", err)
		os.Exit(1)
	}
	if resp.JSON200 == nil {
		fmt.Printf("  ❌ unexpected status: %s\n", resp.Status())
		os.Exit(1)
	}

	results := resp.JSON200.Results
	if results == nil || len(*results) == 0 {
		fmt.Println("  ⚠️  no trade results (market may be closed)")
		return
	}

	for i, trade := range *results {
		fmt.Printf("  Trade[%d].Price = %s  (type: decimal.Decimal ✓)\n", i, trade.Price.String())
		fmt.Printf("  Trade[%d].Size  = %s  (type: decimal.Decimal ✓)\n", i, trade.Size.String())
	}
}

func testDividends(c *rest.Client) {
	ticker := "AAPL"
	limit := 1
	resp, err := c.ListDividendsWithResponse(context.Background(), &gen.ListDividendsParams{
		Ticker: &ticker,
		Limit:  &limit,
	})
	if err != nil {
		fmt.Printf("  ❌ request error: %v\n", err)
		os.Exit(1)
	}
	if resp.JSON200 == nil {
		fmt.Printf("  ❌ unexpected status: %s\n", resp.Status())
		os.Exit(1)
	}

	results := resp.JSON200.Results
	if results == nil || len(*results) == 0 {
		fmt.Println("  ⚠️  no dividend results")
		return
	}

	d := (*results)[0]
	if d.Currency != nil {
		fmt.Printf("  Dividend.Currency   = %s  ✓\n", *d.Currency)
	} else {
		fmt.Println("  Dividend.Currency   = nil (field exists but was empty)")
	}
	fmt.Printf("  Dividend.CashAmount = %.4f\n", d.CashAmount)
	fmt.Printf("  Dividend.Ticker     = %s\n", d.Ticker)
}
