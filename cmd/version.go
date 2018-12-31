package cmd

import (
  "fmt"

  "github.com/spf13/cobra"

  carboncfglatest "github.com/starofservice/carbon/pkg/schema/carboncfg/latest"
  pkgcfglatest "github.com/starofservice/carbon/pkg/schema/pkgcfg/latest"
  "github.com/starofservice/carbon/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the current Carbon versions",
  Long: `
Print the current Carbon versions`,
  Run: func(cmd *cobra.Command, args []string) {
    // version.PrintVersion()
    runVersion()
  },
}

func init() {
  RootCmd.AddCommand(versionCmd)
}

func runVersion() {
  fmt.Println("Carbon version:", version.GetVersion())
  fmt.Println("Latest Carbon config apiVersion:", carboncfglatest.Version)
  fmt.Println("Latest package config apiVersion:", pkgcfglatest.Version)
}