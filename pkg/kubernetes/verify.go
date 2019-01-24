package kubernetes

import (
  "bytes"
  "fmt"
  "path/filepath"
  "strings"
  "text/template"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/util/tojson"
)

func (self *KubeDeployment) VerifyAll(path string) error {
  log.Debug("Verifying kubernetes templates")
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

func (self *KubeDeployment) VerifyTpl(path string) error {
  log.Tracef("Verifying %s", path)
  tpl, err := template.ParseFiles(path)
  if err != nil {
    return errors.Wrapf(err, "parsing Kuberentese manifests template '%s'", path)
  }

  tpl.Option("missingkey=error")
  
  var data bytes.Buffer
  err = tpl.Execute(&data, self.Variables)
  if err != nil {
    if strings.Index(err.Error(), "no entry for key") != -1 || strings.Index(err.Error(), "can't evaluate field") != -1 {
      return fmt.Errorf("%s\nPlease make sure that all variables are defined at carbon.yaml and use `.Var` prefix for the variables at kubernetes config files", err.Error())
    }
    return errors.Wrap(err, "building Kuberentese manifests template")
  }

  _, err = tojson.ToJSON(data.Bytes())
  if err != nil {
    return errors.Wrap(err, "converting Kuberentese manifests to JSON")
  }
  
  return nil
}
