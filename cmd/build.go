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
  // "reflect"
  "github.com/spf13/cobra"
  "io/ioutil"
  "os"

  // // "github.com/spf13/viper"
  "path/filepath"

  // // "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
  "github.com/starofservice/carbon/pkg/schema/rootcfg"
  dockerbuild "github.com/starofservice/carbon/pkg/docker/build"
  // "github.com/starofservice/carbon/pkg/kubernetes/manifest"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/kubernetes"
  // "github.com/starofservice/carbon/pkg/util"
  log "github.com/sirupsen/logrus"
)

var cfgFile string

// buildCmd represents the build command
var buildCmd = &cobra.Command{
  Use:   "build",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  // Run: func(cmd *cobra.Command, args []string) {
  //   fmt.Println("build called")
  // },
  Run: func(cmd *cobra.Command, args []string) {
    runBuild()
  },
}

func init() {
  RootCmd.AddCommand(buildCmd)

  // cobra.OnInitialize(initConfig)

  buildCmd.Flags().StringVarP(&cfgFile, "config", "c", "carbon.yaml", "config file (default is carbon.yaml)")
  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // buildCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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


  kubeManif, err := kubernetes.ReadTemplates(cfg.KubeManifests)
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

  err = kd.VerifyAll(cfg.KubeManifests)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to verify Kubernetes configs due to the error: %s", err.Error())
    // fmt.Println(err.Error())
    os.Exit(1)

  }

  metaMap, err := pkgmeta.SerializeMeta(*meta)
  if err != nil {
    // panic(err.Error())
    log.Fatalf("Failed to serialize Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

  // fmt.Println("cfg processed")
  // fmt.Println(cfg.Dockerfile)
  // fmt.Println(filepath.Dir(cfgPath))
  // fmt.Println(reflect.TypeOf(cfg.Artifacts))
  // fmt.Println(cfg.Artifacts)

  // fmt.Println(cfg)
  // fmt.Println(metaMap)

  log.Info("Building Carbon package")
  bo := dockerbuild.NewBuildOptions()
  bo.Build(cfg, filepath.Dir(cfgPath), metaMap)
  // fmt.Println("docker build processed")
  log.Info("Carbon package has been built successfully")
}