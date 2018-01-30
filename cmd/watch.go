package cmd

import (
	"fmt"
	"github.com/muka/btrader/service"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch market events",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := service.Watch(baseCoin, args)
		fail(err)

		fmt.Println("Asset\tSpread\tMid\t\tAsks\t\tBid")
		for _, avg := range list {
			if avg.Spread > 1 {
				fmt.Printf("%s\t%.2f%%\t%.8f\t%.8f\t%.8f\n", avg.Asset, avg.Spread, avg.Mid, avg.Ask, avg.Bid)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
