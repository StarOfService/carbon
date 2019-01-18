package rootcfg

import (
  "github.com/pkg/errors"
  // "fmt"
  "gopkg.in/yaml.v2"
  // "io/ioutil"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/util/command"
  "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
  // "github.com/starofservice/carbon/pkg/schema/rootcfg/util"
  "github.com/starofservice/carbon/pkg/schema/versioned"
)

var schemaVersions = map[string]func() versioned.VersionedConfig{
  latest.Version: latest.NewCarbonConfig,
}

const (
  HookPreBuild = "pre-build"
  HookPostBuild = "post-build"
)

func GetCurrentVersion(data []byte) (string, error) {
  type APIVersion struct {
    Version string `yaml:"apiVersion"`
  }
  apiVersion := &APIVersion{}
  if err := yaml.Unmarshal(data, apiVersion); err != nil {
    return "", errors.Wrap(err, "parsing api version")
  }
  return apiVersion.Version, nil
}

type CarbonConfig struct {
 Data latest.CarbonConfig
}

func ParseConfig(data []byte) (*CarbonConfig, error) {
  log.Debug("Processing Carbon config")

  current, err := GetCurrentVersion(data)
  if err != nil {
    return nil, err
  }

  sh := versioned.NewSchemaHandler(current, latest.Version)
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }

  cfg, err := sh.GetLatestConfig(data)
  if err != nil {
    return nil, err
  }

  parsedCfg := cfg.(*latest.CarbonConfig)
  pc := &CarbonConfig{
    Data: *parsedCfg,
  }
  return pc, nil
}

func (self *CarbonConfig) HookDefined(hookType string) bool {
  switch hookType {
  case HookPreBuild:
    if len(self.Data.Hooks.PreBuild) > 0 {
      return true
    } else {
      return false
    }
  case HookPostBuild:
    if len(self.Data.Hooks.PostBuild) > 0 {
      return true
    } else {
      return false
    }
  default:
    return false
  }
}

func (self *CarbonConfig) RunHook(hookType string) error {
  var cmds []string
  switch hookType {
  case HookPreBuild:
    cmds = self.Data.Hooks.PreBuild
  case HookPostBuild:
    cmds = self.Data.Hooks.PostBuild
  default:
    // return fmt.Errorf("Unsupported hook type: %s", hookType)
    return errors.Errorf("Unsupported hook type: %s", hookType)
  }
  for _, i := range cmds {
    err := command.Run(i)
    if err != nil {
      return err
      // return errors.Wrapf(err, "command %s", i)
      // return fmt.Errorf("Failed to run command '%s' due to the error: %s", i, err.Error())

    }
  }
  return nil
}
