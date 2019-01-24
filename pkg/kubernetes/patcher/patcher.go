package patcher

import (
  "encoding/json"
  "regexp"
  "strings"

  jsonpatch "github.com/evanphx/json-patch" // RFC6902 and RFC7386
  "github.com/Jeffail/gabs" // patch filters
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
)

type Patcher struct {
  Filters map[string]string `json:"filters"`
  Type string               `json:"type"`
  Patch json.RawMessage     `json:"patch"`
}

func (self *Patcher) Apply(original []byte) ([]byte, error) {
  log.Trace("Processing patch for object")
  log.Tracef("Patch: %s", string(self.Patch))
  log.Tracef("Original object: %s", string(original))
  modified := original
  
  match, err := self.MatchObj(original)
  if err != nil {
    return original, err
  }
  if !match {
    log.Trace("The object doesn't match patch filters. Skipping it")
    return modified, nil
  }

  switch strings.ToLower(self.Type) {
  case "merge":
    modified, err = jsonpatch.MergePatch(original, self.Patch)
    if err != nil {
      log.Fatalf("Patch data: %s", string(self.Patch))
      log.Fatalf("Kubernetes manifest data: %s", original)
      return original, errors.Wrap(err, "applying merge patch for Kuberentese manifest")
    }
  case "json":
    jp, err := jsonpatch.DecodePatch(self.Patch)
    if err != nil {
      log.Fatalf("Patch data: %s", string(self.Patch))
      return original, errors.Wrap(err, "decoding json patch")
    }
    modified, err = jp.Apply(original)
    if err != nil {
      log.Fatalf("Patch data: %s", string(self.Patch))
      log.Fatalf("Kubernetes manifest data: %s", original)
      return original, errors.Wrap(err, "applying json patch for Kubernetes manifest")
    }
  default:
    log.Fatalf("Patch data: %s", string(self.Patch))
    return original, errors.Errorf("Unknown patch type: %s", self.Type)
  }
  log.Trace("Modified object: %s", string(modified))
  return modified, nil
}

func (self *Patcher) MatchObj(data []byte) (bool, error) {
  gjp, err := gabs.ParseJSON(data)
  if err != nil {
    log.Fatalf("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    log.Fatalf("Kubernetes manifest data: %s", data)
    return false, errors.Wrap(err, "parsing Kubernetes manifest JSON data")
  }

  for k, f := range self.Filters {
    path := strings.Replace(k, "/", ".", -1)
    if !gjp.ExistsP(path) { return false, nil }
    d := gjp.Path(path).String()
    re := regexp.MustCompile(f)
    if !re.MatchString(d) { return false, nil }
  }
  return true, nil
}
