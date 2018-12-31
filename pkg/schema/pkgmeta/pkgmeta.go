package pkgmeta

import (
  "github.com/pkg/errors"
  // "fmt"
  "time"
  // "encoding/base64"
  // "github.com/starofservice/flapper"

  log "github.com/sirupsen/logrus"

  rootcfglatest "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta/util"
  "github.com/starofservice/carbon/pkg/util/base64"
  // "github.com/starofservice/carbon/pkg/kubernetes"
)

// type MainConfigVersion struct {
//   apiVersion string
// }

// const (
//   MetaPrefix = "c6"
//   MetaDelimiter = "."
// )

type APIVersion struct {
  ApiVersion string
}

var schemaVersions = versions{
  // {v1alpha1.Version, v1alpha1.NewPackageConfig},
  {latest.Version, latest.NewPackageConfig},
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


func New(mainCfg *rootcfglatest.CarbonConfig, rawMainCfg, rawKubeCfg []byte) *latest.PackageConfig {
  mainCfgB64 := base64.Encode(rawMainCfg)
  kubeCfgB64 := base64.Encode(rawKubeCfg)

  p := latest.NewPackageConfigWithVersion()
  p.Name = mainCfg.Name
  p.Version = mainCfg.Version
  p.BuildTime = time.Now().Unix()
  p.MainConfigB64 = mainCfgB64
  p.KubeConfigB64 = kubeCfgB64

  for _, i := range mainCfg.Variables {
    p.Variables = append(p.Variables, latest.PackageConfigVariable{Name: i.Name, Default: i.Default, Description: i.Description})
  }
  return p
}

func SerializeMeta(p latest.PackageConfig) (map[string]string, error) {
  log.Debug("Serializing carbon package metadata")
  fh, err := util.NewFlapper()
  if err != nil {
    panic(err.Error())
  }
  resp, err := fh.Marshal(p)
  if err != nil {
    panic(err.Error())
  }
  return resp, nil
}

// ParseConfig reads a configuration file.
func DeserializeMeta(metaMap map[string]string) (*latest.PackageConfig, error) {
  log.Debug("Deserializing carbon package metadata")
  // buf, err := misc.ReadConfiguration(filename)
  // if err != nil {
  //   return nil, errors.Wrap(err, "read skaffold config")
  // }
  // cfgSrc, err := ioutil.ReadFile(cfgPath)
  // if err != nil {
  //   // log.Fatal(err)
  // }

  apiv := &APIVersion{}
  fh, err := util.NewFlapper()
  if err != nil {
    panic(err.Error())
  }

  if err := fh.Unmarshal(metaMap, apiv); err != nil {
    return nil, errors.Wrap(err, "parsing api version")
  }

  factory, present := schemaVersions.Find(apiv.ApiVersion)
  if !present {
    return nil, errors.Errorf("unknown api version: '%s'", apiv.ApiVersion)
  }

  cfg := factory()
  if err := cfg.Parse(metaMap); err != nil {
    return nil, errors.Wrap(err, "unable to parse config")
  }

  // if err := yamltags.ProcessStruct(cfg); err != nil {
  //   return nil, errors.Wrap(err, "invalid config")
  // }

  parsedCfg := cfg.(*latest.PackageConfig)
  return parsedCfg, nil
}