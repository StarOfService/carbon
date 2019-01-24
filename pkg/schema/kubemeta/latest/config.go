package latest

import (
  "encoding/json"

  "github.com/starofservice/carbon/pkg/schema/versioned"
)

const Version string = "v1alpha1"

type KubeMetadata struct {
  APIVersion    string            `json:"apiVersion"`
  Name          string            `json:"name"`
  Version       string            `json:"version"`
  Source        string            `json:"source"`
  Variables     map[string]string `json:"variables"`
  Patches       string            `json:"patches"`
  Manifest      string            `json:"manifest"`
  // hooks
  // dependencies
}

func NewKubeMetadata() versioned.ConfigHandler {
  return new(KubeMetadata)
}

func (c *KubeMetadata) GetVersion() string {
  return c.APIVersion
}

func (c *KubeMetadata) Parse(data []byte) error {
  if err := json.Unmarshal(data, c); err != nil {
    return err
  }

  // if useDefaults {
  //   if err := c.SetDefaultValues(); err != nil {
  //     return errors.Wrap(err, "applying default values")
  //   }
  // }

  return nil
}