package cmd

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  dockerbuild "github.com/starofservice/carbon/pkg/docker/build"
  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/schema/rootcfg"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
)

var cfgFile string
var buildPush bool

var buildCmd = &cobra.Command{
  Use:   "build",
  Short: "Build Carbon package",
  Long: `
Builds a Carbon package based on the provided carbon.yaml config.`,
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

  buildCmd.Flags().StringVarP(&cfgFile, "config", "c", "carbon.yaml", "config file (default is carbon.yaml)")
  buildCmd.Flags().BoolVar(&buildPush, "push", false, "Push built images to the repositories (disabled by default)")
}

func runBuild() {
  log.Info("Starting Carbon build")

  if minikube.Enabled {
    err := minikube.SetDockerEnv()
    if err != nil {
      log.Fatal(err.Error())
      os.Exit(1)
    }
  }

  log.Info("Reading Carbon config")
  cfgPath, err := filepath.Abs(cfgFile)
  if err != nil {
    log.Fatalf("Failed to find Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

  cfgBody, err := ioutil.ReadFile(cfgPath)
  if err != nil {
    log.Fatalf("Failed to read Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

  cfg, err := rootcfg.ParseConfig(cfgBody)
  if err != nil {
    log.Fatalf("Failed to parse Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

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
    log.Fatalf("Failed to read Kubernetes configs due to the error: %s", err.Error())
    os.Exit(1)
  }

  meta := pkgmeta.New(cfg, cfgBody, kubeManif)

  kd, err := kubernetes.NewKubeDeployment(meta, "image", "tag")
  if err != nil {
    log.Fatalf("Failed to create new instance of KubeDeploy due to the error: %s", err.Error())
    os.Exit(1)
  }

  err = kd.VerifyAll(cfg.Data.KubeManifests)
  if err != nil {
    log.Fatalf("Failed to verify Kubernetes configs due to the error: %s", err.Error())
    os.Exit(1)

  }

  metaMap, err := meta.Serialize()
  if err != nil {
    log.Fatalf("Failed to serialize Carbon config due to the error: %s", err.Error())
    os.Exit(1)
  }

  log.Info("Building Carbon package")
  dockerBuild, err := dockerbuild.NewBuildOptions(cfg, filepath.Dir(cfgPath))
  if err != nil {
    log.Fatalf("Failed to create Docker build handler due to the error: %s", err.Error())
    os.Exit(1)
  }

  err = dockerBuild.Build(metaMap)
  if err != nil {
    log.Fatalf("Failed to build Carbon package due to the error: %s", err.Error())
    os.Exit(1)
  }

  if cfg.HookDefined(rootcfg.HookPostBuild) {
    log.Info("Running post-build hook")
    err = cfg.RunHook(rootcfg.HookPostBuild)
    if err != nil {
      log.Fatalf("Failed to run post-build hook due to the error: %s", err.Error())
      os.Exit(1)
    }
  }

  if buildPush {
    log.Info("Pushing built docker images")
    dockerBuild.Push()
  }

  log.Info("Carbon package has been built successfully")
}
