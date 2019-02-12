package latest

import (
  "encoding/json"

  "github.com/pkg/errors"
  "github.com/starofservice/vconf"
)

const Version string = "v1alpha1"


type KubeConfig struct {
  APIVersion    string     `json:"apiVersion"`
  CarbonScope   string     `json:"carbonScope"`
}

func NewKubeConfig() vconf.ConfigInterface {
  return new(KubeConfig)
}

func (c *KubeConfig) GetVersion() string {
  return c.APIVersion
}

func (c *KubeConfig) Parse(data []byte) error {
  if err := json.Unmarshal(data, c); err != nil {
    return err
  }

  return nil
}

func (c *KubeConfig) Upgrade() (vconf.ConfigInterface, error) {
  return nil, errors.New("not implemented yet")
}
