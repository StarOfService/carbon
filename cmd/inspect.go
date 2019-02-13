package cmd

import (
  "fmt"
  "os"

  "github.com/olekukonko/tablewriter"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
)

var inspectCmd = &cobra.Command{
  Use:   "inspect dockerImage [dockerImage ...]",
  Short: "Show information for a Carbon package",
  Long: `
Expose Carbon metadata for a given Docker image.
The image may be either local or remote.`,
  SilenceErrors: true,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
      return errors.New("This command requires at least one argument")
    }
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    cmd.SilenceUsage = true
    var success bool = true
    for _, i := range args {
      if err := runInspect(i); err != nil {
        success = false
        log.Warnf("Skipping image '%s' due to the error: %s", i, err.Error())
      }
    }
    if !success {
      return errors.New("Carbon failed to insepct some packages")
    }
    return nil
  },
}

func init() {
  RootCmd.AddCommand(inspectCmd)
}

func runInspect(image string) error {
  log.Debug("Getting carbon package metadata")
  dm, err := dockermeta.NewDockerMeta(image)
  if err != nil {
    return errors.Wrap(err, "creating new Docker metadata")
  }
  
  labels, err := dm.GetLabels()
  if err != nil {
    return errors.Wrapf(err, "getting Carbon metadata")
  }

  meta, err := pkgmeta.Deserialize(labels)
  if err != nil {
    return errors.Wrapf(err, "deserializing Carbon metadata")
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

  return nil
}
