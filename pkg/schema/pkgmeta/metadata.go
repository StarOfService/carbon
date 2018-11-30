package pkgmeta

import (
  "github.com/pkg/errors"
  "fmt"
  "time"
  "encoding/base64"
  // "github.com/starofservice/flapper"


  rootcfglatest "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta/util"
  "github.com/starofservice/carbon/pkg/kubernetes"
)

// type MainConfigVersion struct {
//   apiVersion string
// }

// const (
//   MetaPrefix = "c6"
//   MetaDelimiter = "."
// )

type APIVersion struct {
  Version string
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



// ParseConfig reads a configuration file.
func ParseMetadata(metaMap map[string]string) (*latest.PackageConfig, error) {
  // buf, err := misc.ReadConfiguration(filename)
  // if err != nil {
  //   return nil, errors.Wrap(err, "read skaffold config")
  // }
  // cfgSrc, err := ioutil.ReadFile(cfgPath)
  // if err != nil {
  //   // log.Fatal(err)
  // }

  apiVersion := &APIVersion{}
  fh, err := util.NewFlapper()
  if err != nil {
    panic(err.Error())
  }

  if err := fh.Unmarshal(metaMap, apiVersion); err != nil {
    return nil, errors.Wrap(err, "parsing api version")
  }

  factory, present := schemaVersions.Find(apiVersion.Version)
  if !present {
    return nil, errors.Errorf("unknown api version: '%s'", apiVersion.Version)
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

func New(mainCfg *rootcfglatest.CarbonConfig, rawMainCfg, rawKubeCfg []byte) *latest.PackageConfig {
  mainCfgB64 := B64Encode(rawMainCfg)
  kubeCfgB64 := B64Encode(rawKubeCfg)

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

func NewKubeDeploy(p *latest.PackageConfig) (*kubernetes.KubeDeploy, error) {
  k8sManifest, err := B64Decode(p.KubeConfigB64)
  if err != nil {return nil, err}

  kd := &kubernetes.KubeDeploy{
    Manifest: k8sManifest,
    Variables: kubernetes.DeployVars{
      Pkg: make(map[string]string),
      Var: make(map[string]string),
    },
  }

  for _, v := range p.Variables {
    kd.Variables.Var[v.Name] = v.Default
  }

  return kd, nil
}

func Map(p latest.PackageConfig) (map[string]string, error) {
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

func B64Encode(data []byte) string {
  return base64.StdEncoding.EncodeToString(data)
}

func B64Decode(data string) ([]byte, error) {
  resp, err := base64.StdEncoding.DecodeString(data)
  if err != nil {
    return nil, fmt.Errorf("Unable to decode string `%s` due to the error: %s", data, err.Error())
  }
  return resp, nil
}