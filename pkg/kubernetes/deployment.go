package kubernetes

import (
  "bytes"
  "fmt"
  // "io/ioutil"
  // "path/filepath"
  // "strings"
  "text/template"
  "os"
  // "time"
  // "regexp"
  // "encoding/base64"
  "github.com/starofservice/carbon/pkg/util/tojson"
  "github.com/starofservice/carbon/pkg/util/base64"
  pkgmetalatest "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
  log "github.com/sirupsen/logrus"
  // "github.com/starofservice/carbon/pkg/util/tojson"
  // sigsk8syaml "sigs.k8s.io/yaml"
)

type DepVarsPkg struct {
  Name string
  Version string
  DockerName string
  DockerTag string
}

type DepVars struct {
  Var map[string]string
  Pkg DepVarsPkg
}

type KubeDeployment struct {
  BuiltManifest []byte
  RawManifest []byte
  Variables DepVars
}

func NewKubeDeployment(meta *pkgmetalatest.PackageConfig, dname string, dtag string) (*KubeDeployment, error) {
  log.Debug("Creating new kubernetes deployment handler")

  k8sManifest, err := base64.Decode(meta.KubeConfigB64)
  if err != nil {
    return nil, err
  }

  self := &KubeDeployment{
    RawManifest: k8sManifest,
    Variables: DepVars{
      Pkg: DepVarsPkg{
        Name: meta.Name,
        Version: meta.Version,
        DockerName: dname,
        DockerTag: dtag,
      }, //make(map[string]string),
      Var: make(map[string]string),
    },
  }

  for _, v := range meta.Variables {
    self.Variables.Var[v.Name] = v.Default
  }

  return self, nil
}

func (self *KubeDeployment) UpdateVars(vars map[string]string) {
  log.Debug("Applying carbon variables")

  for k, v := range vars {
    log.Tracef("%s: %s", k, v)
    self.Variables.Var[k] = v
  }
}


func (self *KubeDeployment) Build() {
  log.Debug("Building kubernetes manifest based on the template from the package and provided variables")

  tpl, err := template.New("kubeManifest").Option("missingkey=zero").Parse(string(self.RawManifest))
  if err != nil {
    log.Fatalf("Failed to parse kuberentese manifests teamplate due to the error: %s", err.Error())
    os.Exit(1)
  }

  buf := &bytes.Buffer{}
  err = tpl.Execute(buf, self.Variables)
  if err != nil {
    log.Fatalf("Failed to build kuberentese manifests teamplate due to the error: %s", err.Error())
    os.Exit(1)
  }

  self.BuiltManifest = tojson.ToJson(buf.Bytes())
}

func (self *KubeDeployment) SetAppLabel() error {
  log.Debug("Applying carbon lables for kubernetes manifests")
  ops := fmt.Sprintf(`---
filters:
  kind: .*
type: merge
patch:
  metadata:
    labels:
      managed-by: carbon
      carbon/component-name: %s
      carbon/component-version: %s
`, self.Variables.Pkg.Name, self.Variables.Pkg.Version)
  
  // patch := [][]byte{[]byte(ops)}
  patch := tojson.ToJson([]byte(ops))
  if err := self.ProcessPatches(patch); err != nil {
    return err
  }
  return nil
}