package minikube

import (
  "os"
  "os/exec"
  "strings"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
)

const K8sContext = "minikube"
var Enabled = false
var envVars []string

func CheckStatus() error {
  out, err := exec.Command("minikube", "status").Output()
  if err != nil {
    log.Error(string(out))
    return err
  }

  for _, i := range strings.Split(string(out),"\n") {
    raw := strings.TrimSpace(i)
    if strings.HasPrefix(i, "minikube") {
      status := strings.TrimSpace(strings.Replace(raw, "minikube:", "", -1))
      if strings.EqualFold(status, "Running") {
        return nil
      }
      return errors.Errorf("Wrong minikube status. Expected 'Running', but got: %s", status)
    }
  }
  return errors.Errorf("Unable to define minikube status")
}

func SetDockerEnv() error {
  out, err := exec.Command("minikube", "docker-env", "--shell", "bash").Output()
  if err != nil {
    return errors.Wrap(err, "running `minikube docker-env` command")
  }

  for _, i := range strings.Split(string(out),"\n") {
    raw := strings.TrimSpace(i)
    if strings.HasPrefix(i, "export") {
      envVar := strings.TrimSpace(strings.Replace(raw, "export", "", -1))
      envVarSlice := strings.SplitN(envVar, "=", 2)
      envKey := envVarSlice[0]
      envValue := strings.Trim(envVarSlice[1], "\"'")

      envVars = append(envVars, envKey)
      os.Setenv(envKey, envValue)
    }
  }

  return nil
}


func UnsetDockerEnv() {
  for _, i := range envVars {
    err := os.Unsetenv(i)
    if err != nil {
      log.Errorf("Unable to unset environment variable '%s'", i)
    }
  }
}