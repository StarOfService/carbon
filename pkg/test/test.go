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

// func TestMain(m *testing.M) {
//     log.Info("Starting minikube")
//     err := MinikubeStart()
//     if err != nil {
//       log.Fatal("Failed to start minikube due ot the error: ", err.Error())
//     }

//     log.Info("Starting docker registry")
//     err = DockerRegistryStart()
//     if err != nil {
//       log.Fatal("Failed to start docker registry due ot the error: ", err.Error())
//     }

//     code := m.Run()

//     // TODO Defer
//     log.Info("Deleting minikube")
//     err = MinikubeDelete()
//     if err != nil {
//       log.Fatal("Failed to stop minikube due ot the error: ", err.Error())
//     }

//     // TODO Defer
//     log.Info("Deleting docker registry")
//     err = DockerRegistryDelete()
//     if err != nil {
//       log.Fatal("Failed to delete docker registry container due ot the error: ", err.Error())
//     }
    
//     os.Exit(code)
// }

func MinikubeStart() error {
  log.SetLevel(log.FatalLevel)
  defer log.SetLevel(log.InfoLevel)

  err := minikube.CheckStatus()
  if err != nil {
    err = command.Run("minikube start", "", nil, nil)
    if err != nil {
      return err
    }
  }

  err = command.Run("kubectl config use-context minikube", "", nil, nil)
  if err != nil {
    log.Errorf("Failed to switch kubectl context to 'minikube' due to the error: %s", err.Error())
    os.Exit(1)
  }
  return nil
}

func MinikubeStop() error {
  return command.Run("minikube stop", "", nil, nil)
}

func MinikubeDelete() error {
  return command.Run("minikube delete", "", nil, nil)
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
    return errors.Wrap(err, "deleting docker registry mock container")
  }

  cmd = fmt.Sprintf("docker run -it -d --rm -p %s:5000 --name %s %s", DockerMockRepo, DockerMockContainerName, DockerMockSrcImg)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "starting docker registry mock")
  }
  time.Sleep(5 * time.Second)

  cmd = fmt.Sprintf("echo '%s' | docker login -u %s --password-stdin %s", DockerMockTestPassword, DockerMockTestUser, DockerMockRepo)
  cmdHandler = exec.Command("sh", "-c", cmd)
  // cmdHandler.Stdout = os.Stdout
  cmdHandler.Stderr = os.Stderr
  err = cmdHandler.Run()
  if err != nil {
    return errors.Wrap(err, "logging in to the docker registry mock")
  }

  cmd = fmt.Sprintf("docker tag %s %s", DockerMockSrcImg, DockerMockTestImage)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "assigning tag for the testing image")
  }
  
  cmd = fmt.Sprintf("docker push %s", DockerMockTestImage)
  err = command.Run(cmd, "", nil, os.Stderr)
  if err != nil {
    return errors.Wrap(err, "pushing testing image to the mock docker registry")
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
      log.Error("Failed to stop docker registry mock due to the error: ", err)
  }
}
