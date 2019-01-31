package kubernetes_test

import (
  "os"
  "testing"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/kubernetes"
  "github.com/starofservice/carbon/pkg/test"
)

func TestMain(m *testing.M) {
    log.Info("Starting minikube")
    err := test.MinikubeStart()
    if err != nil {
      log.Error("Failed to start minikube due ot the error: ", err.Error())
    }
    defer test.DeferMinikubeDelete()

    code := m.Run()
    
    os.Exit(code)
}


func TestPurgeList(t *testing.T) {
  suite := map[string]bool{
    "rbac.authorization.k8s.io/v1/ClusterRole": false,
    "apiextensions.k8s.io/v1beta1/CustomResourceDefinition": false,
    "/v1/Namespace": false,
    "extensions/v1beta1/Deployment": false,
    "/v1/ServiceAccount": false,
  }

  allRes, err := kubernetes.GetAllResources()
  if err != nil {
    t.Errorf("Failed to receive all resources list due to the error: %s", err.Error())
    return 
  }

  for _, i := range allRes {
    for k := range suite {
      if k == i {
        suite[k] = true
      }
    }
  }
  
  for k, v := range suite {
    if !v {
      t.Errorf("Missing resource: %s", k)
    }
  }
}