package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/RRethy/kubectl-x/internal/cmd/shell"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open up a shell.",
	Long: `Open up a shell.

Usage:
  kubectl x shell [resource] [resource-name]

Args:
  resource       Resource type to open a shell on.
  resource-name  Name of the resource to open a shell on.

Example:
  kubectl-pi shell
  kubectl-pi shell pod
  kubectl-pi shell pod mypod`,
	Run: func(cmd *cobra.Command, args []string) {
		var resource string
		var resourceName string
		if len(args) > 0 {
			resource = args[0]
			if len(args) > 1 {
				resourceName = args[1]
			}
		}

		checkErr(shell.Shell(context.Background(), configFlags, resourceBuilderFlags, resource, resourceName, exactMatch))
	},
}

func init() {
	rootCmd.AddCommand(ctxCmd)
	ctxCmd.Flags().BoolVarP(&exactMatch, "exact", "e", false, "Exact match")
	resourceBuilderFlags.AddFlags(ctxCmd.Flags())
}
