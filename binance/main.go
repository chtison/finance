package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/chtison/finance/binance/pkg/wss"
)

var rootCmd = &cobra.Command{
	Version: "0.0.1",
	Use:     fmt.Sprintf("%s", filepath.Base(os.Args[0])),
	Short:   "Access binance APIs",
}

func main() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.AddCommand(wss.NewCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
