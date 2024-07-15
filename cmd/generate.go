package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/RRethy/kubectl-x/internal/cmd/generate"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Use ChatGPT to generate a shell command to run.",
	Long: `Use ChatGPT to generate a shell command to run.

Usage:
  kubectl x gen [description]

Args:
  description  Description of the command to generate.

Example:
  kubectl x gen list all pods sorted by date
  kubectl x gen "list all pods sorted by date"
  kubectl x gen scale deployment web to 8 relicas`,
	Run: func(cmd *cobra.Command, args []string) {
		checkErr(generate.Generate(context.Background(), strings.Join(args, " ")))
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
