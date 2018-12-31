package tojson

import (
  "fmt"
  "testing"
  // jsonpatch "github.com/evanphx/json-patch"
  // ghodssyaml "github.com/ghodss/yaml"
)

func TestSingleDocJson(t *testing.T) {
  resp := ToJson([]byte(`{"k1": "v1"}`))
  fmt.Println(string(resp))
}
func TestSingleMultiJson(t *testing.T) {
  // fmt.Println(string(ToJson([]byte(`{"k1": "v1"}{"k2": "v2"}`))))

  resp := ToJson([]byte(`{"k1": "v1"}{"k2": "v2"}`))
  fmt.Println(string(resp))
}
func TestSingleDocYaml(t *testing.T) {
  // fmt.Println(string(ToJson([]byte(`k1: v1`))))

  resp := ToJson([]byte(`k1: v1`))
  fmt.Println(string(resp))
}
func TestSingleMultiYaml(t *testing.T) {
//   fmt.Println(string(ToJson([]byte(`k1: v1
// ---
// k2: v2`))))


  resp := ToJson([]byte(`k1: v1
---
k2: v2`))
  fmt.Println(string(resp))
}
func TestUnknownFormat(t *testing.T) {
  // fmt.Println(string(ToJson([]byte(`k1v1
  //   asdf`))))


  resp := ToJson([]byte(`k1v1
    asdf`))
  fmt.Println(string(resp))
}