package wss

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wss",
		Short: "Access binance secured websocket API",
		Long: `Access binance secured websocket API

https://github.com/binance-exchange/binance-official-api-docs/blob/master/web-socket-streams.md`,
	}
	cmd.AddCommand(getKlineStreamCommand())
	return cmd
}
