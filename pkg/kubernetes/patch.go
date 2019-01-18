package kubernetes

import (
  // "fmt"
  // "strings"
  // "github.com/Jeffail/gabs" // patch filters
  // jsonpatch "github.com/evanphx/json-patch" //RFC6902 and RFC7386
  "bytes"
  "io"
  // "os"
  "github.com/pkg/errors"
  // "regexp"
  "encoding/json"
  // "github.com/starofservice/carbon/pkg/util/tojson"
  log "github.com/sirupsen/logrus"

  "github.com/starofservice/carbon/pkg/kubernetes/patcher"
)


func (self *KubeDeployment) ProcessPatches(data []byte) error {
  log.Debug("Processing patches")

  // normData := normalizePatchesData(data)
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

// func normalizePatchesData(patches [][]byte) []byte {
//   log.Debug("Normalizing patches data")
//   var resp []byte
//   for _, i := range patches {
//     ni := tojson.ToJson(i)
//     resp = append(resp, ni...)
//   }
//   return resp
// }

func deserialPatchers(data []byte) ([]patcher.Patcher, error) {
  log.Debug("Deserialzing patch handlers")
  // fmt.Println(string(data))
  var resp []patcher.Patcher
  dec := json.NewDecoder(bytes.NewReader(data))
  for {
    var ph patcher.Patcher
    err := dec.Decode(&ph)
    if err == io.EOF {
        break
    }
    if err != nil {
      // TODO DEBUG
      // log.Fatalf("Failed to deserialize patch data due to the error: %s", err.Error())
      log.Fatalf("Document:\n%s", string(data))
      return nil, errors.Wrap(err, "deserializing patch data")
      // os.Exit(1)

      // fmt.Println("Failed to unmarshal Patches data due to the error: ", err.Error())
      // return nil, err
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
        // TODO ERROR
        // fmt.Println("Failed to unmarshal JSON data due to the error: ", err.Error())
        // return err
        // log.Fatalf("Failed to deserialize kubernetes manifest data due to the error: %s", err.Error())
        // os.Exit(1)
        return errors.Wrap(err, "deserializing Kubernetes manifest")
      }
      original, err := json.Marshal(obj)
      if err != nil {
        // log.Fatalf("Failed to serialize kubernetes manifest data due to the error: %s", err.Error())
        log.Fatal("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
        // os.Exit(1)
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


