// Crypto - Daily Open/Close
// https://polygon.io/docs/crypto/get_v1_open-close_crypto__from___to___date
// https://github.com/polygon-io/client-go/blob/master/rest/aggs.go
package main

import (
	"context"
	"time"
	"log"
	"os"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

func main() {
	// init client
	c := polygon.New(os.Getenv("POLYGON_API_KEY"))

	// set params
	params := models.GetDailyOpenCloseAggParams{
		Ticker: "X:BTCUSD",
		Date:   models.Date(time.Date(2023, 4, 11, 0, 0, 0, 0, time.Local)),
	}.WithAdjusted(true)

	// make request
	res, err := c.GetDailyOpenCloseAgg(context.Background(), params)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the result
	log.Print(res)
}
