package service

import (
	"github.com/pdepip/go-binance/binance"
)

//LastPrice return the last price for a symbol
func LastPrice(symbol string) (float64, error) {
	query := binance.SymbolQuery{
		Symbol: symbol,
	}
	client := binance.New("", "")
	res, err := client.GetLastPrice(query)
	if err != nil {
		return -1, err
	}
	return res.Price, nil
}

//LastTrade return the last price bought
func LastTrade(symbol string) (float64, error) {
	client := getClient()
	res, err := client.GetTrades(symbol)
	if err != nil {
		return -1, err
	}
	if len(res) == 0 {
		return -1, nil
	}
	return res[0].Price, nil
}

//GetUSD return the current price in USDT of a coin
func GetUSD(baseCoin string) (float64, error) {
	return LastPrice(baseCoin + "USDT")
}
