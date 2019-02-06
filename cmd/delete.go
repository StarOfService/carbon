package cmd

import (
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/kubemeta"
)

var DeleteMetadataNamespace string

var deleteCmd = &cobra.Command{
  Use:   "delete packageName [packageName ...]",
  Short: "Delete Carbon packages from your Kubernetes cluster",
  Long: `
`,
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
      err := runDelete(i)
      if err != nil {
        success = false
        log.Errorf("Failed to delete package '%s' due to the error: %s", i, err.Error())
      }
    }
    if success {
      log.Info("Carbon packages has been deleted successfully")  
    } else {
      return errors.New("Carbon failed to delete some packages")
    }

    return nil
  },
}

func init() {
  RootCmd.AddCommand(deleteCmd)

  deleteCmd.Flags().StringVar(&DeleteMetadataNamespace, "metadata-namespace", "", "Namespace where Carbon has to keep its metadata. Current parameter has precendance over `namespace` and should be used for muli-namespaced environments")
}

func runDelete(pkg string) error {
  log.Info("Deleting Carbon package", pkg)

  installed, err := kubemeta.IsInstalled(
    pkg,
    MetadataNamespace(DeleteMetadataNamespace, ""),
  )
  if err != nil {
    return errors.Wrap(err, "checking if the package is installed")
  }
  if !installed {
    return errors.New("The package isn't installed")
  }

  kmeta, err := kubemeta.Get(
    pkg,
    MetadataNamespace(DeleteMetadataNamespace, ""),
  )
  if err != nil {
    return errors.Wrap(err, "getting package metadata")
  }

  err = kubernetes.Delete(kmeta.Data.Manifest, kmeta.Data.Namespace)
  if err != nil {
    return errors.Wrap(err, "deleting Kubernetes resources")
  }

  err = kubemeta.Delete(pkg, MetadataNamespace(DeleteMetadataNamespace, ""))
  if err != nil {
    log.Error("Failed to delete metadata for the package due to the error: ", err.Error())
    return errors.Wrap(err, "deleting Carbon package metadata")
  }

  log.Debugf("Carbon package '%s' has been deleted", pkg)
  return nil
}
