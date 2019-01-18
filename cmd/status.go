package cmd

import (
  "fmt"
  "github.com/spf13/cobra"

  "os"
  log "github.com/sirupsen/logrus"


  // "github.com/starofservice/carbon/pkg/kubernetes/manifest"
  "github.com/olekukonko/tablewriter"
  "github.com/starofservice/carbon/pkg/schema/kubemeta"
)

var statusNamespace string
var statusFull bool
// var statusMetadataNamespace string

var statusCmd = &cobra.Command{
  Use:   "status",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

  statusCmd.Flags().StringVarP(&statusNamespace, "namespace", "n", "", "If present, defineds the Kubernetes namespace scope for the deployed resources and Carbon metadata")
  statusCmd.Flags().BoolVarP(&statusFull, "full", "f", false, "Print full information (including patches and manifests) for a given package. Disabled by default.")
  // deployCmd.Flags().StringVar(&statusMetadataNamespace, "metadata-namespace", "", "Namespace where Carbon has to keep its metadata. Current parameter has precendance over `namespace` and should be used for muli-namespaced environments")
}

func runStatusAll() {
  meta, err := kubemeta.GetAll(getStatusMetadataNamespace())
  if err != nil {
    panic("panic")
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
    panic("panic")
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
  if statusFull {
    fmt.Println("Patches:", meta.Data.Patches)
    fmt.Println("Manifest:", meta.Data.Manifest)
  }
}

func getStatusMetadataNamespace() string {
  // if statusMetadataNamespace != "" {
  //   return statusMetadataNamespace
  // }
  if statusNamespace != "" {
    return statusNamespace
  }
  return "default"
}