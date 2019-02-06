package latest

import (
  "encoding/json"
  "errors"
  
  "github.com/starofservice/vconf"
)

const Version string = "v1alpha1"

type PackageConfigVariable struct {
  Name        string `json:"name"`
  Default     string `json:"default"`
  Description string `json:"description"`
}

type PackageConfig struct {
  APIVersion    string                  `json:"apiVersion"`
  PkgName       string                  `json:"pkgName"`
  PkgVersion    string                  `json:"pkgVersion"`
  BuildTime     int64                   `json:"buildtime"`
  MainConfigB64 string                  `json:"mainConfigB64"`
  KubeConfigB64 string                  `json:"kubeConfigB64"`
  Variables     []PackageConfigVariable `json:"variables"`
}

func NewPackageConfig() vconf.ConfigInterface {
  return new(PackageConfig)
}

func (c *PackageConfig) GetVersion() string {
  return c.APIVersion
}


func (c *PackageConfig) Parse(data []byte) error {
  if err := json.Unmarshal(data, c); err != nil {
    return err
  }

  return nil
}

func (c *PackageConfig) Upgrade() (vconf.ConfigInterface, error) {
  return nil, errors.New("not implemented yet")
}
