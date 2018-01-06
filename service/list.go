package service

import (
	"fmt"
	"strings"

	"github.com/pdepip/go-binance/binance"
	"github.com/spf13/viper"
)

var bclient *binance.Binance

func getClient() *binance.Binance {
	if bclient == nil {
		apiKey := viper.GetString("apiKey")
		apiSecret := viper.GetString("apiSecret")
		bclient = binance.New(apiKey, apiSecret)
	}
	return bclient
}

//List account status
func List(baseCoin string, filter []string) error {

	symbolFilter := make(map[string]bool)
	if len(filter) > 0 {
		for _, f := range filter {
			symbolFilter[strings.ToUpper(f)] = true
		}
	}

	client := getClient()

	usdChange, err := GetUSD(baseCoin)
	if err != nil {
		return err
	}

	fmt.Printf("Last %s price: %0.8f", baseCoin, usdChange)

	var total float64

	positions, err := client.GetPositions()
	if err != nil {
		return err
	}
	for _, balance := range positions {

		if len(symbolFilter) > 0 {
			if _, ok := symbolFilter[balance.Asset]; !ok {
				continue
			}
		}

		if balance.Free == 0 && balance.Locked == 0 {
			continue
		}

		availStrf := "%.4f%s\t"
		if balance.Free < 1 { //BTC
			availStrf = "%.8f %s\t"
		}
		availStr := fmt.Sprintf(availStrf, balance.Free, balance.Asset)

		if balance.Asset == baseCoin {
			fmt.Println(availStr)
			continue
		}

		exchangeSymbol := balance.Asset + baseCoin

		lockedStr := ""
		if balance.Locked > 0 {
			lockedStr = fmt.Sprintf("\tLocked: %.8f", balance.Locked)
		}

		lastPriceStr := ""
		lastPrice, err := LastPrice(exchangeSymbol)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		} else {
			lastPriceStr = fmt.Sprintf("\t\tPrice: %.8f %s", lastPrice, baseCoin)
		}

		lastTradeStr := ""
		lastTrade, err := LastTrade(exchangeSymbol)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		} else {
			if lastTrade > -1 {
				lastTradeStr = fmt.Sprintf("\tBought: %.8f %s", lastTrade, baseCoin)
			} else {
				lastTradeStr = "\tBought:           "
			}

		}

		deltaStr := fmt.Sprintf("\tDelta: %.8f", (lastPrice - lastTrade))

		usdVal := ((lastPrice * balance.Free) * usdChange)
		total += usdVal

		valueStr := fmt.Sprintf("\tValue: %.8f", usdVal)

		fmt.Printf("%s%s%s%s%s%s\n", availStr, deltaStr, lockedStr, lastPriceStr, lastTradeStr, valueStr)
	}

	fmt.Printf("---\nTotal value %.2f\n", total)

	return nil
}
