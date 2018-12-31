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
  // "fmt"
  // "strings"
  "io/ioutil"
  "os"

  "github.com/spf13/cobra"


  // "github.com/starofservice/carbon/pkg/util/argparser"
  "github.com/starofservice/carbon/pkg/variables"
  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  // pkgmetalatest "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
  "github.com/starofservice/carbon/pkg/kubernetes"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/carbon/pkg/util/tojson"
)

var buildVars []string
var buildVarFiles []string
var buildPatches []string
var buildPatchFiles []string
var buildDefaultPWL bool

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
  Use:   "deploy",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  // Args: func(cmd *cobra.Command, args []string) {
  //   if len(args) != 1 {
  //     return fmt.Errorf("Currently only one argument is sapported")
  //   }
  //   if strings.Index(args[0], "://") != -1 {
  //     return fmt.Errorf("Image must not contain any schema like http:// or https://") 
  //   }
  //   if len(strings.Split(args[0], ":")) != 2 {
  //     return fmt.Errorf("Image must contain image name and tag delimited by a colon") 
  //   }
  // },
  Run: func(cmd *cobra.Command, args []string) {
    // fmt.Println("deploy called")
    runDeploy(args[0])
  },
}

func init() {
  RootCmd.AddCommand(deployCmd)

  deployCmd.Flags().StringArrayVar(&buildVars, "var", []string{}, "Define a value for a package variable")
  deployCmd.Flags().StringArrayVar(&buildVarFiles, "var-file", []string{}, "Define file with values for a package variables")
  deployCmd.Flags().StringArrayVar(&buildPatches, "patch", []string{}, "Apply directly typed patch for the manifest")
  deployCmd.Flags().StringArrayVar(&buildPatchFiles, "patch-file", []string{}, "Apply patch from a file for the manifest")
  deployCmd.Flags().BoolVar(&buildDefaultPWL, "default-prune-white-list", false, "Use the default prune white-list for the kubect apply operation. Enabling this option speeds-up deployment, but not all resource versions are pruned")

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
  patches := parsePatches()
  

  // var image string // should be replaced by input data
  // labels, err := dockermeta.GetLabels(image)
  log.Info("Getting carbon package metadata")
  dm := dockermeta.NewDockerMeta(image)
  labels := dm.GetLabels()
  // if err != nil {
  //   // panic(err.Error())
  //   log.Fatal("Failed to extract Carbon metadata from the Docker image '%s' due to the error: %s", image, err.Error())
  //   os.Exit(1)
  // }

  meta, err := pkgmeta.DeserializeMeta(labels)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to deserialize Carbon metadata from the Docker image '%s' due to the error: %s", image, err.Error())
    os.Exit(1)
  }

  kd, err := kubernetes.NewKubeDeployment(meta, dm.Name(), dm.Tag())
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to create new instance of KubeDeploy due to the error: %s", err.Error())
    os.Exit(1)
  }

  kd.UpdateVars(vars)

  log.Info("Building kubernetes configuration")
  kd.Build()
  // if err != nil {
  //   // panic(err.Error())
  //   log.Fatal("Failed to build kubernetes configuration due to the error: %s", err.Error())
  //   os.Exit(1)
  // }

  log.Info("Applying patches")
  // var patches [][]byte
  err = kd.ProcessPatches(patches)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to apply user patches for kubernetes resources due to the error: %s", err.Error())
    os.Exit(1)
  }



  err = kd.SetAppLabel()
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to apply Carbon labels for kubernetes resources due to the error: %s", err.Error())
    os.Exit(1) 
  }



  log.Info("Applying kubernetes configuration")
  // fmt.Println(string(kd.BuiltManifest))
  err = kd.Apply(buildDefaultPWL)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to apply Kubernetes configuration due to the error: %s", err.Error())
    os.Exit(1)
  }
  log.Info("Carbon package has been deployed successfully")
}

func parseVars() map[string]string {
  log.Debug("Parsing variables")

  vars := variables.NewVars()
  vars.ParseVarFiles(buildVarFiles)
  vars.ParseFlags(buildVars)
  
  return vars.Data
}

// func parsePatches() [][]byte {
//   log.Debug("Parsing patches")
//   // if len(buildPatches) > 0 && len(buildPatchFiles) > 0 {
//   //   fmt.Prinln("`patch` and `patch-file` parameters can't be used at the same time")
//   //   os.Exit(1)
//   // }
//   var patches [][]byte
//   for _, i := range buildPatchFiles {
//     d, err := ioutil.ReadFile(i)
//     if err != nil {
//       // panic(err.Error())
//       log.Fatalf("Failed to read a patch file '%s', due to the error: '%s'. Skipping it.", i, err.Error())
//       os.Exit(1)
//       // continue
//     }
//     patches = append(patches, d)
//   }

//   for _, i := range buildPatches {
//     patches = append(patches, []byte(i))
//   }

//   return patches
// }


func parsePatches() []byte {
  log.Debug("Parsing patches")
  // if len(buildPatches) > 0 && len(buildPatchFiles) > 0 {
  //   fmt.Prinln("`patch` and `patch-file` parameters can't be used at the same time")
  //   os.Exit(1)
  // }
  var rawPatches [][]byte
  for _, i := range buildPatchFiles {
    d, err := ioutil.ReadFile(i)
    if err != nil {
      // panic(err.Error())
      log.Fatalf("Failed to read a patch file '%s', due to the error: '%s'. Skipping it.", i, err.Error())
      os.Exit(1)
      // continue
    }
    rawPatches = append(rawPatches, d)
  }

  for _, i := range buildPatches {
    rawPatches = append(rawPatches, []byte(i))
  }

  var resp []byte
  for _, i := range rawPatches {
    ni := tojson.ToJson(i)
    resp = append(resp, ni...)
  }

  return resp
}