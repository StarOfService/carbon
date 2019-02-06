package latest

import (
  yaml "gopkg.in/yaml.v2"
  "errors"
  
  "github.com/starofservice/vconf"
)

const Version string = "v1alpha1"

type CarbonConfigVariable struct {
  Name        string `yaml:"name"`
  Default     string `yaml:"default"`
  Description string `yaml:"description"`
}

type CarbonConfigHooks struct {
  PreBuild  []string `yaml:"pre-build"`
  PostBuild []string `yaml:"post-build"`
}

type CarbonConfig struct {
  APIVersion    string     `yaml:"apiVersion"`
  Name          string     `yaml:"name"`
  Version       string     `yaml:"version"`
  Dockerfile    string     `yaml:"dockerfile"`
  KubeManifests string     `yaml:"kubeManifests"`
  Artifacts     []string   `yaml:"artifacts"`
  Variables     []CarbonConfigVariable `yaml:"variables"`
  Hooks         CarbonConfigHooks `yaml:"hooks"`
  // dependencies
}

func NewCarbonConfig() vconf.ConfigInterface {
  return new(CarbonConfig)
}

func (c *CarbonConfig) GetVersion() string {
  return c.APIVersion
}

func (c *CarbonConfig) Parse(contents []byte) error {
  if err := yaml.UnmarshalStrict(contents, c); err != nil {
    return err
  }

  return nil
}

func (c *CarbonConfig) Upgrade() (vconf.ConfigInterface, error) {
  return nil, errors.New("not implemented yet")
}
