package rootcfg

import (
  "github.com/pkg/errors"
  // "fmt"
  "gopkg.in/yaml.v2"
  // "io/ioutil"

  "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
  "github.com/starofservice/carbon/pkg/schema/rootcfg/util"
)

// type MainConfigVersion struct {
//   apiVersion string
// }

type APIVersion struct {
  Version string `yaml:"apiVersion"`
}

var schemaVersions = versions{
  // {v1alpha1.Version, v1alpha1.NewCarbonConfig},
  {latest.Version, latest.NewCarbonConfig},
}

type version struct {
  apiVersion string
  factory    func() util.VersionedConfig
}

type versions []version

// Find search the constructor for a given api version.
func (v *versions) Find(apiVersion string) (func() util.VersionedConfig, bool) {
  for _, version := range *v {
    if version.apiVersion == apiVersion {
      return version.factory, true
    }
  }

  return nil, false
}

// func ParseConfigPath(cfgPath string) (*latest.CarbonConfig, error) {
//   // buf, err := misc.ReadConfiguration(filename)
//   // if err != nil {
//   //   return nil, errors.Wrap(err, "read skaffold config")
//   // }
//   cfgBody, err := ioutil.ReadFile(cfgPath)
//   if err != nil {
//     // log.Fatal(err)
//   }
//   return ParseConfig
// }

// ParseConfig reads a configuration file.
func ParseConfig(cfgBody []byte) (*latest.CarbonConfig, error) {
  apiVersion := &APIVersion{}
  if err := yaml.Unmarshal(cfgBody, apiVersion); err != nil {
    return nil, errors.Wrap(err, "parsing api version")
  }

  factory, present := schemaVersions.Find(apiVersion.Version)
  if !present {
    return nil, errors.Errorf("unknown api version: '%s'", apiVersion.Version)
  }

  cfg := factory()
  if err := cfg.Parse(cfgBody); err != nil {
    return nil, errors.Wrap(err, "unable to parse config")
  }

  // if err := yamltags.ProcessStruct(cfg); err != nil {
  //   return nil, errors.Wrap(err, "invalid config")
  // }

  parsedCfg := cfg.(*latest.CarbonConfig)
  return parsedCfg, nil
}
