package carboncfg

import (
	"encoding/json"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/starofservice/vconf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
	"github.com/starofservice/carbon/pkg/schema/carboncfg/latest"
	"github.com/starofservice/carbon/pkg/util/tojson"
)

const (
	CarbonConfigItem   = "carbon-config"
	CarbonConfigKey    = "config"
	defaultCarbonScope = "cluster"
)

var schemaVersions = map[string]func() vconf.ConfigInterface{
	latest.Version: latest.NewCarbonConfig,
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

type CarbonConfig struct {
	Data latest.CarbonConfig
}

func New() *CarbonConfig {
	return &CarbonConfig{
		Data: latest.CarbonConfig{
			APIVersion:  latest.Version,
			CarbonScope: defaultCarbonScope,
		},
	}
}

func parseConfig(ns string) (*CarbonConfig, error) {
	log.Debug("Reading Kubernetes config")

	cmHandler, err := kubecommon.GetConfigMapHandler(ns)
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

	parsedCfg := cfg.(*latest.CarbonConfig)
	pc := &CarbonConfig{
		Data: *parsedCfg,
	}
	return pc, nil
}

func FindAndParseConfig() (*CarbonConfig, error) {

	configNamespace := ""
	cmExists, err := configCMExists(kubecommon.CurrentNamespace)
	if err != nil {
		return nil, errors.Wrapf(err, "checking ConfigMap with carbon config at '%s' namespace", kubecommon.CurrentNamespace)
	}
	if cmExists {
		configNamespace = kubecommon.CurrentNamespace
	}

	if configNamespace == "" {
		cmExists, err = configCMExists(kubecommon.GlobalCarbonNamespace)
		if err != nil {
			return nil, errors.Wrapf(err, "checking ConfigMap with carbon config at '%s' namespace", kubecommon.GlobalCarbonNamespace)
		}
		if cmExists {
			configNamespace = kubecommon.GlobalCarbonNamespace
		}
	}

	if configNamespace == "" {
		return New(), nil
	}

	return parseConfig(configNamespace)
}

func configCMExists(ns string) (bool, error) {
	nsExists, err := kubecommon.CheckCarbonNamespace(ns)
	if err != nil {
		return false, errors.Wrapf(err, "checking '%s' namespace from the context", ns)
	}
	if !nsExists {
		return false, nil
	}

	cmHandler, err := kubecommon.GetConfigMapHandler(ns)
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

func (self *CarbonConfig) CarbonScope() (string, error) {
	if self.Data.CarbonScope == "" {
		return defaultCarbonScope, nil
	}

	if self.Data.CarbonScope != "cluster" && self.Data.CarbonScope != "namespace" {
		return "", errors.Errorf("Unknown carbonScope value '%s' at the Kubernetes config", self.Data.CarbonScope)
	}

	return self.Data.CarbonScope, nil
}

func MetaNamespace() (string, error) {
	kcfg, err := FindAndParseConfig()
	if err != nil {
		return "", errors.Wrap(err, "getting Carbon config for Kubernetes cluster")
	}

	scope, err := kcfg.CarbonScope()
	if err != nil {
		return "", err
	}
	if scope == "cluster" {
		return kubecommon.GlobalCarbonNamespace, nil
	}
	return kubecommon.CurrentNamespace, nil
}
