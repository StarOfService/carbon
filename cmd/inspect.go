package cmd

import (
  "fmt"
  "os"

  "github.com/olekukonko/tablewriter"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
)

var inspectCmd = &cobra.Command{
  Use:   "inspect docker_image [docker_image ...]",
  Short: "Show information for a Carbon package",
  Long: `
Expose Carbon metadata for a given Docker image.
The image may be either local or remote.`,
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
    log.Fatalf("Failed to extract Carbon metadata from the Docker image '%s' due to the error: %s", image, err.Error())
    // os.Exit(1)
  }

  meta, err := pkgmeta.Deserialize(labels)
  if err != nil {
    log.Fatalf("Failed to deserialize Carbon metadata from the Docker image '%s' due to the error: %s", image, err.Error())
    // os.Exit(1)
  }

  table := tablewriter.NewWriter(os.Stdout)
  for _, vh := range meta.Data.Variables {
    table.Append([]string{vh.Name, vh.Default, vh.Description})
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
