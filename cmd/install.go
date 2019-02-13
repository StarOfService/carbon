package cmd

import (
  "io/ioutil"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/homecfg"
  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/kubemeta"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/util/tojson"
  "github.com/starofservice/carbon/pkg/variables"
)

var InstallVarFlags []string
var InstallVarFiles []string
var InstallPatches []string
var InstallPatchFiles []string
var InstallDefaultPWL bool

var installCmd = &cobra.Command{
  Use:   "install image",
  Short: "Install Carbon package",
  Long: `
Install specified Carbon package to the currently active Kubernetes cluster.
The Docker image with Carbon metadata must be specified as a command argument.
It may be either local or remote Docker image.`,
  SilenceErrors: true,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 1 {
      return errors.New("This command requires exactly one argument")
    }
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    cmd.SilenceUsage = true
    return errors.Wrap(runInstall(args[0]), "install")
  },
}

func init() {
  RootCmd.AddCommand(installCmd)

  installCmd.Flags().StringArrayVar(&InstallVarFlags, "var", []string{}, "Define a value for a package variable")
  installCmd.Flags().StringArrayVar(&InstallVarFiles, "var-file", []string{}, "Define file with values for a package variables")
  installCmd.Flags().StringArrayVar(&InstallPatches, "patch", []string{}, "Apply directly typed patch for the manifest")
  installCmd.Flags().StringArrayVar(&InstallPatchFiles, "patch-file", []string{}, "Apply patch from a file for the Kubernetes manifest")
  installCmd.Flags().BoolVar(&InstallDefaultPWL, "default-prune-white-list", false, "Use the default prune white-list for the kubect apply operation. Enabling this option speeds-up deployment, but not all resource versions are pruned")
}

func runInstall(image string) error {
  log.Info("Starting Carbon install")

  vars, err := parseVars()
  if err != nil {
    return errors.Wrap(err, "parsing variables")
  }
  patches, err := parsePatches()
  if err != nil {
    return errors.Wrap(err, "parsing patches")
  }

  log.Info("Getting Carbon package metadata")
  dm, err := dockermeta.NewDockerMeta(image)
  if err != nil {
    return errors.Wrap(err, "creating new Docker metadata")
  }

  labels, err := dm.GetLabels()
  if err != nil {
    return errors.Wrap(err, "extract Carbon metadata from the Docker image")
  }

  pmeta, err := pkgmeta.Deserialize(labels)
  if err != nil {
    return errors.Wrap(err, "deserializing Carbon metadata from the Docker image")
  }

  kinstall, err := kubernetes.NewKubeInstall(pmeta, dm.Name(), dm.Tag())
  if err != nil {
    return errors.Wrap(err, "creating new KubeInstall instance")
  }

  kinstall.UpdateVars(vars)

  log.Info("Building Kubernetes configuration")
  if err = kinstall.Build(); err != nil {
    return errors.Wrap(err, "building Kubernetes configuration")
  }

  log.Info("Applying patches")
  if err = kinstall.ProcessPatches(patches); err != nil {
    return errors.Wrap(err, "applying provided patches for Kubernetes resources")
  }

  if err = kinstall.SetAppLabel(); err != nil {
    return errors.Wrap(err, "applying Carbon labels for Kubernetes resources")
  }

  log.Info("Applying kubernetes configuration")
  if err = kinstall.Apply(InstallDefaultPWL); err != nil {
    log.Error("Failed to apply Kubernetes configuration due to the error: ", err.Error())
    return revertInstall(kinstall)
  }

  log.Info("Saving Carbon package metadata")
  kmeta, err := kubemeta.New(kinstall, patches)
  if err != nil {
    return errors.Wrap(err, "creating package metadata for Kubernetes")
  }

  if err = kmeta.Apply(); err != nil {
    log.Error("Failed to save Carbon packge metadata due to the error: ", err.Error())
    return revertInstall(kinstall)
  }

  log.Info("Carbon package has been installed successfully")

  return nil
}

func parseVars() (map[string]string, error) {
  log.Debug("Parsing variables")

  vars := variables.NewVars()

  InstallVarFiles = append([]string{homecfg.HomeConfigVarfilePath()}, InstallVarFiles...)
  if err := vars.ParseFiles(InstallVarFiles); err != nil {
    return vars.Data, errors.Wrap(err, "parsing variable files")
  }

  if err := vars.ParseFlags(InstallVarFlags); err != nil {
    return vars.Data, errors.Wrap(err, "parsing variable flags")
  }

  return vars.Data, nil
}

func revertInstall(kinstall *kubernetes.KubeInstall) error {
  log.Error("Trying to revert changes")

  installed, err := kubemeta.IsInstalled(kinstall.Variables.Pkg.Name)
  if err != nil {
    return errors.Wrap(err, "checking if a preveous version of the package is installed")
  }
  if !installed {
    return errors.Errorf("The package '%s' has never been installed yet. Nothing to do.", kinstall.Variables.Pkg.Name)
  }

  kmeta, err := kubemeta.Get(kinstall.Variables.Pkg.Name)
  if err != nil {
    return errors.Wrap(err, "getting previous Kubernetes configuration")
  }

  kinstall.BuiltManifest = []byte(kmeta.Data.Manifest)

  if err = kinstall.Apply(InstallDefaultPWL); err != nil {
    return errors.Wrap(err, "applying previous Kubernetes configuration")
  }

  return errors.New("Revert has been applied successfully")
}

func parsePatches() ([]byte, error) {
  log.Debug("Parsing patches")

  var rawPatches [][]byte
  for _, i := range InstallPatchFiles {
    d, err := ioutil.ReadFile(i)
    if err != nil {
      return nil, errors.Wrapf(err, "reading a patch file '%s'", i)
    }
    rawPatches = append(rawPatches, d)
  }

  for _, i := range InstallPatches {
    rawPatches = append(rawPatches, []byte(i))
  }

  var resp []byte
  for _, i := range rawPatches {
    ni, err := tojson.ToJSON(i)
    if err != nil {
      return nil, err
    }
    resp = append(resp, ni...)
  }

  return resp, nil
}
