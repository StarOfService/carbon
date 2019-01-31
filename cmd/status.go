package cmd

import (
  "fmt"
  "os"
  "github.com/olekukonko/tablewriter"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  "github.com/starofservice/carbon/pkg/schema/kubemeta"
)

var StatusNamespace string
var StatusFull bool
// var statusMetadataNamespace string

var statusCmd = &cobra.Command{
  Use:   "status [package_name [package_name ...]]",
  Short: "Show information about installed Carbon packages for your Kubernetes cluster",
  Long: `
Show information about installed Carbon packages at the current Kubernetes cluster.
Without arguments all installed packages are listed.
You can provide name for a specific package(s). In this case a detailed information for the requested package(s) is printed.`,
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
      runStatusAll()  
    } else {
      for n, i := range args {
        if n != 0 {
          fmt.Println("")
        }
        runStatusSingle(i)
      }
    } 
  },
}

func init() {
  RootCmd.AddCommand(statusCmd)

  statusCmd.Flags().StringVarP(&StatusNamespace, "namespace", "n", "", "If present, defineds the Kubernetes namespace scope for the deployed resources and Carbon metadata")
  statusCmd.Flags().BoolVarP(&StatusFull, "full", "f", false, "Print a full information for a given package(s) (including patches and applied manifests). Disabled by default.")
  // deployCmd.Flags().StringVar(&statusMetadataNamespace, "metadata-namespace", "", "Namespace where Carbon has to keep its metadata. Current parameter has precendance over `namespace` and should be used for muli-namespaced environments")
}

func runStatusAll() {
  meta, err := kubemeta.GetAll(getStatusMetadataNamespace())
  if err != nil {
    // panic("panic")
    log.Fatal("Failed to get status for the deployed carbon packages due to the error: ", err.Error())
  }

  if len(meta) == 0 {
    log.Info("Current Kubernetes context doesn't have any Carbon packages installed")
    return
  }

  table := tablewriter.NewWriter(os.Stdout)
  for _, i := range meta {
    table.Append([]string{i.Data.Name, i.Data.Version, i.Data.Source})
  }
  table.SetHeader([]string{"Name", "Version", "Source"})
  table.SetCenterSeparator("")
  table.SetColumnSeparator("")
  table.SetRowSeparator("")
  table.SetHeaderLine(false)
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetBorder(false)
  table.Render()
}

func runStatusSingle(pkg string) {
  meta, err := kubemeta.Get(pkg, getStatusMetadataNamespace())
  if err != nil {
    // panic("panic")
    log.Fatalf("Failed to get status for the carbon package '%s' due to the error: %s", pkg, err.Error())
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

  fmt.Println("Name:", meta.Data.Name)
  fmt.Println("Version:", meta.Data.Version)
  fmt.Println("Source:", meta.Data.Source)
  fmt.Println("Variables:")
  table.Render()
  if StatusFull {
    fmt.Println("Patches:", meta.Data.Patches)
    fmt.Println("Manifest:", meta.Data.Manifest)
  }
}

func getStatusMetadataNamespace() string {
  // if statusMetadataNamespace != "" {
  //   return statusMetadataNamespace
  // }
  if StatusNamespace != "" {
    return StatusNamespace
  }
  return "default"
}