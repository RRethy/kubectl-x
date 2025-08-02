package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/RRethy/kubectl-x/pkg/cli/cur"
)

var curCmd = &cobra.Command{
	Use:   "cur",
	Short: "Print current context and namespace.",
	Long: `Print current context and namespace.

Usage:
  kubectl x cur

Example:
  kubectl x cur`,
	Run: func(cmd *cobra.Command, args []string) {
		checkErr(cur.Cur(context.Background()))
	},
}

func init() {
	rootCmd.AddCommand(curCmd)
}
