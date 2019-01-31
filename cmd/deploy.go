package cmd

import (
  "fmt"
  "io/ioutil"
  "os"

  "github.com/spf13/cobra"
  log "github.com/sirupsen/logrus"
  "github.com/pkg/errors"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/schema/kubemeta"
  "github.com/starofservice/carbon/pkg/variables"
  "github.com/starofservice/carbon/pkg/util/homedir"
  "github.com/starofservice/carbon/pkg/util/tojson"
)

var DeployVarFlags []string
var DeployVarFiles []string
var DeployPatches []string
var DeployPatchFiles []string
var DeployDefaultPWL bool
var DeployNamespace string
var DeployMetadataNamespace string

var deployCmd = &cobra.Command{
  Use:   "deploy image",
  Short: "Deploy Carbon package",
  Long: `
Deploy specified Carbon package to the currently active Kubernetes cluster.
The docker image with carbon metadata must be specified as a command argument.
It may be either local or remote docker image.`,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 1 {
      return fmt.Errorf("This command requires exactly one argument")
    }
    return nil
  },
  Run: func(cmd *cobra.Command, args []string) {
    runDeploy(args[0])
  },
}

func init() {
  RootCmd.AddCommand(deployCmd)

  deployCmd.Flags().StringArrayVar(&DeployVarFlags, "var", []string{}, "Define a value for a package variable")
  deployCmd.Flags().StringArrayVar(&DeployVarFiles, "var-file", []string{}, "Define file with values for a package variables")
  deployCmd.Flags().StringArrayVar(&DeployPatches, "patch", []string{}, "Apply directly typed patch for the manifest")
  deployCmd.Flags().StringArrayVar(&DeployPatchFiles, "patch-file", []string{}, "Apply patch from a file for the Kubernetes manifest")
  deployCmd.Flags().StringVarP(&DeployNamespace, "namespace", "n", "", "If present, defineds the Kubernetes namespace scope for the deployed resources and Carbon metadata")
  deployCmd.Flags().StringVar(&DeployMetadataNamespace, "metadata-namespace", "", "Namespace where Carbon has to keep its metadata. Current parameter has precendance over `namespace` and should be used for muli-namespaced environments")
  deployCmd.Flags().BoolVar(&DeployDefaultPWL, "default-prune-white-list", false, "Use the default prune white-list for the kubect apply operation. Enabling this option speeds-up deployment, but not all resource versions are pruned")
}

func runDeploy(image string) {
  log.Info("Starting Carbon deploy")

  if minikube.Enabled {
    err := minikube.SetDockerEnv()
    if err != nil {
      log.Fatal(err.Error())
      os.Exit(1)
    }
  }

  vars := parseVars()
  patches, err := parsePatches()
  if err != nil {
    log.Fatal("Failed to parse patches due to the error: ", err.Error())
    // os.Exit(1)
  }  

  log.Info("Getting carbon package metadata")
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

  kdeploy, err := kubernetes.NewKubeDeployment(meta, dm.Name(), dm.Tag())
  if err != nil {
    log.Fatal("Failed to create new instance of KubeDeploy due to the error: ", err.Error())
    // os.Exit(1)
  }

  kdeploy.UpdateVars(vars)

  log.Info("Building Kubernetes configuration")
  err = kdeploy.Build()
  if err != nil {
    log.Fatal("Failed to build Kubernetes configuration due to the error: ", err.Error())
    // os.Exit(1)
  }

  log.Info("Applying patches")
  err = kdeploy.ProcessPatches(patches)
  if err != nil {
    // TODO: Check amount of spaces between the custom text and the error message.
    log.Fatal("Failed to apply user patches for kubernetes resources due to the error: ", err.Error())
    // os.Exit(1)
  }

  err = kdeploy.SetAppLabel()
  if err != nil {
    log.Fatal("Failed to apply Carbon labels for kubernetes resources due to the error: ", err.Error())
    // os.Exit(1)
  }

  log.Info("Applying kubernetes configuration")
  err = kdeploy.Apply(DeployDefaultPWL, DeployNamespace)
  if err != nil {
    log.Error("Failed to apply Kubernetes configuration due to the error: ", err.Error())
    revert(kdeploy)
    os.Exit(1)
  }

  log.Info("Applying Carbon metadata for the package")
  kmeta := kubemeta.New(kdeploy, patches, getDeployMetadataNamespace())
  err = kmeta.Apply()
  if err != nil {
    log.Error("Failed to update Carbon metadata due to the error: ", err.Error())
    revert(kdeploy)
    os.Exit(1)
  }

  log.Info("Carbon package has been deployed successfully")
}

func parseVars() map[string]string {
  log.Debug("Parsing variables")

  homeVarsPath := homedir.Path() + "/" + "carbon.vars"

  vars := variables.NewVars()
  if _, err := os.Stat(homeVarsPath); err == nil {
   DeployVarFiles = append([]string{homeVarsPath}, DeployVarFiles...)
  }
  err := vars.ParseFiles(DeployVarFiles)
  if err != nil {
    log.Fatal("Failed to parse variable files due to the error: ", err.Error())
    // os.Exit(1)
  }
  vars.ParseFlags(DeployVarFlags)
  if err != nil {
    log.Fatal("Failed to parse variable flags due to the error: ", err.Error())
    // os.Exit(1)
  }

  return vars.Data
}

func revert(kdeploy *kubernetes.KubeDeployment) {
  log.Error("Trying to revert changes")
  deployed, err := kubemeta.IsDeployed(kdeploy.Variables.Pkg.Name, getDeployMetadataNamespace())
  if err != nil {
    log.Fatal("Failed to check if the package is already deployed due to the error: ", err.Error())
    // os.Exit(1)
  }
  if !deployed {
    log.Fatalf("Looks like the package '%s' has never been deployed yet. Nothing to do.", kdeploy.Variables.Pkg.Name)
  }

  kmeta, err := kubemeta.Get(kdeploy.Variables.Pkg.Name, getDeployMetadataNamespace())
  if err != nil {
    log.Fatal("Failed to revert Kubernetes configuration due to the error: ", err.Error())
    // os.Exit(1)
  }
  kdeploy.BuiltManifest = []byte(kmeta.Data.Manifest)
  err = kdeploy.Apply(DeployDefaultPWL, DeployNamespace)
  if err != nil {
    log.Fatal("Failed to revert Kubernetes configuration due to the error: ", err.Error())
    // os.Exit(1)
  }
  log.Error("Revert has ran successfully")
  os.Exit(1)
}

func parsePatches() ([]byte, error) {
  log.Debug("Parsing patches")

  var rawPatches [][]byte
  for _, i := range DeployPatchFiles {
    d, err := ioutil.ReadFile(i)
    if err != nil {
      return nil, errors.Wrapf(err, "reading a patch file '%s'", i)
    }
    rawPatches = append(rawPatches, d)
  }

  for _, i := range DeployPatches {
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

func getDeployMetadataNamespace() string {
  if DeployMetadataNamespace != "" {
    return DeployMetadataNamespace
  }
  if DeployNamespace != "" {
    return DeployNamespace
  }
  return "default"
}
