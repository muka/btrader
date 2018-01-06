package cmd

import (
	"github.com/muka/binance-cli/service"
	"github.com/spf13/cobra"
)

var exchangeCoin string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long: `list`,
	Run: func(cmd *cobra.Command, args []string) {
		err := service.List(exchangeCoin, args)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&exchangeCoin, "exchange-coin", "c", "BTC", "Set one of the exchangeable coins")
}
