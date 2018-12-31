package patcher

import (
  // "fmt"
  "os"
  "strings"
  "github.com/Jeffail/gabs" // patch filters
  jsonpatch "github.com/evanphx/json-patch" //RFC6902 and RFC7386
  // "bytes"
  // "io"
  "regexp"
  "encoding/json"
  // "github.com/starofservice/carbon/pkg/util/tojson"
  log "github.com/sirupsen/logrus"
)


type Patcher struct {
  Filters map[string]string `json:"filters"`
  Type string               `json:"type"`
  Patch json.RawMessage     `json:"patch"`
}

func (self *Patcher) Apply(original []byte) []byte {
  log.Trace("Processing patch for object")
  log.Tracef("Patch: %s", string(self.Patch))
  log.Tracef("Original object: %s", string(original))
  modified := original
  if !self.MatchObj(original) {
    log.Trace("The object doesn't match patch filters. Skipping it")
    return modified
  }

  var err error
  switch strings.ToLower(self.Type) {
  case "merge":
    // self.patchApplyMerge(p)
    // var patch json.RawMessage
    // err = json.Unmarshal(self.Patch, &patch)
    // if err != nil {
    //   // return nil, err
    //   log.Fatal("Failed to deserialize patch data due to the error: %s", err.Error())
    //   log.Fatal("Failed to deserialize patch data due to the error: %s", err.Error())
    //   log.Fatal("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    //   os.Exit(1)
    // }
    // for _, i := range patches {
    modified, err = jsonpatch.MergePatch(original, self.Patch)
    if err != nil {
      // return nil, err
      log.Fatalf("Failed to apply merge patch for kubernetes manifest due to the error: %s", err.Error())
      log.Fatalf("Patch data: %s", string(self.Patch))
      log.Fatalf("Kubernetes manifest data: %s", original)
      os.Exit(1)
    }
    // }
  case "json":
    // self.patchApplyJson(p)
    jp, err := jsonpatch.DecodePatch(self.Patch)
    if err != nil {
      // panic(err)
      // return nil, err
      log.Fatalf("Failed to decode json patch data due to the error: %s", err.Error())
      log.Fatalf("Patch data: %s", string(self.Patch))
      os.Exit(1)
    }
    modified, err = jp.Apply(original)
    if err != nil {
      // panic(err)
      // return nil, err
      log.Fatalf("Failed to apply json patch for kubernetes manifest due to the error: %s", err.Error())
      log.Fatalf("Patch data: %s", string(self.Patch))
      log.Fatalf("Kubernetes manifest data: %s", original)
      os.Exit(1)
    }
  default:
    // TODO ERROR
    // fmt.Println("Unknown patch type: ", self.Type)
    // return fmt.Errorf("Unknown patch type: %s", self.Type)
    log.Fatalf("Unknown patch type: %s", self.Type)
    os.Exit(1)
  }
  log.Trace("Modified object: %s", string(modified))
  return modified
}


func (self *Patcher) MatchObj(data []byte) bool {
  gjp, err := gabs.ParseJSON(data)
  if err != nil {
    log.Fatalf("Failed to parse JSON data: %s", data)
    log.Fatalf("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    os.Exit(1)
    // return false
  }

  for k, f := range self.Filters {
    path := strings.Replace(k, "/", ".", -1)
    if !gjp.ExistsP(path) { return false }
    d := gjp.Path(path).String()
    re := regexp.MustCompile(f)
    if !re.MatchString(d) { return false }
  }
  return true
}


