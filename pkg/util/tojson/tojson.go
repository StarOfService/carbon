package tojson

import (
  "bytes"
  "fmt"
  "io"
  "os"
  "encoding/json"
  "gopkg.in/yaml.v2"
  ghodssyaml "github.com/ghodss/yaml"
  log "github.com/sirupsen/logrus"
)

const (
  jsonPrefix = "{"
  yamlSep = "---"
)

func ToJson(data []byte) []byte {
  log.Debug("Normalizing config data")
  log.Tracef("Input data: %s", string(data))

  if err := checkJson(data); err == nil {
    return data
  }

  if err := checkYaml(data); err != nil {
    log.Fatal("Failed to verify YAML document due to the error: ", err.Error())
    log.Fatalf("Document:\n%s", data)
    os.Exit(1)
  }

  log.Debug("Converting YAML to JSON")
  yamlDocs := bytes.Split(data, []byte(yamlSep))
  var jsonDocs [][]byte
  for _, d := range yamlDocs {
    d = bytes.TrimSpace(d)
    if len(d) == 0 {
      continue
    }
    jd, err := ghodssyaml.YAMLToJSON(d)
    if err != nil {
      log.Fatal("Failed to convert YAML document to JSON format due to the error: ", err.Error())
      log.Fatalf("Document:\n%s", d)
      os.Exit(1)
    }
    jsonDocs = append(jsonDocs, jd)
  }

  resp := bytes.Join(jsonDocs, []byte(""))

  // Checking response format because plain text may be considered as a valid YAML format
  if err := checkJson(resp); err != nil {
    log.Fatal("Failed to verify JSON document after the YAML->JSON convertion due to the error: ", err.Error())
    log.Fatalf("Document:\n%s", data)
    os.Exit(1)
  }

  log.Trace("Output data: %s", string(resp))
  return resp
}

// json.Valid considers multi-document as a wrong JSON format.
// That's why we need a custom function.
func checkJson(data []byte) error {
  log.Debug("Trying to verify JSON document")
  trimData := bytes.TrimSpace(data)
  if !bytes.HasPrefix(trimData, []byte(jsonPrefix)) {
    return fmt.Errorf("Wrong JSON format")
  }

  dec := json.NewDecoder(bytes.NewReader(trimData))
  for {
    var d interface{}
    err := dec.Decode(&d)
    if err == io.EOF {
        return nil
    }
    if err != nil {
      // TODO DEBUG
      log.Debugf("Failed to unmarshal JSON data due to the error: ", err.Error())
      return err
    }
  }
}

func checkYaml(data []byte) error {
  log.Debug("Trying to verify YAML document")
  dec := yaml.NewDecoder(bytes.NewReader(data))
  for {
    var d interface{}
    err := dec.Decode(&d)
    if err == io.EOF {
        return nil
    }
    if err != nil {
      // TODO DEBUG
      log.Debugf("Failed to unmarshal YAML data due to the error: ", err.Error())
      
      return err
    }
  }
}