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
  "os"

  // homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  log "github.com/sirupsen/logrus"
  // "github.com/spf13/viper"

  // "github.com/starofservice/carbon/cmd/deploy"
  // "github.com/starofservice/carbon/cmd/status"
)

// var cfgFile string
var logLevel string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
  Use:   "carbon",
  Short: "A brief description of your application",
  Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //  Run: func(cmd *cobra.Command, args []string) { },
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    // fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
    log.SetFormatter(&log.TextFormatter{
      FullTimestamp: true,
    })
    switch logLevel {
    case "trace":
      log.SetLevel(log.TraceLevel)
    case "debug":
      log.SetLevel(log.DebugLevel)
    case "info":
      log.SetLevel(log.InfoLevel)
    case "warn":
      log.SetLevel(log.WarnLevel)
    case "error":
      log.SetLevel(log.ErrorLevel)
    case "fatal":
      log.SetLevel(log.FatalLevel)
    default:
      log.Fatal("Unsupported log level: %s", logLevel)
      os.Exit(1)
    }
  },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set the logging level ('trace'|'debug'|'info'|'warn'|'error'|'fatal') (default 'info')")


  // cobra.OnInitialize(initConfig)

  // // Here you will define your flags and configuration settings.
  // // Cobra supports persistent flags, which, if defined here,
  // // will be global for your application.
  // RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.carbon.yaml)")

  // // Cobra also supports local flags, which will only run
  // // when this action is called directly.
  // RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// // initConfig reads in config file and ENV variables if set.
// func initConfig() {
//  if cfgFile != "" {
//    // Use config file from the flag.
//    viper.SetConfigFile(cfgFile)
//  } else {
//    // Find home directory.
//    home, err := homedir.Dir()
//    if err != nil {
//      fmt.Println(home)
//      os.Exit(1)
//    }

//    // Search config in home directory with name ".cobra" (without extension).
//    viper.AddConfigPath(home)
//    viper.SetConfigName(".cobra")
//  }

//  viper.AutomaticEnv() // read in environment variables that match

//  // If a config file is found, read it in.
//  if err := viper.ReadInConfig(); err == nil {
//    fmt.Println("Using config file:", viper.ConfigFileUsed())
//  }
// }
