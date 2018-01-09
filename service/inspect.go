package service

import (
	"github.com/pdepip/go-binance/binance"
)

//InspectView contains coin details
type InspectView struct {
	USDValue float64
	Change   binance.ChangeStats
	Info     *CoinInfo
}

//Inspect a coin an return state
func Inspect(baseCoin, symbol string) (*InspectView, error) {
	view := InspectView{}

	change, err := GetChangeStats(symbol + baseCoin)
	if err != nil {
		return nil, err
	}

	view.Change = change

	list, err := List(ListFilter{
		BaseCoin:      baseCoin,
		Asset:         []string{symbol},
		USDValueLimit: 0.001,
	})
	if err != nil {
		return nil, err
	}

	if len(list.Coins) == 1 {
		view.Info = &list.Coins[0]
	}

	view.USDValue = list.USDUnitValue

	return &view, nil
}
