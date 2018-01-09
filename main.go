package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/parnurzeal/gorequest"

	"github.com/leekchan/accounting"
)

var (
	app      = kingpin.New("bitcoin", "A command-line Bitcoin application.")
	price    = app.Command("price", "Output the current price of Bitcoin in USD.")
	priceRaw = price.Flag("raw", "Display only the number without currency symbols.").Bool()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case price.FullCommand():
		runPrice(*priceRaw)
	}
}

func runPrice(rawFormat bool) {
	_, body, errs := gorequest.New().Get("https://api.gdax.com/products/BTC-USD/ticker").End()
	if errs != nil {
		app.Errorf("Unable to get Bitcoin price in USD.\n%v\n", errs)
		return
	}

	ticker := BTCTicker{}
	err := json.Unmarshal([]byte(body), &ticker)
	if err != nil {
		app.Errorf("Unable to parse ticker data. %v\n", err)
		return
	}

	if price, err := strconv.ParseFloat(ticker.Price, 32); err != nil {
		app.Errorf("Unable to parse the ticker price. %v\n", err)
	} else {
		if rawFormat {
			fmt.Printf("%f\n", price)
		} else {
			bigPrice := big.NewFloat(price)
			ac := accounting.Accounting{Symbol: "$", Precision: 2}
			fmt.Printf("%s USD\n", ac.FormatMoneyBigFloat(bigPrice))
		}
	}

}
