// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
  "fmt"
  // "strings"
  "io/ioutil"
  "os"

  "github.com/spf13/cobra"
  "github.com/pkg/errors"

  // "github.com/starofservice/carbon/pkg/util/argparser"
  "github.com/starofservice/carbon/pkg/variables"
  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/schema/kubemeta"
  // pkgmetalatest "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
  "github.com/starofservice/carbon/pkg/kubernetes"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/carbon/pkg/util/tojson"
  "github.com/starofservice/carbon/pkg/util/homedir"
)

var deployVarFlags []string
var deployVarFiles []string
var deployPatches []string
var deployPatchFiles []string
var deployDefaultPWL bool
var deployNamespace string
var deployMetadataNamespace string

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
  Use:   "deploy",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 1 {
      return fmt.Errorf("This command requires exactly one argument")
    }
    return nil
  },
  Run: func(cmd *cobra.Command, args []string) {
    // fmt.Println("deploy called")

    // TODO check number of arguments
    runDeploy(args[0])
  },
}

func init() {
  RootCmd.AddCommand(deployCmd)

  deployCmd.Flags().StringArrayVar(&deployVarFlags, "var", []string{}, "Define a value for a package variable")
  deployCmd.Flags().StringArrayVar(&deployVarFiles, "var-file", []string{}, "Define file with values for a package variables")
  deployCmd.Flags().StringArrayVar(&deployPatches, "patch", []string{}, "Apply directly typed patch for the manifest")
  deployCmd.Flags().StringArrayVar(&deployPatchFiles, "patch-file", []string{}, "Apply patch from a file for the Kubernetes manifest")
  deployCmd.Flags().StringVarP(&deployNamespace, "namespace", "n", "", "If present, defineds the Kubernetes namespace scope for the deployed resources and Carbon metadata")
  deployCmd.Flags().StringVar(&deployMetadataNamespace, "metadata-namespace", "", "Namespace where Carbon has to keep its metadata. Current parameter has precendance over `namespace` and should be used for muli-namespaced environments")
  deployCmd.Flags().BoolVar(&deployDefaultPWL, "default-prune-white-list", false, "Use the default prune white-list for the kubect apply operation. Enabling this option speeds-up deployment, but not all resource versions are pruned")

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // deployCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runDeploy(image string) {
  log.Info("Starting Carbon deploy")

  vars := parseVars()
  patches, err := parsePatches()
  if err != nil {
    log.Fatal("Failed to parse patches due to the error: %s", err.Error())
    os.Exit(1)
  }  

  // var image string // should be replaced by input data
  // labels, err := dockermeta.GetLabels(image)
  log.Info("Getting carbon package metadata")
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

  kdeploy, err := kubernetes.NewKubeDeployment(meta, dm.Name(), dm.Tag())
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to create new instance of KubeDeploy due to the error: %s", err.Error())
    os.Exit(1)
  }

  kdeploy.UpdateVars(vars)

  log.Info("Building Kubernetes configuration")
  err = kdeploy.Build()
  if err != nil {
    // panic(err.Error())
    log.Fatal("Failed to build Kubernetes configuration due to the error: %s", err.Error())
    os.Exit(1)
  }

  log.Info("Applying patches")
  // var patches [][]byte
  err = kdeploy.ProcessPatches(patches)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to apply user patches for kubernetes resources due to the error: %s", err.Error())
    os.Exit(1)
  }



  err = kdeploy.SetAppLabel()
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to apply Carbon labels for kubernetes resources due to the error: %s", err.Error())
    os.Exit(1) 
  }



  log.Info("Applying kubernetes configuration")
  // fmt.Println(string(kdeploy.BuiltManifest))
  err = kdeploy.Apply(deployDefaultPWL, deployNamespace)
  if err != nil {
    // panic(err.Error())
    log.Errorf("Failed to apply Kubernetes configuration due to the error: %s", err.Error())
    revert(kdeploy)
    os.Exit(1)
  }

  log.Info("Applying Carbon metadata for the package")
  kmeta := kubemeta.New(kdeploy, patches, getDeployMetadataNamespace())
  err = kmeta.Apply()
  if err != nil {
    // panic(err.Error())
    log.Errorf("Failed to update Carbon metadata due to the error: %s", err.Error())
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
    // vars.ParseFiles([]string{homeVarsPath})
   deployVarFiles = append([]string{homeVarsPath}, deployVarFiles...)
  }
  err := vars.ParseFiles(deployVarFiles)
  if err != nil {
    log.Fatalf("Failed to parse variable files due to the error: %s", err.Error())
    os.Exit(1)
  }
  vars.ParseFlags(deployVarFlags)
  if err != nil {
    log.Fatalf("Failed to parse variable flags due to the error: %s", err.Error())
    os.Exit(1)
  }

  return vars.Data
}

func revert(kdeploy *kubernetes.KubeDeployment) {
  log.Error("Trying to revert changes")
  kmeta, err := kubemeta.Get(kdeploy.Variables.Pkg.Name, getDeployMetadataNamespace())
  if err != nil {
    log.Fatalf("Failed to revert Kubernetes configuration due to the error: %s", err.Error())
    os.Exit(1)
  }
  kdeploy.BuiltManifest = []byte(kmeta.Data.Manifest)
  err = kdeploy.Apply(deployDefaultPWL, deployNamespace)
  if err != nil {
    log.Fatalf("Failed to revert Kubernetes configuration due to the error: %s", err.Error())
    os.Exit(1)
  }
  log.Error("Revert has ran successfully")
  os.Exit(1)
}

// func parsePatches() [][]byte {
//   log.Debug("Parsing patches")
//   // if len(deployPatches) > 0 && len(deployPatchFiles) > 0 {
//   //   fmt.Prinln("`patch` and `patch-file` parameters can't be used at the same time")
//   //   os.Exit(1)
//   // }
//   var patches [][]byte
//   for _, i := range deployPatchFiles {
//     d, err := ioutil.ReadFile(i)
//     if err != nil {
//       // panic(err.Error())
//       log.Fatalf("Failed to read a patch file '%s', due to the error: '%s'. Skipping it.", i, err.Error())
//       os.Exit(1)
//       // continue
//     }
//     patches = append(patches, d)
//   }

//   for _, i := range deployPatches {
//     patches = append(patches, []byte(i))
//   }

//   return patches
// }


func parsePatches() ([]byte, error) {
  log.Debug("Parsing patches")
  // if len(deployPatches) > 0 && len(deployPatchFiles) > 0 {
  //   fmt.Prinln("`patch` and `patch-file` parameters can't be used at the same time")
  //   os.Exit(1)
  // }
  var rawPatches [][]byte
  for _, i := range deployPatchFiles {
    d, err := ioutil.ReadFile(i)
    if err != nil {
      // panic(err.Error())
      return nil, errors.Wrapf(err, "reading a patch file '%s'", i)
      // log.Fatalf("Failed to read a patch file '%s', due to the error: '%s'. Skipping it.", i, err.Error())
      // os.Exit(1)
      // continue
    }
    rawPatches = append(rawPatches, d)
  }

  for _, i := range deployPatches {
    rawPatches = append(rawPatches, []byte(i))
  }

  var resp []byte
  for _, i := range rawPatches {
    ni, err := tojson.ToJson(i)
    if err != nil {
      // panic(err.Error())
      return nil, err
      // continue
    }
    resp = append(resp, ni...)
  }

  return resp, nil
}

func getDeployMetadataNamespace() string {
  if deployMetadataNamespace != "" {
    return deployMetadataNamespace
  }
  if deployNamespace != "" {
    return deployNamespace
  }
  return "default"
}
