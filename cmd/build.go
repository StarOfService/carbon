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
  // "reflect"
  "github.com/spf13/cobra"
  "io/ioutil"
  "os"

  // // "github.com/spf13/viper"
  "path/filepath"

  // // "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
  "github.com/starofservice/carbon/pkg/schema/rootcfg"
  dockerbuild "github.com/starofservice/carbon/pkg/docker/build"
  // dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  // "github.com/starofservice/carbon/pkg/kubernetes/manifest"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/kubernetes"
  // "github.com/starofservice/carbon/pkg/util"
  log "github.com/sirupsen/logrus"
)

var cfgFile string
// var buildPush bool

// buildCmd represents the build command
var buildCmd = &cobra.Command{
  Use:   "build",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 0 {
      return fmt.Errorf("This command doesn't use any arguments")
    }
    return nil
  },
  Run: func(cmd *cobra.Command, args []string) {
    runBuild()
  },
}

func init() {
  RootCmd.AddCommand(buildCmd)

  // cobra.OnInitialize(initConfig)

  buildCmd.Flags().StringVarP(&cfgFile, "config", "c", "carbon.yaml", "config file (default is carbon.yaml)")
  // buildCmd.Flags().BoolVar(&buildPush, "push", false, "Push built images to the repositories (disabled by default)")
}

// func initConfig() {
//  if cfgFile != "" {
//    viper.SetConfigFile("config.yaml")
//  } else {
//    // Use config file from the flag.
//    viper.SetConfigFile(cfgFile)

//    // // Find home directory.
//    // home, err := homedir.Dir()
//    // if err != nil {
//    //  fmt.Println(home)
//    //  os.Exit(1)
//    // }

//    // // Search config in home directory with name ".cobra" (without extension).
//    // viper.AddConfigPath(home)
//    // viper.SetConfigName(".cobra")
//  }

//  viper.AutomaticEnv() // read in environment variables that match

//  // If a config file is found, read it in.
//  if err := viper.ReadInConfig(); err == nil {
//    fmt.Println("Using config file:", viper.ConfigFileUsed())
//  }
// }

func runBuild() {
  log.Info("Starting Carbon build")

  log.Info("Reading Carbon config")
  cfgPath, err := filepath.Abs(cfgFile)
  if err != nil {
    log.Fatalf("Failed to find Carbon config due to the error: %s", err.Error())
    os.Exit(1)
    // panic(err.Error())
  }
  // fmt.Println(cfgPath)
  cfgBody, err := ioutil.ReadFile(cfgPath)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to read Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

  cfg, err := rootcfg.ParseConfig(cfgBody)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to parse Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }
  // cfgB64 := util.EncodeMetadata(cfgBody)
  // cfgB64 := pkgmeta.B64Encode(cfgBody)

  if cfg.HookDefined(rootcfg.HookPreBuild) {
    log.Info("Running pre-build hook")
    err = cfg.RunHook(rootcfg.HookPreBuild)
    if err != nil {
      log.Fatalf("Failed to run pre-build hook due to the error: %s", err.Error())
      os.Exit(1)
    }    
  }



  kubeManif, err := kubernetes.ReadTemplates(cfg.Data.KubeManifests)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to read Kubernetes configs due to the error: %s", err.Error())
    os.Exit(1)
  }
  // kubeManifB64 := pkgmeta.B64Encode(kubeManif)

  meta := pkgmeta.New(cfg, cfgBody, kubeManif)

  // fmt.Println(meta.Variables)

  kd, err := kubernetes.NewKubeDeployment(meta, "image", "tag")
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to create new instance of KubeDeploy due to the error: %s", err.Error())
    os.Exit(1)
  }

  err = kd.VerifyAll(cfg.Data.KubeManifests)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to verify Kubernetes configs due to the error: %s", err.Error())
    // fmt.Println(err.Error())
    os.Exit(1)

  }

  metaMap, err := meta.Serialize()
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to serialize Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

  log.Info("Building Carbon package")
  buildOpts := dockerbuild.NewBuildOptions()
  err = buildOpts.Build(cfg, filepath.Dir(cfgPath), metaMap)
  if err != nil {
    log.Fatalf("Failed to build Carbon package due to the error: %s", err.Error())
    os.Exit(1)
  }
  // fmt.Println("docker build processed")

  if cfg.HookDefined(rootcfg.HookPostBuild) {
    log.Info("Running post-build hook")
    err = cfg.RunHook(rootcfg.HookPostBuild)
    if err != nil {
      log.Fatalf("Failed to run post-build hook due to the error: %s", err.Error())
      os.Exit(1)
    }
  }

  // if buildPush {
  //   log.Info("Pushing built docker images")
  //   // TODO
  //   var tagsToPush []string
  //   for _, i := range cfg.Data.Artifacts {
  //     dm := dockermeta.NewDockerMeta(i)
  //     // fmt.Println(dm.Domain())
  //     log.Warnf("domain: %s", dm.Domain())
  //     if dm.Domain() != "" {
  //       tagsToPush = append(tagsToPush, i)
  //     }
  //   }
  //   log.Warnf("domains length: %s", len(tagsToPush))
  // }

  log.Info("Carbon package has been built successfully")
}