package cmd

import (
	"fmt"
	"os"
	"github.com/muka/btrader/service"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a coin",
	Long: `Inspect a coin`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println("Provide a coin symbol as an argument")
			os.Exit(1)
		}

		res, err := service.Inspect(baseCoin, args[0])
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}


		if res.Info != nil {
			fmt.Printf("Avail: %.2f %s\n", res.Info.Free, res.Info.Asset)
			fmt.Printf("Avg. Traded: %.8f %s\n", res.Info.AvgTraded, res.Info.Asset)
		}

		fmt.Printf("Ask: %.8f %s\t\tBid: %.8f %s\n", res.Change.AskPrice, baseCoin, res.Change.BidPrice, baseCoin)
		fmt.Printf("Price: %.8f %s\t\tHigh: %.8f %s\t\tLow: %.8f %s\n", res.Change.LastPrice, baseCoin, res.Change.HighPrice, baseCoin,res.Change.LowPrice, baseCoin)
		fmt.Printf("Volume: %.2f\n", res.Change.Volume)
		fmt.Printf("Change: %.8f %s (%.2f%%)\n", res.Change.PriceChange, baseCoin, res.Change.PriceChangePercent)

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
