package minikube_test

import (
  "os"
  "testing"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/test"
)

const minikubeDriverVar = "CARBON_TEST_MINIKUBE_NONE_DRIVER"

var dockerEnv = []string{
  "DOCKER_TLS_VERIFY",
  "DOCKER_HOST",
  "DOCKER_CERT_PATH",
  "DOCKER_API_VERSION",
}

func TestMinikubeRunning(t *testing.T) {
  err := test.MinikubeStart()
  if err != nil {
    t.Errorf("Failed to start Minikube due ot the error: %s", err.Error())
    return
  }

  err = minikube.CheckStatus()
  if err != nil {
    t.Errorf("Failed to check Minikube status. Got the error: %s", err.Error())
    return
  }

  // 'none' driver does not support 'minikube docker-env' command
  md := os.Getenv(minikubeDriverVar)
  if len(md) == 0 {
    err = minikube.SetDockerEnv()
    if err != nil {
      t.Errorf("Failed to set Docker environment variables. Got the error: %s", err.Error())
      return
    }
    defer minikube.UnsetDockerEnv()

    for _, i := range dockerEnv {
      v := os.Getenv(i)
      if len(v) == 0 {
        t.Errorf("Docker environment variable '%s' is undefined", i)
      }
    }
  }
}

func TestMinikubeStopped(t *testing.T) {
  err := test.MinikubeStop()
  if err != nil {
    t.Errorf("Failed to stop Minikube due ot the error: %s", err.Error())
    return
  }

  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)

  err = minikube.CheckStatus()
  if err == nil {
    t.Errorf("Expected a failure, but the Minikube statue check has been passed successfully")
  }
}

func TestMinikubeDeleted(t *testing.T) {
  err := test.MinikubeDelete()
  if err != nil {
    t.Errorf("Failed to delete Minikube due ot the error: %s", err.Error())
    return
  }

  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)

  err = minikube.CheckStatus()
  if err == nil {
    t.Errorf("Expected a failure, but the Minikube statue check has been passed successfully")
  }
}
