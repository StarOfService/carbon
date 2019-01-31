package cmd_test

import (
  // "bytes"
  "os"
  "testing"
  log "github.com/sirupsen/logrus"

  carboncmd "github.com/starofservice/carbon/cmd"
  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/test"
)

const (
  lifecyclePkgName = "carbon-test"
  lifecyclePkgTag = "test"
)
var lifecyclePkgFull = lifecyclePkgName + ":" + lifecyclePkgTag

func TestMain(m *testing.M) {
    log.Info("Starting minikube")
    err := test.MinikubeStart()
    if err != nil {
      log.Error("Failed to start minikube due ot the error: ", err.Error())
    }

    code := m.Run()
    
    minikube.UnsetDockerEnv()
    test.DeferMinikubeDelete()
    
    os.Exit(code)
}

func TestBuildAndPush(t *testing.T) {
  minikube.UnsetDockerEnv()

  t.Log("Starting docker registry")
  err := test.DockerRegistryStart()
  if err != nil {
    test.DockerRegistryDelete()
    // log.Fatal("Failed to start docker registry due ot the error: ", err.Error())
    t.Error("Failed to start docker registry due ot the error: ",err.Error())
  }
  defer test.DockerRegistryDelete()

  carboncmd.RootCmd.SetArgs([]string{
    "build",
    "-l", "error",
    "-c", "../test/carbon.yaml",
    "--tag", test.DockerMockRepo+"/carbon-test",
    "--tag-prefix", "pref-",
    "--tag-suffix", "-suf",
    "--push",
    "--rm",
  })
  _, err = carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyBuild(t *testing.T) {
  carboncmd.BuildTags = []string{}

  carboncmd.RootCmd.SetArgs([]string{
    "build",
    "-l", "error",
    "-m",
    "-c", "../test/carbon.yaml",
    "--tag", lifecyclePkgFull,
    "--tag-prefix", "",
    "--tag-suffix", "",
    "--push=false",
    "--rm=false",
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyInspect(t *testing.T) {
  carboncmd.RootCmd.SetArgs([]string{
    "inspect",
    "-l", "error",
    "-m",
    lifecyclePkgFull,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyDeploy(t *testing.T) {
  carboncmd.RootCmd.SetArgs([]string{
    "deploy",
    "-l", "error",
    "-m",
    lifecyclePkgFull,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyStatusAll(t *testing.T) {
  carboncmd.RootCmd.SetArgs([]string{
    "status",
    "-l", "error",
    "-m",
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyStatusSuite(t *testing.T) {
  carboncmd.RootCmd.SetArgs([]string{
    "status",
    "-l", "error",
    "-m",
    lifecyclePkgName,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyStatusSuiteFull(t *testing.T) {
  carboncmd.RootCmd.SetArgs([]string{
    "status",
    "-l", "error",
    "-m",
    "-f",
    lifecyclePkgName,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestVersion(t *testing.T) {
  carboncmd.RootCmd.SetArgs([]string{
    "version",
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}
