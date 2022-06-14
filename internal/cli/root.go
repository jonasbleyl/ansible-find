package cli

import (
	"github.com/spf13/cobra"
)

var vaultFile string

func Setup() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "ansible-vars",
		Long: `A CLI tool to help manage ansible variables.

This tool will only use variable YAML files that reside within one of
the following directories: [group_vars, host_vars, defaults, vars]`,
	}

	rootCmd.PersistentFlags().StringVarP(&vaultFile, "vault", "v", ".vault", "ansible vault password file")

	setupFind(rootCmd)
	return rootCmd
}
