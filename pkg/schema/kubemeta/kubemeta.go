package kubemeta

import (
  "encoding/json"
  "github.com/pkg/errors"
  "fmt"
  "strings"
  // "io/ioutil"
  // log "github.com/sirupsen/logrus"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  apicorev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/versioned"

  "github.com/starofservice/carbon/pkg/schema/kubemeta/latest"
)

const (
  // defaultKubeNamespace = "default"
  metaObjectPrefix = "carbon-package-metadata-"
  metaObjectLabelKey = "carbon/type"
  metaObjectLabelValue = "package-metadata"
  metaObjectKey = "metadata"
)

var schemaVersions = map[string]func() versioned.VersionedConfig{
  latest.Version: latest.NewKubeMetadata,
}

func GetCurrentVersion(data []byte) (string, error) {
  type APIVersion struct {
    Version string `json:"apiVersion"`
  }
  apiVersion := &APIVersion{}
  if err := json.Unmarshal(data, apiVersion); err != nil {
    return "", errors.Wrap(err, "parsing api version")
  }
  return apiVersion.Version, nil
}

// struct KubeMetadata latest.KubeMetadata
type KubeMetaHandler struct {
 Data latest.KubeMetadata
 Namespace string
}

// func ParseKubeMetaHandler(body []byte) (*KubeMetaHandler, error) {
func Get(name, ns string) (*KubeMetaHandler, error) {
  // log.Debug("Processing Kubernete metadata")
  // metaName := metaObjectPrefix + name
  // if ns == "" {
  //   ns = defaultKubeNamespace
  // }

  secretHandler, err := getSecretHandler(ns)
  if err != nil {
    return nil, err
  }

  secretObject, err := secretHandler.Get(getMetaName(name), metav1.GetOptions{})
  if err != nil {
    return nil, err
  }
  // data := secretObject.Data[metaObjectKey]

  // current, err := GetCurrentVersion(data)
  // if err != nil {
  //   return nil, err
  // }

  // sh := versioned.NewSchemaHandler(current, latest.Version)
  // for k, v := range schemaVersions {
  //   sh.RegVersion(k, v)
  // }
  // // sh.RegVersion(latest.Version, latest.NewKubeMetadata)

  // cfg, err := sh.GetLatestConfig(data)
  // if err != nil {
  //   return nil, err
  // }

  // parsedCfg := cfg.(*latest.KubeMetadata)
  // km := &KubeMetaHandler{
  //   Data: *parsedCfg,
  //   Namespace: "",
  // }
  // // parsedCfg := cfg.(*latest.KubeMetadata)
  // // return parsedCfg, nil
  // return km, nil
  return secretToKubeMetaHandler(secretObject)

}


func GetAll(ns string) ([]*KubeMetaHandler, error) {
  secretHandler, err := getSecretHandler(ns)
  if err != nil {
    return nil, err
  }

  label := fmt.Sprintf("%s=%s", metaObjectLabelKey, metaObjectLabelValue)
  slist, err := secretHandler.List(metav1.ListOptions{LabelSelector: label})
  if err != nil {
    return nil, err
  }

  var resp []*KubeMetaHandler
  for _, i := range slist.Items {
    km, err := secretToKubeMetaHandler(&i)
    if err != nil {
      return nil, err
    }
    resp = append(resp, km)
  }

  return resp, nil
}

func secretToKubeMetaHandler(secret *apicorev1.Secret) (*KubeMetaHandler, error) {
  data := secret.Data[metaObjectKey]

  current, err := GetCurrentVersion(data)
  if err != nil {
    return nil, err
  }

  sh := versioned.NewSchemaHandler(current, latest.Version)
  for k, v := range schemaVersions {
    sh.RegVersion(k, v)
  }
  // sh.RegVersion(latest.Version, latest.NewKubeMetadata)

  cfg, err := sh.GetLatestConfig(data)
  if err != nil {
    return nil, err
  }

  parsedCfg := cfg.(*latest.KubeMetadata)
  km := &KubeMetaHandler{
    Data: *parsedCfg,
    Namespace: "",
  }
  // parsedCfg := cfg.(*latest.KubeMetadata)
  // return parsedCfg, nil
  return km, nil
}

func New(kd *kubernetes.KubeDeployment, patches []byte, ns string) *KubeMetaHandler {
  // if ns == "" {
  //   ns = defaultKubeNamespace
  // }

  source := kd.Variables.Pkg.DockerName + ":" + kd.Variables.Pkg.DockerTag
  return &KubeMetaHandler{
    Data: latest.KubeMetadata{
      ApiVersion: latest.Version,
      Name: kd.Variables.Pkg.Name,
      Version: kd.Variables.Pkg.Version,
      Source: source,
      Variables: kd.Variables.Var,
      Patches: string(patches),
      Manifest: string(kd.BuiltManifest),
    },
    Namespace: ns,
  }
}

func (self *KubeMetaHandler) Apply() error {
  data, err := json.Marshal(self.Data)
  if err != nil {
    return err
  }

  // metaName := metaObjectPrefix + self.Data.Name
  meta := &apicorev1.Secret{
    ObjectMeta: metav1.ObjectMeta{
      Labels: map[string]string{metaObjectLabelKey: metaObjectLabelValue},
      Name: getMetaName(self.Data.Name),
      Namespace: self.Namespace,
    },
    Data: map[string][]byte{metaObjectKey: data},
  }

  // kubeConfig, err := kubernetes.GetKubeConfig()
  // if err != nil {
  //   return err
  // }
  // coreV1Client, err := typedcorev1.NewForConfig(kubeConfig)
  // if err != nil {
  //   return err
  // }
  // secretHandler := coreV1Client.Secrets(self.Namespace)
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
  // fmt.Println("Getting item:")
  // fmt.Println(item)
  // fmt.Println(err)


  // _, err = secretHandler.Create(meta)
  // if err != nil {
  //   return err
  // }
  return nil
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
