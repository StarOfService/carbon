package kubemeta

import (
  "encoding/json"
  "fmt"
  "strings"

  apicorev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/vconf"

  "github.com/starofservice/carbon/pkg/kubernetes"
  kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
  "github.com/starofservice/carbon/pkg/schema/carboncfg"
  "github.com/starofservice/carbon/pkg/schema/kubemeta/latest"
)

const (
  metaObjectPrefix = "carbon-package-metadata-"
  metaObjectLabelKey = "carbon/type"
  metaObjectLabelValue = "package-metadata"
  metaObjectKey = "metadata"
)

var schemaVersions = map[string]func() vconf.ConfigInterface{
  latest.Version: latest.NewKubeMetadata,
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

type Handler struct {
  Data latest.KubeMetadata
  Namespace string
}

func IsInstalled(name string) (bool, error) {
  log.Debugf("Checking if a component '%s' is installed", name)

  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    return false, err
  }

  slist, err := getAllSecrets(mns)
  if err != nil {
    return false, err
  }

  for _, i := range slist.Items {
    log.Trace("component: ", i.ObjectMeta.Name)
    if i.ObjectMeta.Name == metaObjectPrefix + name {
      return true, nil
    }
  }
  return false, nil
}

func (self *Handler) Delete() error {
  log.Debug("Deleting Carbon meatadata for package ", self.Data.Name)

  secretHandler, err := kubecommon.GetSecretHandler(self.Namespace)
  if err != nil {
    return err
  }

  err = secretHandler.Delete(metaObjectPrefix + self.Data.Name, &metav1.DeleteOptions{})
  if err != nil {
    return err
  }

  return nil
}

func Get(name string) (*Handler, error) {
  // log.Debug("Processing Kubernete metadata")
  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    return nil, err
  }

  secretHandler, err := kubecommon.GetSecretHandler(mns)
  if err != nil {
    return nil, err
  }

  secretObject, err := secretHandler.Get(getMetaName(name), metav1.GetOptions{})
  if err != nil {
    return nil, err
  }

  return secretToHandler(secretObject, mns)
}

func GetAll() ([]*Handler, error) {
  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    return nil, err
  }

  var resp []*Handler
  slist, err := getAllSecrets(mns)
  if err != nil {
    return nil, errors.Wrap(err, "getting secret list")
  }
  for _, si := range slist.Items {
    km, err := secretToHandler(&si, mns)
    if err != nil {
      return nil, errors.Wrapf(err, "extracting installed package metadata from a secret '%s' in namespace '%s'", si.ObjectMeta.Name, mns)
    }
    resp = append(resp, km)
  }
  return resp, nil
}

func secretToHandler(secret *apicorev1.Secret, mns string) (*Handler, error) {
  data := secret.Data[metaObjectKey]

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

  parsedCfg := cfg.(*latest.KubeMetadata)
  km := &Handler{
    Data: *parsedCfg,
    Namespace: mns,
  }

  return km, nil
}

// func New(kd *kubernetes.KubeInstall, patches []byte, ns, mns string) *Handler {
func New(kd *kubernetes.KubeInstall, patches []byte) (*Handler, error) {
  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    return nil, err
  }

  source := kd.Variables.Pkg.DockerName + ":" + kd.Variables.Pkg.DockerTag
  return &Handler{
    Data: latest.KubeMetadata{
      APIVersion: latest.Version,
      Name: kd.Variables.Pkg.Name,
      Version: kd.Variables.Pkg.Version,
      Source: source,
      Variables: kd.Variables.Var,
      Patches: string(patches),
      Namespace: kubecommon.CurrentNamespace,
      Manifest: string(kd.BuiltManifest),
    },
    Namespace: mns,
  }, nil
}

func (self *Handler) Apply() error {
  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    return err
  }

  if mns == kubecommon.GlobalCarbonNamespace {
    nsExists, err := kubecommon.CheckCarbonNamespace(kubecommon.GlobalCarbonNamespace)
    if err != nil {
      return errors.Wrapf(err, "looking for '%s' namespace", kubecommon.GlobalCarbonNamespace)
    }

    if !nsExists {
      err = kubecommon.CreateGlobalCarbonNamespace()
      if err != nil {
        return errors.Wrapf(err, "creating '%s' namespace", kubecommon.GlobalCarbonNamespace)
      }
    }
  }

  data, err := json.Marshal(self.Data)
  if err != nil {
    return err
  }

  meta := &apicorev1.Secret{
    ObjectMeta: metav1.ObjectMeta{
      Labels: map[string]string{metaObjectLabelKey: metaObjectLabelValue},
      Name: getMetaName(self.Data.Name),
      Namespace: self.Namespace,
    },
    Data: map[string][]byte{metaObjectKey: data},
  }

  secretHandler, err := kubecommon.GetSecretHandler(self.Namespace)
  if err != nil {
    return err
  }

  _, err = secretHandler.Get(getMetaName(self.Data.Name), metav1.GetOptions{})
  if err != nil {
    if strings.Contains(err.Error(), "not found") {
      _, err = secretHandler.Create(meta)
      if err != nil {
        return err
      }
    } else {
      return err
    }
  } else {
    _, err = secretHandler.Update(meta)
    if err != nil {
      return err
    }
  }

  return nil
}

func getAllSecrets(ns string) (*apicorev1.SecretList, error) {
  secretHandler, err := kubecommon.GetSecretHandler(ns)
  if err != nil {
    return nil, err
  }

  label := fmt.Sprintf("%s=%s", metaObjectLabelKey, metaObjectLabelValue)
  slist, err := secretHandler.List(metav1.ListOptions{LabelSelector: label})
  if err != nil {
    return nil, err
  }

  return slist, nil
}

func getMetaName(name string) string {
  return metaObjectPrefix + name
}
