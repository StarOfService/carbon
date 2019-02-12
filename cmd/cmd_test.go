package cmd_test

import (
  "io/ioutil"
  "os"
  "strings"
  "testing"
  
  log "github.com/sirupsen/logrus"

  carboncmd "github.com/starofservice/carbon/cmd"
  "github.com/starofservice/carbon/pkg/minikube"
  rclatest "github.com/starofservice/carbon/pkg/schema/pkgcfg/latest"
  "github.com/starofservice/carbon/pkg/test"
  "github.com/starofservice/carbon/pkg/version"
)

const (
  lifecyclePkgName = "carbon-test"
  lifecyclePkgTag = "test"
)
var lifecyclePkgFull = lifecyclePkgName + ":" + lifecyclePkgTag

func TestMain(m *testing.M) {
    log.Info("Starting Minikube")
    err := test.MinikubeStart()
    if err != nil {
      log.Error("Failed to start Minikube due ot the error: ", err.Error())
    }

    code := m.Run()

    log.SetLevel(log.InfoLevel)
    minikube.UnsetDockerEnv()
    // err = test.MinikubeDelete()
    // if err != nil {
    //   log.Error("Failed to delete Minikube due ot the error: ", err.Error())
    // }

    os.Exit(code)
}

func TestBuildAndPush(t *testing.T) {
  minikube.UnsetDockerEnv()

  t.Log("Starting Docker registry")
  err := test.DockerRegistryStart()
  if err != nil {
    test.DockerRegistryDelete()
    t.Error("Failed to start Doocker registry due ot the error: ",err.Error())
  }
  defer test.DockerRegistryDelete()

  realStdout := os.Stdout
  defer func() { os.Stdout = realStdout }()
  os.Stdout, _ = os.Open(os.DevNull)

  carboncmd.RootCmd.SetArgs([]string{
    "build",
    "-l", "fatal",
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

  realStdout := os.Stdout
  defer func() { os.Stdout = realStdout }()
  os.Stdout, _ = os.Open(os.DevNull)

  carboncmd.RootCmd.SetArgs([]string{
    "build",
    "-l", "fatal",
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
  origStdout := os.Stdout
  r, w, _ := os.Pipe()
  os.Stdout = w

  carboncmd.RootCmd.SetArgs([]string{
    "inspect",
    "-l", "fatal",
    "-m",
    lifecyclePkgFull,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error("Failed to run carbon inspect: ", err.Error())
    return
  }

  w.Close()
  out, _ := ioutil.ReadAll(r)
  os.Stdout = origStdout

  if !strings.Contains(string(out), lifecyclePkgName) {
    t.Error("Command output doesn't contain information about the package")
    t.Error("Output:\n", string(out))
  }
}

func TestLifecyclyInstall(t *testing.T) {
  realStdout := os.Stdout
  defer func() { os.Stdout = realStdout }()
  os.Stdout, _ = os.Open(os.DevNull)

  carboncmd.RootCmd.SetArgs([]string{
    "install",
    "-l", "fatal",
    "-m",
    lifecyclePkgFull,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestLifecyclyInstallAndRevert(t *testing.T) {
  realStdout := os.Stdout
  defer func() { os.Stdout = realStdout }()
  os.Stdout, _ = os.Open(os.DevNull)

  carboncmd.RootCmd.SetArgs([]string{
    "install",
    "-l", "fatal",
    "-m",
    "--patch", `{"filters":{"kind":".*"},"type":"merge","patch":{"kind":"foo"}}`,
    lifecyclePkgFull,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if !strings.Contains(err.Error(), "Revert has been applied successfully") {
    t.Error(err.Error())
  }
}

func TestLifecyclyStatusAll(t *testing.T) {
  origStdout := os.Stdout
  r, w, _ := os.Pipe()
  os.Stdout = w

  carboncmd.RootCmd.SetArgs([]string{
    "status",
    "-l", "fatal",
    "-m",
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }

  w.Close()
  out, _ := ioutil.ReadAll(r)
  os.Stdout = origStdout

  if !strings.Contains(string(out), lifecyclePkgFull) {
    t.Error("Command output doesn't contain information about the package")
    t.Error("Output:\n", string(out))
  }
}

func TestLifecyclyStatusSuite(t *testing.T) {
  origStdout := os.Stdout
  r, w, _ := os.Pipe()
  os.Stdout = w

  carboncmd.RootCmd.SetArgs([]string{
    "status",
    "-l", "fatal",
    "-m",
    lifecyclePkgName,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }

  w.Close()
  out, _ := ioutil.ReadAll(r)
  os.Stdout = origStdout

  if !strings.Contains(string(out), lifecyclePkgFull) {
    t.Error("Command output doesn't contain information about the package")
    t.Error("Output:\n", string(out))
  }
}

func TestLifecyclyStatusSuiteFull(t *testing.T) {
  origStdout := os.Stdout
  r, w, _ := os.Pipe()
  os.Stdout = w

  carboncmd.RootCmd.SetArgs([]string{
    "status",
    "-l", "fatal",
    "-m",
    "-f",
    lifecyclePkgName,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }

  w.Close()
  out, _ := ioutil.ReadAll(r)
  os.Stdout = origStdout

  if !strings.Contains(string(out), lifecyclePkgFull) {
    t.Error("Command output doesn't contain information about the package")
    t.Error("Output:\n", string(out))
  }
}

func TestLifecyclyDelete(t *testing.T) {
  realStdout := os.Stdout
  defer func() { os.Stdout = realStdout }()
  os.Stdout, _ = os.Open(os.DevNull)

  carboncmd.RootCmd.SetArgs([]string{
    "delete",
    "-l", "fatal",
    "-m",
    lifecyclePkgName,
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }
}

func TestVersion(t *testing.T) {
  origStdout := os.Stdout
  r, w, _ := os.Pipe()
  os.Stdout = w

  carboncmd.RootCmd.SetArgs([]string{
    "version",
  })
  _, err := carboncmd.RootCmd.ExecuteC()
  if err != nil {
    t.Error(err.Error())
  }

  w.Close()
  out, _ := ioutil.ReadAll(r)
  os.Stdout = origStdout

  if !strings.Contains(string(out), version.DefaultVersion) || !strings.Contains(string(out), rclatest.Version) {
    t.Error("Command output doesn't contain information about the package")
    t.Error("Output:\n", string(out))
  }
  log.Info("The endzzz")
}
