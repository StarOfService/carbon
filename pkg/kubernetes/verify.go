package kubernetes

import (
  "bytes"
  "path/filepath"
  "strings"
  "text/template"

  "github.com/Masterminds/sprig"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/util/tojson"
)

func (self *KubeInstall) VerifyAll(path string) error {
  log.Debug("Verifying Kubernetes templates")
  files, err := filepath.Glob(path)
  if err != nil {
    return err
  }

  for _, f := range files {
    if err = self.VerifyTpl(f); err != nil {
      return err
    }
  }

  return nil
}

func (self *KubeInstall) VerifyTpl(path string) error {
  log.Tracef("Verifying %s", path)
  tpl, err := template.New(filepath.Base(path)).Option("missingkey=zero").Funcs(sprig.TxtFuncMap()).ParseFiles(path)
  if err != nil {
    return errors.Wrapf(err, "parsing Kubernetes manifests template '%s'", path)
  }

  var data bytes.Buffer
  err = tpl.Execute(&data, self.Variables)
  if err != nil {
    if strings.Index(err.Error(), "no entry for key") != -1 || strings.Index(err.Error(), "can't evaluate field") != -1 {
      return errors.Errorf("%s\nPlease make sure that all variables are defined at carbon.yaml and prefixed by `.Var` at Kubernetes manifest templates", err.Error())
    }
    return errors.Wrap(err, "building Kubernetes manifests template")
  }

  _, err = tojson.ToJSON(data.Bytes())
  if err != nil {
    return errors.Wrap(err, "converting Kubernetes manifests to JSON")
  }

  return nil
}
