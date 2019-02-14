package carboncfg_test

import (
  "os"
  "testing"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
  "github.com/stretchr/testify/assert"
  apicorev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/schema/carboncfg"
  // "github.com/starofservice/carbon/pkg/test"
)

func TestMain(m *testing.M) {
    // log.Info("Starting Minikube")
    // err := test.MinikubeStart()
    // if err != nil {
    //   log.Error("Failed to start Minikube due ot the error: ", err.Error())
    // }

    err := kubernetes.SetNamespace("")
    if err != nil {
      log.Error("Failed to set current namespace due to the error ", err.Error())
      return
    }

    code := m.Run()

    os.Exit(code)
}

func TestCarboncfgWithoutConfig(t *testing.T) {
  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    t.Errorf("Failed to get metadata namespace due to the error %s", err.Error())
    return
  }

  assert.Equal(t, mns, kubernetes.GlobalCarbonNamespace, "they should be equal")
}

func TestCarboncfgGlobalConfig(t *testing.T) {
  config := `{"apiVersion": "v1alpha1", "carbonScope": "namespace"}`
  err := createCarbonConfig(kubernetes.GlobalCarbonNamespace, config)
  if err != nil {
    t.Errorf("Failed to create Kubernetes Carbon config due to the error %s", err.Error())
    return
  }
  defer deleteCarbonConfig(t, kubernetes.GlobalCarbonNamespace)

  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    t.Errorf("Failed to get metadata namespace due to the error %s", err.Error())
    return
  }

  assert.Equal(t, mns, kubernetes.CurrentNamespace, "they should be equal")
}

func TestCarboncfgGlobalAndLocalConfig(t *testing.T) {
  config := `{"apiVersion": "v1alpha1", "carbonScope": "cluster"}`
  err := createCarbonConfig(kubernetes.GlobalCarbonNamespace, config)
  if err != nil {
    t.Errorf("Failed to create Kubernetes Carbon config due to the error %s", err.Error())
    return
  }
  defer deleteCarbonConfig(t, kubernetes.GlobalCarbonNamespace)

  config = `{"apiVersion": "v1alpha1", "carbonScope": "namespace"}`
  err = createCarbonConfig(kubernetes.CurrentNamespace, config)
  if err != nil {
    t.Errorf("Failed to create Kubernetes Carbon config due to the error %s", err.Error())
    return
  }
  defer deleteCarbonConfig(t, kubernetes.CurrentNamespace)

  mns, err := carboncfg.MetaNamespace()
  if err != nil {
    t.Errorf("Failed to get metadata namespace due to the error %s", err.Error())
    return
  }

  assert.Equal(t, mns, kubernetes.CurrentNamespace, "they should be equal")
}

func createCarbonConfig(ns, data string) error {
  cmh, err := kubernetes.GetConfigMapHandler(ns)
  if err != nil {
    return errors.Wrap(err, "getting ConfigMap handler")
  }

  cm := &apicorev1.ConfigMap{
    ObjectMeta: metav1.ObjectMeta{
      Name: carboncfg.CarbonConfigItem,
    },
    Data: map[string]string{carboncfg.CarbonConfigKey: data},
  }

  _, err = cmh.Create(cm)
  if err != nil {
    return errors.Wrapf(err, "creating ConfigMap item")
  }

  return nil
}

func deleteCarbonConfig(t *testing.T, ns string) {
  cmh, err := kubernetes.GetConfigMapHandler(ns)
  if err != nil {
    t.Errorf("Failed to get ConfigMap handler due to the error %s", err.Error())
    return
  }

  _ = cmh.Delete(carboncfg.CarbonConfigItem, &metav1.DeleteOptions{})
  if err != nil {
    t.Errorf("Failed to delete ConfigMap config itme '%s' due to the error %s", carboncfg.CarbonConfigItem, err.Error())
    return
  }
}
