package service

import (
	"sort"
	"strings"

	"github.com/pdepip/go-binance/binance"
)

// spread %
// spread = ((askPrice - bidPrice) / askPrice) x 100
func spread(askPrice, bidPrice float64) float64 {
	return ((askPrice - bidPrice) / askPrice) * 100
}

// AvgDepth contains market depth average for an asset
type AvgDepth struct {
	Asset  string
	Mid    float64
	Ask    float64
	Bid    float64
	Spread float64
}

func getCoinList(baseCoin string) (assets []string, err error) {
	cli := getClient()
	pricesList, err := cli.GetAllPrices()
	if err != nil {
		return nil, err
	}
	for _, c := range pricesList {
		len := len(c.Symbol) - len(baseCoin)
		baseC := c.Symbol[len:]
		sym := c.Symbol[:len]
		if baseC == baseCoin {
			assets = append(assets, sym)
		}
	}
	return assets, nil
}

//Watch market updates
func Watch(baseCoin string, assets []string) (list []AvgDepth, err error) {

	if len(assets) == 0 {
		assets, err = getCoinList(baseCoin)
		if err != nil {
			return list, err
		}
	}

	for _, coin := range assets {
		avg, err := getAvgDepth(baseCoin, coin)
		list = append(list, avg)
		if err != nil {
			return list, err
		}
	}

	sort.Slice(list, func(i int, j int) bool {
		return list[i].Spread > list[j].Spread
	})

	return list, nil
}

func getAvgDepth(baseCoin, coin string) (avg AvgDepth, err error) {

	cli := getClient()
	q := binance.OrderBookQuery{
		Limit:  10,
		Symbol: strings.ToUpper(coin) + baseCoin,
	}

	res, err := cli.GetOrderBook(q)
	if err != nil {
		return avg, err
	}

	sumA := .0
	sumB := .0
	l := len(res.Asks)
	for i := 0; i < l; i++ {
		a := res.Asks[i]
		b := res.Bids[i]
		sumA += a.Price
		sumB += b.Price
		// fmt.Printf("%.0fx %.8f %s\t%.0fx %.8f %s\t%.2f%%\n",
		// 	a.Quantity, a.Price, baseCoin,
		// 	b.Quantity, b.Price, baseCoin,
		// 	spread(a.Price, b.Price),
		// )
	}

	avgA := sumA / float64(l)
	avgB := sumB / float64(l)
	mid := (avgA + avgB) / 2
	avgSpread := spread(avgA, avgB)

	avg.Asset = coin
	avg.Ask = avgA
	avg.Bid = avgB
	avg.Mid = mid
	avg.Spread = avgSpread

	return avg, nil
}
