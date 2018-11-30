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

  "github.com/spf13/cobra"

  // dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  // "github.com/starofservice/carbon/pkg/schema/pkgmeta"
)

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




  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // deployCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runDeploy(image string) {
  // var image string // should be replaced by input data
  // labels, err := dockermeta.GetLabels(image)
  // meta, err := pkgmeta.ParseMetadata(labels)
}