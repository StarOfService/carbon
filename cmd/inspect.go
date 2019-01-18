package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
  // "io/ioutil"
  "os"
  log "github.com/sirupsen/logrus"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  // "github.com/starofservice/carbon/pkg/kubernetes/manifest"
  "github.com/olekukonko/tablewriter"
)

var inspectCmd = &cobra.Command{
  Use:   "inspect",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
      return fmt.Errorf("This command requires at least one argument")
    }
    return nil
  },
  Run: func(cmd *cobra.Command, args []string) {
    for _, i := range args {
      runInspect(i)  
    }
  },
}

func init() {
  RootCmd.AddCommand(inspectCmd)
}

func runInspect(image string) {
  log.Debug("Getting carbon package metadata")
  dm := dockermeta.NewDockerMeta(image)
  labels, err := dm.GetLabels()
  if err != nil {
    log.Fatal("Failed to extract Carbon metadata from the Docker image '%s' due to the error: %s", image, err.Error())
    os.Exit(1)
  }

  meta, err := pkgmeta.Deserialize(labels)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to deserialize Carbon metadata from the Docker image '%s' due to the error: %s", image, err.Error())
    os.Exit(1)
  }

  // var varTable [][]string
  table := tablewriter.NewWriter(os.Stdout)
  for _, vh := range meta.Data.Variables {
    table.Append([]string{vh.Name, vh.Default, vh.Description})
    // vrow := []string{vh.Name, vh.Default, vh.Description}
    // varTable = append(varTable, vrow)
  }
  table.SetHeader([]string{"Name", "Default", "Description"})
  table.SetCenterSeparator("")
  table.SetColumnSeparator("")
  table.SetRowSeparator("")
  table.SetHeaderLine(false)
  table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
  table.SetBorder(false)

  fmt.Println("Name:", meta.Data.PkgName)
  fmt.Println("Version:", meta.Data.PkgVersion)
  fmt.Println("Variables:")
  table.Render()
}