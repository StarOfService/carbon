package kubernetes

import (
  "bytes"
  "encoding/json"
  "io"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/kubernetes/patcher"
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

func deserialPatchers(data []byte) ([]patcher.Patcher, error) {
  log.Debug("Deserialzing patch handlers")

  var resp []patcher.Patcher
  dec := json.NewDecoder(bytes.NewReader(data))
  for {
    var ph patcher.Patcher
    err := dec.Decode(&ph)
    if err == io.EOF {
        break
    }
    if err != nil {
      log.Fatalf("Document:\n%s", string(data))
      return nil, errors.Wrap(err, "deserializing patch data")
    }
    resp = append(resp, ph)
  }
  return resp, nil
}

func (self *KubeDeployment) processPatchers(phs []patcher.Patcher) error {
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
        log.Fatal("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
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
