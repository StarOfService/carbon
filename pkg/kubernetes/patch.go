package kubernetes

import (
  "bytes"
  "encoding/json"
  "io"
  "regexp"
  "strings"

  jsonpatch "github.com/evanphx/json-patch" // RFC6902 and RFC7386
  "github.com/Jeffail/gabs" // patch filters
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
)

func (self *KubeDeployment) ProcessPatches(data []byte) error {
  log.Debug("Processing patches")

  patches, err := deserialPatchers(data)
  if err != nil {
    return err
  }

  err = self.processPatchers(patches)
  if err != nil {
    return err
  }
  return nil  
}

func deserialPatchers(data []byte) ([]Patcher, error) {
  log.Debug("Deserialzing patch handlers")

  var resp []Patcher
  dec := json.NewDecoder(bytes.NewReader(data))
  for {
    var ph Patcher
    err := dec.Decode(&ph)
    if err == io.EOF {
        break
    }
    if err != nil {
      log.Errorf("Document:\n%s", string(data))
      return nil, errors.Wrap(err, "deserializing patch data")
    }
    resp = append(resp, ph)
  }
  return resp, nil
}

func (self *KubeDeployment) processPatchers(phs []Patcher) error {
  log.Debug("Processing patch handlers")
  for _, ph := range phs {
    var newManifest []byte

    dec := json.NewDecoder(bytes.NewReader(self.BuiltManifest))
    for {
      var obj interface{}
      err := dec.Decode(&obj)
      if err == io.EOF {
          break
      }
      if err != nil {
        return errors.Wrap(err, "deserializing Kubernetes manifest")
      }
      original, err := json.Marshal(obj)
      if err != nil {
        log.Error("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
        return errors.Wrap(err, "serialize Kubernetes manifest")
      }
      modified, err := ph.Apply(original)
      if err != nil {
        return err
      }
      newManifest = append(newManifest, modified...)
    }
    self.BuiltManifest = newManifest
  }

  return nil
}

type Patcher struct {
  Filters map[string]string `json:"filters"`
  Type string               `json:"type"`
  Patch json.RawMessage     `json:"patch"`
}

func (self *Patcher) Apply(original []byte) ([]byte, error) {
  log.Trace("Processing patch for object")
  log.Trace("Patch: ", string(self.Patch))
  log.Trace("Original object: ", string(original))
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
      log.Error("Patch data: ", string(self.Patch))
      log.Error("Kubernetes manifest data: ", original)
      return original, errors.Wrap(err, "applying merge patch for Kuberentese manifest")
    }
  case "json":
    jp, err := jsonpatch.DecodePatch(self.Patch)
    if err != nil {
      log.Error("Patch data: ", string(self.Patch))
      return original, errors.Wrap(err, "decoding json patch")
    }
    modified, err = jp.Apply(original)
    if err != nil {
      log.Error("Patch data: ", string(self.Patch))
      log.Error("Kubernetes manifest data: ", original)
      return original, errors.Wrap(err, "applying json patch for Kubernetes manifest")
    }
  default:
    log.Error("Patch data: ", string(self.Patch))
    return original, errors.Error("Unknown patch type: ", self.Type)
  }
  log.Trace("Modified object: ", string(modified))
  return modified, nil
}

func (self *Patcher) MatchObj(data []byte) (bool, error) {
  gjp, err := gabs.ParseJSON(data)
  if err != nil {
    log.Error("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    log.Error("Kubernetes manifest data: ", data)
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
