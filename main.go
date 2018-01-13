package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/leekchan/accounting"
)

func main() {
	app := kingpin.New("bitcoin", "A command-line Bitcoin application.")
	price := app.Command("price", "Output the current price of Bitcoin in USD.")
	priceRaw := price.Flag("raw", "Display only the number without currency symbols.").Bool()

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Errorf("bitcoin: error: %s\n", err)
		app.Usage([]string{})
		return
	}

	switch command {
	case price.FullCommand():
		httpManager := HTTPManager{}
		price, err := runPrice(httpManager, *priceRaw)
		if err != nil {
			app.Errorf("%s\n", err)
			return
		}
		fmt.Print(price)
	}
}

func runPrice(httpManager HTTPRequester, rawFormat bool) (string, error) {
	_, body, errs := httpManager.Get()
	if errs != nil {
		return "", fmt.Errorf("unable to get bitcoin price in USD: %v", errs)
	}

	ticker := BTCTicker{}
	err := json.Unmarshal([]byte(body), &ticker)
	if err != nil {
		return "", fmt.Errorf("unable to parse ticker data: %v", err)
	}

	price, _, err := big.ParseFloat(ticker.Price, 10, 100, big.AwayFromZero)
	if err != nil {
		return "", fmt.Errorf("unable to parse the ticker price: %v", err)
	}

	if rawFormat {
		return fmt.Sprintf("%.8f\n", price), nil
	}

	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return fmt.Sprintf("%s USD\n", ac.FormatMoneyBigFloat(price)), nil
}
