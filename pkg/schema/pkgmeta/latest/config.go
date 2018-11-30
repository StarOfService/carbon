package latest

import (
  // "github.com/starofservice/flapper"
  // "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta/util"
)

const Version string = "v1alpha1"

type PackageConfigVariable struct {
  Name        string
  Default     string
  Description string
}

type PackageConfig struct {
  ApiVersion    string
  Name          string
  Version       string
  BuildTime     int64
  MainConfigB64 string
  KubeConfigB64 string
  Variables     []PackageConfigVariable
}

func NewPackageConfig() util.VersionedConfig {
  return new(PackageConfig)
}

func NewPackageConfigWithVersion() *PackageConfig {
  p := new(PackageConfig)
  p.ApiVersion = Version
  return p
}

func (c *PackageConfig) GetVersion() string {
  return c.ApiVersion
}

func (c *PackageConfig) Parse(contents map[string]string) error {
  fh, err := util.NewFlapper()
  if err != nil {
    panic(err.Error())
  }
  if err := fh.Unmarshal(contents, c); err != nil {
    return err
  }

  // if useDefaults {
  //   if err := c.SetDefaultValues(); err != nil {
  //     return errors.Wrap(err, "applying default values")
  //   }
  // }

  return nil
}