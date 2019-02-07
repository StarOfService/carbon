package pkgmeta

import (
  "encoding/json"
  "time"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/vconf"

  "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
  "github.com/starofservice/carbon/pkg/schema/rootcfg"
  "github.com/starofservice/carbon/pkg/util/base64"
)

const (
  carbonMetaLabel = "carbon-package-metadata"
)

var schemaVersions = map[string]func() vconf.ConfigInterface{
  latest.Version: latest.NewPackageConfig,
}

func GetCurrentVersion(data []byte) (string, error) {
  type VersionStruct struct {
    APIVersion string `json:"apiVersion"`
  }
  version := &VersionStruct{}
  if err := json.Unmarshal(data, version); err != nil {
    // return "", errors.Wrap(err, "parsing api version")
    return "", errors.New("provided image doesn't look like a Carbon package")
  }
  return version.APIVersion, nil
}

type PackageConfig struct {
 Data latest.PackageConfig
}

func New(mainCfg *rootcfg.CarbonConfig, rawMainCfg, rawKubeCfg []byte) *PackageConfig {
  mainCfgB64 := base64.Encode(rawMainCfg)
  kubeCfgB64 := base64.Encode(rawKubeCfg)

  p := &PackageConfig{
    Data: latest.PackageConfig{
      APIVersion: latest.Version,
      PkgName: mainCfg.Data.Name,
      PkgVersion: mainCfg.Data.Version,
      BuildTime: time.Now().Unix(),
      MainConfigB64: mainCfgB64,
      KubeConfigB64: kubeCfgB64,      
    },
  }

  for _, i := range mainCfg.Data.Variables {
    p.Data.Variables = append(p.Data.Variables, latest.PackageConfigVariable{Name: i.Name, Default: i.Default, Description: i.Description})
  }
  return p
}

func (self *PackageConfig) Serialize() (map[string]string, error) {
  log.Debug("Serializing carbon package metadata")
  data, err := json.Marshal(self.Data)
  if err != nil {
    return nil, err
  }

  resp := map[string]string{
    carbonMetaLabel: string(data),
  }
  return resp, nil
}

func Deserialize(metaMap map[string]string) (*PackageConfig, error) {
  log.Debug("Deserializing carbon package metadata")
  
  data := []byte(metaMap[carbonMetaLabel])

  current, err := GetCurrentVersion(data)
  if err != nil {
    return nil, err
  }

  sh := vconf.NewSchemaHandler(latest.Version)
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }

  cfg, err := sh.GetLatestConfig(current, data)
  if err != nil {
    return nil, err
  }

  parsedCfg := cfg.(*latest.PackageConfig)
  pc := &PackageConfig{
    Data: *parsedCfg,
  }
  return pc, nil
}
