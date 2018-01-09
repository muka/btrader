package service

import (
	"fmt"
	"os"
	"time"

	"github.com/pdepip/go-binance/binance"
)

//StatsFilter filter applied on stats retrieval
type StatsFilter struct {
	BaseCoin string
	Asset    []string
}

//StatsView wrap market stats details
type StatsView struct {
}

//Stats elaborate market stats
func Stats(filter StatsFilter) (*StatsView, error) {
	stats := &StatsView{}

	baseCoin := filter.BaseCoin
	cli := getClient()

	dates := []time.Time{
		time.Now(),
	}

	for _, symbol := range filter.Asset {
		exchangeSymbol := symbol + baseCoin
		for _, t := range dates {

			endTime := int64(t.UnixNano() / 1000)
			startTime := endTime - (30 * 60 * 1000) // 30m

			q := binance.AggTradesQuery{
				Symbol:    exchangeSymbol,
				StartTime: startTime,
				EndTime:   endTime,
			}
			agg, err := cli.GetAggTrades(q)
			if err != nil {
				return nil, err
			}

			fmt.Printf("%++v", agg)
			os.Exit(0)
		}
	}

	return stats, nil
}
