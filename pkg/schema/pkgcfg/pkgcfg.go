package pkgcfg

import (
  "encoding/json"
  "os"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/vconf"

  "github.com/starofservice/carbon/pkg/schema/pkgcfg/latest"
  "github.com/starofservice/carbon/pkg/util/command"
  "github.com/starofservice/carbon/pkg/util/tojson"
)

var schemaVersions = map[string]func() vconf.ConfigInterface{
  latest.Version: latest.NewCarbonConfig,
}

const (
  HookPreBuild = "pre-build"
  HookPostBuild = "post-build"
)

func GetCurrentVersion(data []byte) (string, error) {
  type VersionStruct struct {
    APIVersion string `json:"apiVersion"`
  }
  version := &VersionStruct{}
  if err := json.Unmarshal(data, version); err != nil {
    return "", errors.Wrap(err, "parsing api version")
  }
  return version.APIVersion, nil
}

type CarbonConfig struct {
  Cwd  string
  Data latest.CarbonConfig
}

func ParseConfig(dir string, data []byte, vpref, vsuf string) (*CarbonConfig, error) {
  log.Debug("Processing Carbon config")

  jsonData, err := tojson.ToJSON(data)
  if err != nil {
    return nil, err
  }

  current, err := GetCurrentVersion(jsonData)
  if err != nil {
    return nil, err
  }

  sh := vconf.NewSchemaHandler(latest.Version)
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }

  cfg, err := sh.GetLatestConfig(current, jsonData)
  if err != nil {
    return nil, err
  }

  parsedCfg := cfg.(*latest.CarbonConfig)

  parsedCfg.Version = vpref + parsedCfg.Version + vsuf

  pc := &CarbonConfig{
    Cwd: dir,
    Data: *parsedCfg,
  }
  return pc, nil
}

func (self *CarbonConfig) HookDefined(hookType string) bool {
  switch hookType {
  case HookPreBuild:
    if len(self.Data.Hooks.PreBuild) > 0 {
      return true
    }
    return false
  case HookPostBuild:
    if len(self.Data.Hooks.PostBuild) > 0 {
      return true
    }
    return false
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
    return errors.Errorf("Unsupported hook type: %s", hookType)
  }

  for _, i := range cmds {
    err := command.Run(i, self.Cwd, os.Stdout, os.Stderr)
    if err != nil {
      return err
    }
  }
  return nil
}
