package cmd

import (
	"fmt"
	"os"
	"github.com/muka/btrader/service"
	"github.com/spf13/cobra"
)

var baseCoin string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long: `list`,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := service.List(baseCoin, args)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Last %s price: %.8f\n", baseCoin, list.USDUnitValue)

		fmt.Println("---\nCoin\t\t\tDelta\t\t\tPrice\t\t\t\tBought\n---")

		for _, balance := range list.Coins {

			availStrf := "%.4f %s\t"
			if balance.Free < 1 { //BTC
				availStrf = "%.8f %s\t"
			}
			availStr := fmt.Sprintf(availStrf, balance.Free, balance.Asset)

			if balance.Asset == baseCoin {
				fmt.Println(availStr)
				continue
			}

			lockedStr := ""
			// if balance.Locked > 0 {
			// 	lockedStr = fmt.Sprintf("\tLocked: %.8f", balance.Locked)
			// }

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
			lastTraded := balance.LastTraded
			if lastTraded > -1 {
				lastTradedStr = fmt.Sprintf("\t\t%.8f (%.2f$)", lastTraded, balance.USDTradedValue)
			} else {
				lastTradedStr = "\t\t           "
			}

			deltaStr := fmt.Sprintf("\t%.8f (%.2f$)", balance.Delta, balance.USDDeltaValue)

			fmt.Printf("%s%s%s%s%s\n", availStr, deltaStr, lockedStr, lastPriceStr, lastTradedStr)
		}

		fmt.Printf("---\nTotal value %.2f$\n", list.USDTotalValue)
		fmt.Printf("Delta value %.2f$\n", list.USDDeltaValue)
		fmt.Printf("Bought value %.2f$\n", list.USDTradedValue)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&baseCoin, "exchange-coin", "c", "BTC", "Set one of the exchangeable coins")
}
