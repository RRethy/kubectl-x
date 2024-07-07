package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/RRethy/kubectl-x/internal/cmd/ns"
)

var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "Switch namespace.",
	Long: `Switch namespace.

Usage:
  kubectl x ns [namespace]

Args:
  namespace  Partial match to filter namespaces on.

Example:
  kubectl-pi ns
  kubectl-pi ns my-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string
		if len(args) > 0 {
			namespace = args[0]
		}

		checkErr(ns.Ns(context.Background(), configFlags, namespace, exactMatch))
	},
}

func init() {
	rootCmd.AddCommand(nsCmd)
	nsCmd.Flags().BoolVarP(&exactMatch, "exact", "e", false, "Exact match")
}
