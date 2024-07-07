package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	configFlags = genericclioptions.NewConfigFlags(true).WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	rootCmd     = &cobra.Command{
		Use: "kubectl-x",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl x",
		},
		Short: "kubectl (kube-control) plugin with various useful extensions.",
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(cmd.Help())
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error:"), err)
		os.Exit(1)
	}
}
