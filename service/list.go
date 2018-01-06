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

//ListView contains an overview of coin status
type ListView struct {
	BaseCoin       string
	USDUnitValue   float64
	USDTotalValue  float64
	USDDeltaValue  float64
	USDTradedValue float64
	Coins          []CoinInfo
}

//CoinInfo describe the state of a coin
type CoinInfo struct {
	Asset          string
	Free           float64
	Locked         float64
	Price          float64
	LastTraded     float64
	Delta          float64
	USDValue       float64
	USDDeltaValue  float64
	USDTradedValue float64
}

//List account status
func List(baseCoin string, filter []string) (*ListView, error) {

	symbolFilter := make(map[string]bool)
	if len(filter) > 0 {
		for _, f := range filter {
			symbolFilter[strings.ToUpper(f)] = true
		}
	}

	list := ListView{
		Coins: make([]CoinInfo, 0),
	}

	client := getClient()
	usdChange, err := GetUSD(baseCoin)
	if err != nil {
		return nil, err
	}
	list.USDUnitValue = usdChange

	positions, err := client.GetPositions()
	if err != nil {
		return nil, err
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

		coinInfo := CoinInfo{
			Asset:  balance.Asset,
			Free:   balance.Free,
			Locked: balance.Locked,
		}

		if balance.Asset == baseCoin {
			list.Coins = append(list.Coins, coinInfo)
			continue
		}

		exchangeSymbol := balance.Asset + baseCoin
		lastPrice, err := LastPrice(exchangeSymbol)
		if err != nil {
			coinInfo.Price = -1
			fmt.Printf("[%s] Error: %s\n", coinInfo.Asset, err.Error())
			continue
		} else {
			coinInfo.Price = lastPrice
		}

		lastTrade, err := LastTrade(exchangeSymbol)
		if err != nil {
			coinInfo.LastTraded = -1
			fmt.Printf("[%s] Error: %s\n", coinInfo.Asset, err.Error())
			continue
		} else {
			coinInfo.LastTraded = lastTrade
		}

		if balance.Asset == baseCoin {
			continue
		}

		if lastTrade > -1 {
			coinInfo.Delta = (lastPrice - lastTrade)
			coinInfo.USDDeltaValue = coinInfo.Delta * coinInfo.Free * usdChange
			list.USDDeltaValue += coinInfo.USDDeltaValue
		}

		if coinInfo.LastTraded > -1 {
			coinInfo.USDTradedValue = ((coinInfo.LastTraded * balance.Free) * usdChange)
		}

		coinInfo.USDValue = ((lastPrice * balance.Free) * usdChange)

		list.USDTotalValue += coinInfo.USDValue
		list.USDTradedValue += coinInfo.USDTradedValue

		list.Coins = append(list.Coins, coinInfo)
	}

	return &list, nil
}
