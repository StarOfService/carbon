package kubernetes

import (
  "bytes"
  "fmt"
  // "io/ioutil"
  "path/filepath"
  "strings"
  "text/template"
  "github.com/pkg/errors"
  // "gopkg.in/yaml.v2"

  "github.com/starofservice/carbon/pkg/util/tojson"
  log "github.com/sirupsen/logrus"
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
 // tpl := template.Must(template.ParseGlob(path))
  tpl, err := template.ParseFiles(path)
  if err != nil {
    // return fmt.Errorf("Unable to parse kuberentese manifest teamplate %s due to the error: %s", path, err.Error())
    return errors.Wrapf(err, "parsing Kuberentese manifests template '%s'", path)
  }

  tpl.Option("missingkey=error")
  
  var data bytes.Buffer
  err = tpl.Execute(&data, self.Variables)
  // err = tpl.Execute(os.Stdout, self.Variables)
  if err != nil {
    if strings.Index(err.Error(), "no entry for key") != -1 || strings.Index(err.Error(), "can't evaluate field") != -1 {
      return fmt.Errorf("%s\nPlease use make sure that all variables are defined at carbon.yaml and use `.Var` prefix for the variables at kubernetes config files.", err.Error())
    } else {
      // return fmt.Errorf("Unable to execute kuberentese manifests teamplate due to the error: %s", err.Error())
      return errors.Wrap(err, "building Kuberentese manifests template")
    }
  }

  _, err = tojson.ToJson(data.Bytes())
  if err != nil {
    return errors.Wrap(err, "converting Kuberentese manifests to JSON")
  }
  
  // if err = self.VerifyYAML(path, data.Bytes()); err != nil {
  //   return err
  // }
  
  return nil
}

// func (self *KubeDeployment) VerifyYAML(path string, data []byte) error {
//   var obj interface{}
//   err := yaml.Unmarshal(data, &obj)
//   if err != nil {
//     return fmt.Errorf("Failed to parse YAML config `%s` due to the error: %s", path, err.Error())
//   }
//   return nil
// }