package carboncfg

import (
  "encoding/json"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/vconf"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/carboncfg/latest"
  "github.com/starofservice/carbon/pkg/util/tojson"
)

const (
  CarbonConfigItem = "carbon-config"
  CarbonConfigKey = "config"
  defaultCarbonScope = "cluster"
)

var schemaVersions = map[string]func() vconf.ConfigInterface{
  latest.Version: latest.NewKubeConfig,
}

func GetCurrentVersion(data []byte) (string, error) {
  type VersionStruct struct {
    APIVersion string `json:"apiVersion"`
  }
  version := &VersionStruct{}
  if err := json.Unmarshal(data, version); err != nil {
    return "", errors.Wrap(err, "parsing api version")
  }
  return version.APIVersion, nil
}

type KubeConfig struct {
  Data latest.KubeConfig
}

func New() *KubeConfig {
  return &KubeConfig{
    Data: latest.KubeConfig{
      APIVersion: latest.Version,
      CarbonScope: defaultCarbonScope,
    },
  }
}

func parseConfig(ns string) (*KubeConfig, error) {
  log.Debug("Reading Kubernetes config")

  cmHandler, err := kubernetes.GetConfigMapHandler(ns)
  if err != nil {
    return nil, err
  }

  cmObject, err := cmHandler.Get(CarbonConfigItem, metav1.GetOptions{})
  if err != nil {
    return nil, err
  }

  data := []byte(cmObject.Data[CarbonConfigKey])
  jsonData, err := tojson.ToJSON(data)
  if err != nil {
    return nil, err
  }

  current, err := GetCurrentVersion(jsonData)
  if err != nil {
    return nil, err
  }

  sh := vconf.NewSchemaHandler(latest.Version)
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }

  cfg, err := sh.GetLatestConfig(current, jsonData)
  if err != nil {
    return nil, err
  }

  parsedCfg := cfg.(*latest.KubeConfig)
  pc := &KubeConfig{
    Data: *parsedCfg,
  }
  return pc, nil
}

func FindAndParseConfig() (*KubeConfig, error)  {

  configNamespace := ""
  cmExists, err := configCMExists(kubernetes.CurrentNamespace)
  if err != nil {
    return nil, errors.Wrapf(err, "checking ConfigMap with carbon config at '%s' namespace", kubernetes.CurrentNamespace)
  }
  if cmExists {
    configNamespace = kubernetes.CurrentNamespace
  }

  if configNamespace == "" {
    cmExists, err = configCMExists(kubernetes.GlobalCarbonNamespace)
    if err != nil {
      return nil, errors.Wrapf(err, "checking ConfigMap with carbon config at '%s' namespace", kubernetes.GlobalCarbonNamespace)
    }
    if cmExists {
      configNamespace = kubernetes.GlobalCarbonNamespace
    }
  }

  if configNamespace == "" {
    return New(), nil
  }

  return parseConfig(configNamespace)
}


func configCMExists(ns string) (bool, error) {
  nsExists, err := kubernetes.CheckCarbonNamespace(ns)
  if err != nil {
    return false, errors.Wrapf(err, "checking '%s' namespace from the context", ns)
  }
  if !nsExists {
    return false, nil
  }

  cmHandler, err := kubernetes.GetConfigMapHandler(ns)
  if err != nil {
    return false, err
  }

  cmExists := false
  cmList, err := cmHandler.List(metav1.ListOptions{})
  for _, i := range cmList.Items {
    if i.ObjectMeta.Name == CarbonConfigItem {
      cmExists = true
    }
  }

  return cmExists, nil
}

func carbonScope() (string, error) {
  kcfg, err := FindAndParseConfig()
  if err != nil {
    return "", errors.Wrap(err, "getting Carbon config for Kubernetes cluster")
  }

  if kcfg.Data.CarbonScope == "" {
    return defaultCarbonScope, nil
  }

  if kcfg.Data.CarbonScope != "cluster" && kcfg.Data.CarbonScope != "namespace" {
    return "", errors.Errorf("Unknown carbonScope value '%s' at the Kubernetes config", kcfg.Data.CarbonScope)
  }

  return kcfg.Data.CarbonScope, nil
}

func MetaNamespace() (string, error) {
  scope, err := carbonScope()
  if err != nil {
    return "", err
  }
  if scope == "cluster" {
    return kubernetes.GlobalCarbonNamespace, nil
  }
  return kubernetes.CurrentNamespace, nil
}
