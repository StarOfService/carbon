package latest

import (
  yaml "gopkg.in/yaml.v2"

  "github.com/starofservice/carbon/pkg/schema/versioned"
)

const Version string = "carbon/v1alpha1"

type KubeMetadata struct {
  ApiVersion    string            `yaml:"apiVersion"`
  Name          string            `yaml:"name"`
  Version       string            `yaml:"version"`
  Source        string            `yaml:"source"`
  Variables     map[string]string `yaml:"variables"`
  Patches       string            `yaml:"patches"`
  Manifest      string            `yaml:"manifest"`
  // hooks
  // dependencies
}

func NewKubeMetadata() versioned.VersionedConfig {
  return new(KubeMetadata)
}

func (c *KubeMetadata) GetVersion() string {
  return c.ApiVersion
}

func (c *KubeMetadata) Parse(contents []byte) error {
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