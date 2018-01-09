package service

import (
	"fmt"
	"strings"

	"github.com/pdepip/go-binance/binance"
)

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
	AvgTraded      float64
	Delta          float64
	USDValue       float64
	USDDeltaValue  float64
	USDTradedValue float64
	Change         binance.ChangeStats
	Trades         []binance.Trade
}

//ListFilter filter for the list view
type ListFilter struct {
	BaseCoin      string
	Asset         []string
	USDValueLimit float64
}

//List account status
func List(filter ListFilter) (*ListView, error) {

	baseCoin := filter.BaseCoin

	if filter.USDValueLimit == 0 {
		filter.USDValueLimit = 5
	}

	symbolFilter := make(map[string]bool)
	if len(filter.Asset) > 0 {
		for _, f := range filter.Asset {
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

	pricesList, err := client.GetAllPrices()
	if err != nil {
		return nil, err
	}

	prices := make(map[string]float64)
	for _, p := range pricesList {
		prices[p.Symbol] = p.Price
	}

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

		change, err := GetChangeStats(exchangeSymbol)
		if err != nil {
			return nil, err
		}
		coinInfo.Change = change

		lastPrice, ok := prices[exchangeSymbol]
		if !ok {
			coinInfo.Price = -1
			continue
		} else {
			coinInfo.Price = lastPrice
		}

		coinInfo.USDValue = ((lastPrice * balance.Free) + (lastPrice * balance.Locked)) * usdChange

		if coinInfo.USDValue < filter.USDValueLimit {
			continue
		}

		trades, err := LastTrades(exchangeSymbol)
		if err != nil || len(trades) == 0 {
			coinInfo.AvgTraded = -1
			fmt.Printf("[%s] Error: %s\n", coinInfo.Asset, err.Error())
			continue
		} else {
			var avgPrice float64
			var buyOp float64
			for _, t := range trades {
				if t.IsBuyer {
					avgPrice += t.Price
					buyOp++
				}
			}
			if buyOp > 0 {
				coinInfo.AvgTraded = avgPrice / buyOp
			}
		}
		coinInfo.Trades = trades

		if balance.Asset == baseCoin {
			continue
		}

		if coinInfo.AvgTraded > -1 {
			coinInfo.Delta = (lastPrice - coinInfo.AvgTraded)
			coinInfo.USDDeltaValue = ((coinInfo.Delta * coinInfo.Free) + (coinInfo.Delta * coinInfo.Locked)) * usdChange
			list.USDDeltaValue += coinInfo.USDDeltaValue
		}

		if coinInfo.AvgTraded > -1 {
			coinInfo.USDTradedValue = ((coinInfo.AvgTraded * balance.Free) + (coinInfo.AvgTraded * balance.Locked)) * usdChange
		}

		list.USDTotalValue += coinInfo.USDValue
		list.USDTradedValue += coinInfo.USDTradedValue

		list.Coins = append(list.Coins, coinInfo)
	}

	return &list, nil
}
