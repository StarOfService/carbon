package cmd

import (
  "io/ioutil"
  "path/filepath"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"

  dockerbuild "github.com/starofservice/carbon/pkg/docker/build"
  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/schema/pkgcfg"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
)

var BuildConfig string
var BuildPush bool
var BuildRemove bool
var BuildTags []string
var BuildTagPrefix string
var BuildTagSuffix string

var buildCmd = &cobra.Command{
  Use:   "build",
  Short: "Build Carbon package",
  Long: `
Builds a Carbon package based on the provided carbon.yaml config.`,
  SilenceErrors: true,
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) != 0 {
      return errors.New("This command doesn't use any arguments")
    }
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    cmd.SilenceUsage = true
    return errors.Wrap(runBuild(), "build")
  },
}

func init() {
  RootCmd.AddCommand(buildCmd)

  buildCmd.Flags().StringVarP(&BuildConfig, "config", "c", "carbon.yaml", "config file (default is carbon.yaml)")
  buildCmd.Flags().BoolVar(&BuildPush, "push", false, "Push built images to the repositories (disabled by default)")
  buildCmd.Flags().BoolVar(&BuildRemove, "rm", false, "Remove build images after the push operation (disabled by default)")
  buildCmd.Flags().StringArrayVar(&BuildTags, "tag", []string{}, "Name and optionally a tag in the 'name:tag' format. If tag isn't provided, it will be replaced by the component version from carbon.yaml")
  buildCmd.Flags().StringVar(&BuildTagPrefix, "tag-prefix", "", "Prefix which should be added for all tags")
  buildCmd.Flags().StringVar(&BuildTagSuffix, "tag-suffix", "", "Suffix which should be added for all tags")
}

func runBuild() error {
  log.Info("Starting Carbon build")

  if BuildRemove && !BuildPush {
    log.Warn("Images can be removed only when push is enabled (see --push option). Skipping it.")
    BuildRemove = false
  }

  if minikube.Enabled && BuildPush {
    log.Warn("Push can't be used with Minikube mode. Skipping it.")
    BuildPush = false
  }

  log.Info("Reading Carbon config")
  cfgPath, err := filepath.Abs(BuildConfig)
  if err != nil {
    return errors.Wrap(err, "looking for Carbon config")
  }

  cfgBody, err := ioutil.ReadFile(cfgPath)
  if err != nil {
    return errors.Wrap(err, "reading Carbon config")
  }

  cfg, err := pkgcfg.ParseConfig(filepath.Dir(cfgPath), cfgBody)
  if err != nil {
    return errors.Wrap(err, "parsing Carbon config")
  }

  if cfg.HookDefined(pkgcfg.HookPreBuild) {
    log.Info("Running pre-build hook")
    if err = cfg.RunHook(pkgcfg.HookPreBuild); err != nil {
      return errors.Wrap(err, "running pre-biuld hooks")
    }
  }

  kubeManif, err := kubernetes.ReadTemplates(cfg)
  if err != nil {
    return errors.Wrap(err, "reading Kubernetes manifest templates")
  }

  meta := pkgmeta.New(cfg, cfgBody, kubeManif)

  kd, err := kubernetes.NewKubeInstall(meta, "image", "tag")
  if err != nil {
    return errors.Wrap(err, "creating new instance of KubeInstall")
  }

  if err = kd.VerifyAll(cfg.Data.KubeManifests); err != nil {
    return errors.Wrap(err, "verifying Kubernetes configs")
  }

  metaMap, err := meta.Serialize()
  if err != nil {
    return errors.Wrap(err, "serializing Carbon config")
  }

  log.Info("Building Carbon package")
  dockerBuild, err := dockerbuild.NewOptions(cfg, filepath.Dir(cfgPath))
  if err != nil {
    return errors.Wrap(err, "creating Docker build handler")
  }
  
  if err = dockerBuild.ExtendTags(BuildTags, BuildTagPrefix, BuildTagSuffix); err != nil {
    return errors.Wrap(err, "extending tags")
  }

  if err = dockerBuild.Build(metaMap); err != nil {
    return errors.Wrap(err, "building Carbon package")
  }

  if cfg.HookDefined(pkgcfg.HookPostBuild) {
    log.Info("Running post-build hook")
    if err = cfg.RunHook(pkgcfg.HookPostBuild); err != nil {
      return errors.Wrap(err, "running post-biuld hooks")
    }
  }

  if BuildPush {
    log.Info("Pushing built Docker images")
    dockerBuild.Push()
  }

  if BuildRemove {
    log.Info("Removing built images")
    dockerBuild.Remove()
  }

  log.Info("Carbon package has been built successfully")
  return nil
}
