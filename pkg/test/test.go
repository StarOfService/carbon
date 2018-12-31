package test

import (
  "fmt"
  "os"
  "os/exec"
  "time"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/docker/metadata"
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

func MinikubeStart() error {
  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)

  err := minikube.CheckStatus()
  if err != nil {
    _ = command.Run("minikube delete", "", nil, os.Stderr)
    err = command.Run("minikube start", "", nil, os.Stderr)
    if err != nil {
      return errors.Wrap(err, "starting an new Minikube instance")
    }
  }

  err = command.Run("kubectl config use-context minikube", "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "switching kubectl context to 'minikube'")
  }
  return nil
}

func MinikubeStop() error {
  return command.Run("minikube stop", "", nil, os.Stderr)
}

func MinikubeDelete() error {
  return command.Run("minikube delete", "", nil, os.Stderr)
}

func DeferMinikubeDelete() {
  MinikubeDelete()
}

func DockerRegistryStart() error {
  metadata.SkipTLSVerify = true

  var cmd string
  var err error

  cmd = fmt.Sprintf(`if [ "$(docker ps | grep %s | wc -l)" -gt "0" ]; then docker stop %s; fi`, DockerMockContainerName, DockerMockContainerName)
  cmdHandler := exec.Command("sh", "-c", cmd)
  // cmdHandler.Stdout = os.Stdout
  cmdHandler.Stderr = os.Stderr
  err = cmdHandler.Run()
  if err != nil {
    return errors.Wrap(err, "deleting Docker registry mock container")
  }

  cmd = fmt.Sprintf("docker run -it -d --rm -p %s:5000 --name %s %s", DockerMockRepo, DockerMockContainerName, DockerMockSrcImg)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "starting Docker registry mock")
  }
  time.Sleep(5 * time.Second)

  cmd = fmt.Sprintf("echo '%s' | docker login -u %s --password-stdin %s", DockerMockTestPassword, DockerMockTestUser, DockerMockRepo)
  cmdHandler = exec.Command("sh", "-c", cmd)
  // cmdHandler.Stdout = os.Stdout
  cmdHandler.Stderr = os.Stderr
  err = cmdHandler.Run()
  if err != nil {
    return errors.Wrap(err, "logging in to the Docker registry mock")
  }

  cmd = fmt.Sprintf("docker tag %s %s", DockerMockSrcImg, DockerMockTestImage)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "assigning tag for the testing image")
  }

  cmd = fmt.Sprintf("docker push %s", DockerMockTestImage)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "pushing testing image to the mock Docker registry")
  }

  cmd = fmt.Sprintf("docker rmi %s", DockerMockTestImage)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "deleting local copy of the testing image")
  }

  return nil
}

func DockerRegistryDelete() {
  cmd := fmt.Sprintf("docker stop %s", DockerMockContainerName)
  err := command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
      log.Error("Failed to stop Docker registry mock due to the error: ", err)
  }
}
