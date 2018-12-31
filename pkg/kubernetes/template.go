package kubernetes

import (
  "fmt"
  "io/ioutil"
  "path/filepath"
  // "strings"
  // "gopkg.in/yaml.v2"
  log "github.com/sirupsen/logrus"
)

func ReadTemplates(path string) ([]byte, error) {
  log.Debug("Reading kube manifest templates")
  files, err := filepath.Glob(path)
  if err != nil {
    return nil, err
  }

  resp := []byte{}
  for _, f := range files {
    err = addKubeConfig(&resp, f)
    if err != nil {
      return nil, err
    }    
  }
  return resp, nil
}

func addKubeConfig(r *[]byte, path string) error {
  cfg, err := ioutil.ReadFile(path)
  if err != nil {
    return fmt.Errorf("Got an error when was trying to read `%s`. Error message: %s", path, err.Error())
  }

  if len(cfg) > 0 {
    *r = append(*r, "\n---\n"...)
    *r = append(*r, cfg...)
  }

  return nil
}
