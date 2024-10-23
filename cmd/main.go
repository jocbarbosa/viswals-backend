package main

import (
	"log"
	"os"

	"github.com/jocbarbosa/viswals-backend/cmd/api"
	"github.com/jocbarbosa/viswals-backend/cmd/reader"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts the API server",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartAPIServer()
	},
}

var readerCmd = &cobra.Command{
	Use:   "reader",
	Short: "Starts the reader",
	Run: func(cmd *cobra.Command, args []string) {
		reader.StartReader()
	},
}

func main() {
	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(readerCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
