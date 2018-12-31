package kubernetes

import (
  "io/ioutil"
  "path/filepath"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/schema/pkgcfg"
)

func ReadTemplates(cfg *pkgcfg.CarbonConfig) ([]byte, error) {
  log.Debug("Reading kube manifest templates")
  fullPath := filepath.Join(cfg.Cwd, cfg.Data.KubeManifests)
  files, err := filepath.Glob(fullPath)
  if err != nil {
    return nil, err
  }

  if len(files) == 0 {
    return nil, errors.Errorf("Unable to find Kubernetes manifests at %s", fullPath)
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
    return errors.Wrapf(err, "reading '%s", path)
  }

  if len(cfg) > 0 {
    *r = append(*r, "\n---\n"...)
    *r = append(*r, cfg...)
  }

  return nil
}
