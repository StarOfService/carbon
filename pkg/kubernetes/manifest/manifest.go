package manifest

import (
  "fmt"
  "io/ioutil"
  "path/filepath"
  // "strings"
  "gopkg.in/yaml.v2"
)

func ReadTemplates(path string) ([]byte, error) {
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
    return fmt.Errorf("Unable to read file `%s` due to the error: %s", path, err.Error())
  }
  // TODO sanity check for the data
  var data interface{}
  err = yaml.Unmarshal(cfg, &data)
  if err != nil {
    return fmt.Errorf("Failed to parse YAML config `%s` due to the error: %s", path, err.Error())
  }

  if len(cfg) > 0 {
    *r = append(*r, "\n---\n"...)
    *r = append(*r, cfg...)
  }

  return nil
}
