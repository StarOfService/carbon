package latest

import (
  yaml "gopkg.in/yaml.v2"

  "github.com/starofservice/carbon/pkg/schema/rootcfg/util"
)

const Version string = "carbon/v1alpha1"

type CarbonConfig struct {
  ApiVersion    string     `yaml:"apiVersion"`
  Name          string     `yaml:"name"`
  Version       string     `yaml:"version"`
  Dockerfile    string     `yaml:"dockerfile"`
  KubeManifests string     `yaml:"kubeManifests"`
  Artifacts     []string   `yaml:"artifacts"`
  Variables     []CarbonConfigVariable `yaml:"variables"`
  // hooks
  // dependencies
}

type CarbonConfigVariable struct {
  Name        string `yaml:"name"`
  Default     string `yaml:"default"`
  Description string `yaml:"description"`
}

func NewCarbonConfig() util.VersionedConfig {
  return new(CarbonConfig)
}

func (c *CarbonConfig) GetVersion() string {
  return c.ApiVersion
}

func (c *CarbonConfig) Parse(contents []byte) error {
  if err := yaml.UnmarshalStrict(contents, c); err != nil {
    return err
  }

  // if useDefaults {
  //   if err := c.SetDefaultValues(); err != nil {
  //     return errors.Wrap(err, "applying default values")
  //   }
  // }

  return nil
}