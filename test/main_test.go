package test

import (
  "fmt"
  "os"
  "testing"
  "time"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/util/command"
)

const (
  DockerMockSrcImg = "starofservice/docker-registry-mock:latest"
  DockerMockContainerName = "docker-registry-mock"
  DockerMockRepo = "127.0.0.1:5000"
  DockerMockTestUser = "testuser"
  DockerMockTestPassword = "testpassword"
)

var DockerMockTestImage = DockerMockRepo + "/docker-registry-mock:latest"

func TestMain(m *testing.M) {
    log.Info("Starting minikube")
    err := MinikubeStart()
    if err != nil {
      log.Fatal("Failed to start minikube due ot the error: ", err.Error())
    }

    log.Info("Starting docker registry")
    err = DockerRegistryStart()
    if err != nil {
      log.Fatal("Failed to start docker registry due ot the error: ", err.Error())
    }

    code := m.Run()

    // TODO Defer
    log.Info("Deleting minikube")
    err = MinikubeDelete()
    if err != nil {
      log.Fatal("Failed to stop minikube due ot the error: ", err.Error())
    }

    // TODO Defer
    log.Info("Deleting docker registry")
    err = DockerRegistryDelete()
    if err != nil {
      log.Fatal("Failed to delete docker registry container due ot the error: ", err.Error())
    }
    
    os.Exit(code)
}

func MinikubeStart() error {
  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)

  err := minikube.CheckStatus()
  if err != nil {
    err = command.Run("minikube start", nil, nil)
    if err != nil {
      return err
    }
  }

  err = command.Run("kubectl config use-context minikube", nil, nil)
  if err != nil {
    log.Errorf("Failed to switch kubectl context to 'minikube' due to the error: %s", err.Error())
    os.Exit(1)
  }
  return nil
}

func MinikubeStop() error {
  return command.Run("minikube stop", nil, nil)
}

func MinikubeDelete() error {
  minikube.UnsetDockerEnv()
  return command.Run("minikube delete", nil, nil)
}

func DockerRegistryStart() (error) {
  var cmd string
  var err error

  cmd = fmt.Sprintf("docker run -it -d --rm -p %s:5000 --name %s %s", DockerMockRepo, DockerMockContainerName, DockerMockSrcImg)
  log.Info(cmd)
  err = command.Run(cmd, nil, os.Stderr)
  if err != nil {
      log.Fatal("Failed to start docker registry mock due to the error: ", err)
  }
  time.Sleep(10 * time.Second)


  cmd = fmt.Sprintf("docker login -u %s -p %s %s", DockerMockTestUser, DockerMockTestPassword, DockerMockRepo)
  err = command.Run(cmd, nil, os.Stderr)
  if err != nil {
      log.Fatal("Failed to login at the docker registry mock due to the error: ", err)
  }

  cmd = fmt.Sprintf("docker tag %s %s", DockerMockSrcImg, DockerMockTestImage)
  err = command.Run(cmd, nil, os.Stderr)
  if err != nil {
      log.Fatal("Failed to assign tag for the mock docker registry due to the error: ", err)
  }
  
  cmd = fmt.Sprintf("docker push %s", DockerMockTestImage)
  err = command.Run("docker push 127.0.0.1:5000/docker-registry-mock:latest", nil, nil)
  if err != nil {
      log.Fatal("Failed to push docker image to the  mock docker registry due to the error: ", err)
  }

  return nil
}

func DockerRegistryDelete() error {
  // return pool.Purge(resource)
  // var cmd string
  // var err error

  cmd := fmt.Sprintf("docker stop %s", DockerMockContainerName)
  err := command.Run(cmd, nil, os.Stderr)
  if err != nil {
      log.Fatal("Failed to stop docker registry mock due to the error: ", err)
  }
  return nil
}