package kubernetes

import (
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/schema/carboncfg"
  "github.com/starofservice/carbon/pkg/schema/pkgmeta"
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

type KubeInstall struct {
  BuiltManifest []byte
  RawManifest []byte
  Scope string
  Variables DepVars
}

func NewKubeInstall(meta *pkgmeta.PackageConfig, ccfg *carboncfg.CarbonConfig, dname string, dtag string) (*KubeInstall, error) {
  log.Debug("Creating new Kubernetes inatllation handler")

  k8sManifest, err := base64.Decode(meta.Data.KubeConfigB64)
  if err != nil {
    return nil, err
  }

  scope, err := ccfg.CarbonScope()
  if err != nil {
    return nil, errors.Wrap(err, "getting Carbon package installation scope")
  }

  self := &KubeInstall{
    RawManifest: k8sManifest,
    Scope: scope,
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
