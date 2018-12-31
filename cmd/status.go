package cmd

import (
  "fmt"
  "os"

  "github.com/olekukonko/tablewriter"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  "github.com/starofservice/carbon/pkg/schema/kubemeta"
)

var StatusFull bool

var statusCmd = &cobra.Command{
  Use:   "status [packageName [packageName ...]]",
  Short: "Show information about installed Carbon packages for your Kubernetes cluster",
  Long: `
Show information about installed Carbon packages at the current Kubernetes cluster.
Without arguments all installed packages are listed.
You can provide name for a specific package(s). In this case a detailed information for the requested package(s) is printed.`,
  SilenceErrors: true,
  RunE: func(cmd *cobra.Command, args []string) error {
    cmd.SilenceUsage = true
    if len(args) == 0 {
      return errors.Wrap(runStatusAll(), "status")
    } else {
      var success bool = true
      for n, i := range args {
        if n != 0 {
          fmt.Println("")
        }
        if err := runStatusSingle(i); err != nil {
          success = false
          log.Warnf("Failed to get status for the package '%s' due to the error: %s", i, err.Error())
        }
      }
      if !success {
        return errors.New("Carbon failed to get status for some packages")
      }
    }
    return nil
  },
}

func init() {
  RootCmd.AddCommand(statusCmd)

  statusCmd.Flags().BoolVarP(&StatusFull, "full", "f", false, "Print a full information for a given package(s) (including patches and applied manifests). Disabled by default.")
}

func runStatusAll() error {
  meta, err := kubemeta.GetAll()
  if err != nil {
    errors.Wrap(err, "getting information for all installed Carbon packages")
  }

  if len(meta) == 0 {
    log.Info("Current Kubernetes context doesn't have any Carbon packages installed")
    return nil
  }

  table := tablewriter.NewWriter(os.Stdout)
  for _, i := range meta {
    table.Append([]string{i.Data.Namespace, i.Data.Name, i.Data.Version, i.Data.Source})
  }
  table.SetHeader([]string{"Namespace", "Name", "Version", "Source"})
  table.SetCenterSeparator("")
  table.SetColumnSeparator("")
  table.SetRowSeparator("")
  table.SetHeaderLine(false)
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetBorder(false)
  table.Render()

  return nil
}

func runStatusSingle(pkg string) error {
  installed, err := kubemeta.IsInstalled(pkg)
  if err != nil {
    return errors.Wrap(err, "checking if the package is installed")
  }
  if !installed {
    return errors.New("The package isn't installed")
  }

  meta, err := kubemeta.Get(pkg)
  if err != nil {
    errors.Wrap(err, "getting information for the Carbon packages")
  }

  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader([]string{"Name", "Value"})
  table.SetCenterSeparator("")
  table.SetColumnSeparator("")
  table.SetRowSeparator("")
  table.SetHeaderLine(false)
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetBorder(false)

  for k, v := range meta.Data.Variables {
    table.Append([]string{k, v})
  }

  fmt.Println("Namespace:", meta.Data.Namespace)
  fmt.Println("Name:", meta.Data.Name)
  fmt.Println("Version:", meta.Data.Version)
  fmt.Println("Source:", meta.Data.Source)
  fmt.Println("Variables:")
  table.Render()
  if StatusFull {
    fmt.Println("Patches:", meta.Data.Patches)
    fmt.Println("Manifest:", meta.Data.Manifest)
  }

  return nil
}
