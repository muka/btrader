package service

import (
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
func LastTrade(symbol string) ([]binance.Trade, error) {
	t := []binance.Trade{}
	client := getClient()
	res, err := client.GetTrades(symbol)
	if err != nil {
		return t, err
	}
	if len(res) == 0 {
		return t, nil
	}
	return res, nil
}

//GetUSD return the current price in USDT of a coin
func GetUSD(baseCoin string) (float64, error) {
	return LastPrice(baseCoin + "USDT")
}

//GetChangeStats return 24h stats for a coin
func GetChangeStats(symbol string) (binance.ChangeStats, error) {

	client := getClient()

	q := binance.SymbolQuery{
		Symbol: symbol,
	}

	return client.Get24Hr(q)
}
