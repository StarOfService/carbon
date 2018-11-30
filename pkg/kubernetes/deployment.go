package kubernetes

import (
  "fmt"
  "io/ioutil"
  "text/template"
)

type DeployVars struct {
  Var map[string]string
  Pkg map[string]string
}

type KubeDeploy struct {
  Manifest []byte
  Variables DeployVars
}

func (self *KubeDeploy) Verify() error {
  tpl, err := template.New("kubeManifest").Option("missingkey=error").Parse(string(self.Manifest))
  if err != nil {
    return fmt.Errorf("Unable to parse kuberentese manifest teamplate due to the error: %s", err.Error())
  }
  err = tpl.Execute(ioutil.Discard, self.Variables)
  if err != nil {
    return fmt.Errorf("Unable to execute kuberentese manifest teamplate due to the error: %s", err.Error())
  }
  return nil
}