package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/timrourke/maschera/m/v2/deps"
)

var rootCmd = &cobra.Command{
	Use:   "maschera",
	Short: "maschera is a tool for masking PII in streaming data",
	Long:  "maschera is a tool for masking PII in streaming data",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := deps.BuildApp()

		ctx := context.Background()

		return app.Run(ctx)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
