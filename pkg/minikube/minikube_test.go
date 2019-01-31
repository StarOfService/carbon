package minikube_test

import (
  "os"
  "testing"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/test"
)

var dockerEnv = []string{
  "DOCKER_TLS_VERIFY",
  "DOCKER_HOST",
  "DOCKER_CERT_PATH",
  "DOCKER_API_VERSION",
}

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

func TestMinikubeRunning(t *testing.T) {
  err := test.MinikubeStart()
  if err != nil {
    t.Errorf("Failed to start minikube due ot the error: %s", err.Error())
    return
  }

  err = minikube.CheckStatus()
  if err != nil {
    t.Errorf("Failed to check minikube status. Got the error: %s", err.Error())
    return
  }

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

func TestMinikubeStopped(t *testing.T) {
  err := test.MinikubeStop()
  if err != nil {
    t.Errorf("Failed to stop minikube due ot the error: %s", err.Error())
    return
  }
  // defer test.MinikubeStart()

  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)

  err = minikube.CheckStatus()
  if err == nil {
    t.Errorf("Expected a failure, but the minikube statue check has been passed successfully")
  }
}

func TestMinikubeDeleted(t *testing.T) {
  err := test.MinikubeDelete()
  if err != nil {
    t.Errorf("Failed to delete minikube due ot the error: %s", err.Error())
    return
  }
  // defer test.MinikubeStart()

  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)
  
  err = minikube.CheckStatus()
  if err == nil {
    t.Errorf("Expected a failure, but the minikube statue check has been passed successfully")
  }
}
