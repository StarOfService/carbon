package cmd

import (
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/kubemeta"
)

var uninstallCmd = &cobra.Command{
  Use:   "uninstall packageName [packageName ...]",
  Short: "Uninstall Carbon packages from your Kubernetes cluster",
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
      err := runUninstall(i)
      if err != nil {
        success = false
        log.Errorf("Failed to uninstall package '%s' due to the error: %s", i, err.Error())
      }
    }
    if success {
      log.Info("Carbon packages has been uninstalled successfully")
    } else {
      return errors.New("Carbon failed to uninstall some packages")
    }

    return nil
  },
}

func init() {
  RootCmd.AddCommand(uninstallCmd)
}

func runUninstall(pkg string) error {
  log.Info("Uninstalling Carbon package", pkg)

  installed, err := kubemeta.IsInstalled(pkg)
  if err != nil {
    return errors.Wrap(err, "checking if the package is installed")
  }
  if !installed {
    return errors.New("The package isn't installed")
  }

  kmeta, err := kubemeta.Get(pkg)
  if err != nil {
    return errors.Wrap(err, "getting package metadata")
  }

  err = kubernetes.Delete(kmeta.Data.Manifest, kmeta.Data.Namespace)
  if err != nil {
    return errors.Wrap(err, "uninstalling Kubernetes resources")
  }

  err = kmeta.Delete()
  if err != nil {
    log.Error("Failed to uninstall metadata for the package due to the error: ", err.Error())
    return errors.Wrap(err, "uninstalling Carbon package metadata")
  }

  log.Debugf("Carbon package '%s' has been uninstalled", pkg)
  return nil
}
