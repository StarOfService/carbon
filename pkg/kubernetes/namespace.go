package kubernetes

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"

  "github.com/Jeffail/gabs" // patch filters
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
  "github.com/starofservice/carbon/pkg/util/tojson"
)

func (self *KubeInstall) SetNamespace() error {
  log.Debug("Configuring namespace for Kubernetes manifests")
  if self.Scope == "cluster" {
    err := self.setNamespaceClusterScope()
    if err != nil {
      return errors.Wrap(err, "setting up namespace for Kubernetes resources")
    }
  } else if self.Scope == "namespace" {
    err := self.setNamespaceNamespaceScope()
    if err != nil {
      return errors.Wrap(err, "setting up namespace for Kubernetes resources")
    }
  }
  return nil
}

func setNamespacePatch() ([]byte, error) {
    ops := fmt.Sprintf(`---
filters:
  kind: .*
type: merge
patch:
  metadata:
    namespace: %s
`, kubecommon.CurrentNamespace)

  patch, err := tojson.ToJSON([]byte(ops))
  if err != nil {
    log.Error("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    return []byte{}, errors.Wrap(err, "converting Kubernetes patch with Carbon namespace to JSON")
  }

  return patch, nil
}

func (self *KubeInstall) setNamespaceNamespaceScope() error {
  log.Debug("Carbon uses Naemspaced scope")

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
    log.Trace("Original resource: ", string(original))

    jsonParsed, err := gabs.ParseJSON(original)
    if err != nil {
      return errors.Wrap(err, "parsing manifest with Gabs library")
    }

    mkind, _ := jsonParsed.Path("kind").Data().(string)
    mname, _ := jsonParsed.Path("metadata.name").Data().(string)
    mns, ok := jsonParsed.Path("metadata.namespace").Data().(string)
    if ok && mns != kubecommon.CurrentNamespace {
      log.Warnf("Namespace '%s' is defined for '%s' resource of '%s' kind. This namespace will be overriden by '%s' due to namespaced scope",
                mns, mname, mkind, kubecommon.CurrentNamespace)
    }
  }

  patch, err := setNamespacePatch()
  if err != nil {
    return err
  }

  if err := self.ProcessPatches(patch); err != nil {
    return err
  }
  return nil
}

func (self *KubeInstall) setNamespaceClusterScope() error {
  log.Debug("Carbon uses Clustered scope")
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
    log.Trace("Original resource: ", string(original))

    jsonParsed, err := gabs.ParseJSON(original)
    if err != nil {
      return errors.Wrap(err, "parsing manifest with Gabs library")
    }

    // var modified []byte
    ok := jsonParsed.ExistsP("metadata.namespace")
    if ok {
      log.Trace("Namespace is defined, nothing to do")
      newManifest = append(newManifest, original...)
    } else {
      log.Trace("Namespace isn't defined, setting up a namespace for the resource")

      patch, err := setNamespacePatch()
      if err != nil {
        return err
      }

      var ph Patcher
      err = json.Unmarshal(patch, &ph)
      if err != nil {
        log.Errorf("Document:\n%s", string(patch))
        return errors.Wrap(err, "deserializing patch data")
      }

      modified, err := ph.Apply(original)
      if err != nil {
        return err
      }
      log.Trace("Modified resource: ", string(modified))
      newManifest = append(newManifest, modified...)
    }
  }
  self.BuiltManifest = newManifest

  return nil
}
