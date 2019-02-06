package kubemeta

import (
  "encoding/json"
  "fmt"
  "strings"

  apicorev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/starofservice/vconf"

  "github.com/starofservice/carbon/pkg/kubernetes"
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

func IsInstalled(name, ns string) (bool, error) {
  log.Debugf("Checking if a component '%s' is installed", name)
  slist, err := getAllSecrets(ns)
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

func Delete(name, ns string) error {
  log.Debug("Deleting Carbon meatadata for package ", name)
  secretHandler, err := getSecretHandler(ns)
  if err != nil {
    return err
  }

  o := &metav1.DeleteOptions{}

  err = secretHandler.Delete(metaObjectPrefix + name, o)
  if err != nil {
    return err
  }

  return nil
}

func Get(name, ns string) (*Handler, error) {
  // log.Debug("Processing Kubernete metadata")

  secretHandler, err := getSecretHandler(ns)
  if err != nil {
    return nil, err
  }

  secretObject, err := secretHandler.Get(getMetaName(name), metav1.GetOptions{})
  if err != nil {
    return nil, err
  }

  return secretToHandler(secretObject)
}


func GetAll(ns string) ([]*Handler, error) {
  // secretHandler, err := getSecretHandler(ns)
  // if err != nil {
  //   return nil, err
  // }

  // label := fmt.Sprintf("%s=%s", metaObjectLabelKey, metaObjectLabelValue)
  // slist, err := secretHandler.List(metav1.ListOptions{LabelSelector: label})
  // if err != nil {
  //   return nil, err
  // }
  slist, err := getAllSecrets(ns)
  if err != nil {
    return nil, err
  }

  var resp []*Handler
  for _, i := range slist.Items {
    km, err := secretToHandler(&i)
    if err != nil {
      return nil, err
    }
    resp = append(resp, km)
  }

  return resp, nil
}

func secretToHandler(secret *apicorev1.Secret) (*Handler, error) {
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
    Namespace: "",
  }

  return km, nil
}

func New(kd *kubernetes.KubeInstall, patches []byte, ns, mns string) *Handler {

  source := kd.Variables.Pkg.DockerName + ":" + kd.Variables.Pkg.DockerTag
  return &Handler{
    Data: latest.KubeMetadata{
      APIVersion: latest.Version,
      Name: kd.Variables.Pkg.Name,
      Version: kd.Variables.Pkg.Version,
      Source: source,
      Variables: kd.Variables.Var,
      Patches: string(patches),
      Namespace: ns,
      Manifest: string(kd.BuiltManifest),
    },
    Namespace: mns,
  }
}

func (self *Handler) Apply() error {
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

  secretHandler, err := getSecretHandler(self.Namespace)
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
  secretHandler, err := getSecretHandler(ns)
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

func getSecretHandler(namespace string) (typedcorev1.SecretInterface, error) {
  kubeConfig, err := kubernetes.GetKubeConfig()
  if err != nil {
    return nil, err
  }
  coreV1Client, err := typedcorev1.NewForConfig(kubeConfig)
  if err != nil {
    return nil, err
  }
  secretHandler := coreV1Client.Secrets(namespace)
  return secretHandler, nil
}

func getMetaName(name string) string {
  return metaObjectPrefix + name
}
