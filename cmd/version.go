package cmd

import (
  "github.com/spf13/cobra"

  "github.com/starofservice/carbon/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the current Carbon versions",
  Long: `
Print the current Carbon versions`,
  Run: func(cmd *cobra.Command, args []string) {
    version.PrintVersion()
  },
}

func init() {
  RootCmd.AddCommand(versionCmd)
}
