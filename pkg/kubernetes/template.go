package kubernetes

import (
  "fmt"
  "io/ioutil"
  "path/filepath"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/schema/rootcfg"
)

// func ReadTemplates(path string) ([]byte, error) {
func ReadTemplates(cfg *rootcfg.CarbonConfig) ([]byte, error) {
  log.Debug("Reading kube manifest templates")
  fullPath := filepath.Join(cfg.Cwd, cfg.Data.KubeManifests)
  files, err := filepath.Glob(fullPath)
  if err != nil {
    return nil, err
  }
  // TODO: test this:
  if len(files) == 0 {
    return nil, fmt.Errorf("Unable to find Kubernetes manifests at ", fullPath)
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
