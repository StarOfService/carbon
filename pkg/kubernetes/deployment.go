package kubernetes

import (
  "bytes"
  "fmt"
  "text/template"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
  "github.com/starofservice/carbon/pkg/util/tojson"
  "github.com/starofservice/carbon/pkg/util/base64"
)

type DepVarsPkg struct {
  DockerName string
  DockerTag string
  Name string
  Version string
}

type DepVars struct {
  Pkg DepVarsPkg
  Var map[string]string
}

type KubeDeployment struct {
  BuiltManifest []byte
  RawManifest []byte
  Variables DepVars
}

func NewKubeDeployment(meta *pkgmeta.PackageConfig, dname string, dtag string) (*KubeDeployment, error) {
  log.Debug("Creating new kubernetes deployment handler")

  k8sManifest, err := base64.Decode(meta.Data.KubeConfigB64)
  if err != nil {
    return nil, err
  }

  self := &KubeDeployment{
    RawManifest: k8sManifest,
    Variables: DepVars{
      Pkg: DepVarsPkg{
        Name: meta.Data.PkgName,
        Version: meta.Data.PkgVersion,
        DockerName: dname,
        DockerTag: dtag,
      },
      Var: make(map[string]string),
    },
  }

  for _, v := range meta.Data.Variables {
    self.Variables.Var[v.Name] = v.Default
  }

  return self, nil
}

func (self *KubeDeployment) UpdateVars(vars map[string]string) {
  log.Debug("Applying carbon variables")

  for k, v := range vars {
    log.Tracef("%s: %s", k, v)
    if _, ok := self.Variables.Var[k]; ok {
      self.Variables.Var[k] = v  
    } else {
      log.Warnf("Variable '%s' is not supported by the current package", k)
    }
  }
}


func (self *KubeDeployment) Build() error {
  log.Debug("Building kubernetes manifest based on the template from the package and provided variables")

  tpl, err := template.New("kubeManifest").Option("missingkey=zero").Parse(string(self.RawManifest))
  if err != nil {
    return errors.Wrap(err, "parsing Kuberentese manifests teamplate")
  }

  buf := &bytes.Buffer{}
  err = tpl.Execute(buf, self.Variables)
  if err != nil {
    return errors.Wrap(err, "building Kuberentese manifests")
  }

  self.BuiltManifest, err = tojson.ToJSON(buf.Bytes())
  if err != nil {
    return errors.Wrap(err, "converting Kuberentese manifests to JSON")
  }

  return nil
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
  
  patch, err := tojson.ToJSON([]byte(ops))
  if err != nil {
    log.Error("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    return errors.Wrap(err, "converting Kubernetes patch with Carbon labels to JSON")
  }
  if err := self.ProcessPatches(patch); err != nil {
    return err
  }
  return nil
}
