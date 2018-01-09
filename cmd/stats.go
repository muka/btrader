package cmd

import (
	"fmt"
	"github.com/muka/btrader/service"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show market stats for a coin",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		stats, err := service.Stats(service.StatsFilter{
			BaseCoin: baseCoin,
			Asset: args,
		})
		fail(err)

		fmt.Printf("\n%++v\n\n", stats)

	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
