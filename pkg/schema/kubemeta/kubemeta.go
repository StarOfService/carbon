package kubemeta

import (
  "github.com/pkg/errors"
  // "fmt"
  // "gopkg.in/yaml.v2"
  // "io/ioutil"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/kubemeta/latest"
  "github.com/starofservice/carbon/pkg/schema/versioned"
)

struct KubeMetadata latest.KubeMetadata

func GetCurrentVersion(body []byte) (string, error) {
  type APIVersion struct {
    Version string `yaml:"apiVersion"`
  }
  apiVersion := &APIVersion{}
  if err := yaml.Unmarshal(metaBody, apiVersion); err != nil {
    return nil, errors.Wrap(err, "parsing api version")
  }
  return apiVersion.Version
}

func ParseKubeMetadata(body []byte) (*KubeMetadata, error) {
  log.Debug("Processing Kubernete metadata")

  current, err := GetCurrentVersion(body)
  if err != nil {
    return nil, err
  }

  sh := versioned.NewSchemaHandler(current, latest.Version)
  sh.RegVersion(latest.Version, latest.NewCarbonConfig)
  
  cfg, err := sh.GetLatestConfig(body)
  if err != nil {
    return nil, err
  }

  parsedCfg := cfg.(*KubeMetadata)
  return parsedCfg, nil
}

func NewKubeMetadata(kd *kubernetes.KubeDeployment, patches []byte) (*KubeMetadata, error) {
  source := kd.Variables.Pkg.DockerName + ":" + kd.Variables.Pkg.DockerTag
  return &latest.KubeMetadata{
    ApiVersion: latest.Version,
    Name: kd.Variables.Pkg.Name,
    Version: kd.Variables.Pkg.Version,
    Source: source,
    Variables: kd.Variables.Var,
    Patches: string(patches),
    Manifest: kd.BuiltManifest,
  }
}

// func (self *KubeMetadata) ToJson() []byte {
  
// }