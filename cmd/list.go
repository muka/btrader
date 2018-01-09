package cmd

import (
	"fmt"
	"github.com/muka/btrader/service"
	"github.com/spf13/cobra"
)

func renderList(list *service.ListView) {
	fmt.Println("---------------------------------------------------------------------------------------------------")
	fmt.Println("Coin\t\t\tDelta\t\t\tPrice\t\t\t\tAvg.Bought")
	fmt.Println("---------------------------------------------------------------------------------------------------")
	for _, balance := range list.Coins {

		availStrf := "%.2f %s\t"
		if balance.Free < 1 { //BTC
			availStrf = "%.8f %s\t"
		}
		availStr := fmt.Sprintf(availStrf, balance.Free, balance.Asset)

		if balance.Asset == baseCoin {
			fmt.Println(availStr)
			continue
		}

		lockedStr := "   "
		if balance.Locked > 0 {
			// lockedStr = fmt.Sprintf("\tLocked: %.8f", balance.Locked)
			lockedStr = " * "
		}

		lastPriceStr := ""
		lastPrice := balance.Price
		if lastPrice > -1 {
			pad := ""
			if balance.USDValue < 100 {
				pad = "   "
			}
			lastPriceStr = fmt.Sprintf("\t%.8f (%.2f$)%s", lastPrice, balance.USDValue, pad)
		}

		var lastTradedStr string
		lastTraded := balance.AvgTraded
		if lastTraded > -1 {
			lastTradedStr = fmt.Sprintf("\t\t%.8f (%.2f$)", lastTraded, balance.USDTradedValue)
		} else {
			lastTradedStr = "\t\t           "
		}

		deltaStr := fmt.Sprintf("\t%.8f (%.2f$)", balance.Delta, balance.USDDeltaValue)

		fmt.Printf("%s%s%s%s%s\n", availStr, deltaStr, lockedStr, lastPriceStr, lastTradedStr)
	}

}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long: `list`,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := service.List(service.ListFilter{
			BaseCoin: baseCoin,
			Asset: args,
		})
		fail(err)

		fmt.Printf("\nLast %s price: %.8f\n\n", baseCoin, list.USDUnitValue)
		renderList(list)
		fmt.Println("---------------------------------------------------------------------------------------------------")

		color := "33m"
		if list.USDDeltaValue > 0 {
			color = "32m+"
		}
		fmt.Printf("Total %.2f$ (\033[%s%.2f$\033[0m)\n", list.USDTotalValue, color, list.USDDeltaValue)
		// fmt.Printf("Traded %.2f$\n", list.USDTradedValue)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
