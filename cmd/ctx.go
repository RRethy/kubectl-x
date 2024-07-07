package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/RRethy/kubectl-x/internal/cmd/ctx"
)

var ctxCmd = &cobra.Command{
	Use:   "ctx",
	Short: "Switch context.",
	Long: `Switch context.

Usage:
  kubectl x ctx [context] [namespace]

Args:
  context    Partial match to filter contexts on.
  namespace  Partial match to filter namespaces on.

Example:
  kubectl-pi ctx
  kubectl-pi ctx my-context
  kubectl-pi ctx my-context my-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		var contextName string
		var namespace string
		if len(args) > 0 {
			contextName = args[0]
			if len(args) > 1 {
				namespace = args[1]
			}
		}

		checkErr(ctx.Ctx(context.Background(), configFlags, resourceBuilderFlags, contextName, namespace, exactMatch))
	},
}

func init() {
	rootCmd.AddCommand(ctxCmd)
	ctxCmd.Flags().BoolVarP(&exactMatch, "exact", "e", false, "Exact match")
	resourceBuilderFlags.AddFlags(ctxCmd.Flags())
}
