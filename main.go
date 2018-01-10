package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/pkg/errors"

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
		httpManager := HttpManager{}
		if price, err := runPrice(httpManager, *priceRaw); err != nil {
			app.Errorf("%s", err)
		} else {
			fmt.Print(price)
		}
	}
}

func runPrice(httpManager HttpRequester, rawFormat bool) (string, error) {
	_, body, errs := httpManager.Get()
	if errs != nil {
		return "", errors.New(fmt.Sprintf("Unable to get Bitcoin price in USD.\n%v\n", errs))
	}

	ticker := BTCTicker{}
	err := json.Unmarshal([]byte(body), &ticker)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to parse ticker data. %v\n", err))
	}

	if price, err := strconv.ParseFloat(ticker.Price, 32); err != nil {
		return "", errors.New(fmt.Sprintf("Unable to parse the ticker price. %v\n", err))
	} else {
		if rawFormat {
			return fmt.Sprintf("%f\n", price), nil
		} else {
			bigPrice := big.NewFloat(price)
			ac := accounting.Accounting{Symbol: "$", Precision: 2}
			return fmt.Sprintf("%s USD\n", ac.FormatMoneyBigFloat(bigPrice)), nil
		}
	}

}
