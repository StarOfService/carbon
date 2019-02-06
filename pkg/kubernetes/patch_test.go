package kubernetes_test

import (
  "testing"

  "github.com/stretchr/testify/assert"

  "github.com/starofservice/carbon/pkg/kubernetes"
)

func TestKubePatch(t *testing.T) {
  patchSuites := []struct {
    description string
    original    string
    patch       string
    expected    string
  }{
    {
      "(RFC7386) Merge type: add new field for all original objects",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":".*"},"type":"merge","patch":{"metadata":"new"}}`,
      `{"kind":"pod","metadata":"new","name":"podA"}{"kind":"pod","metadata":"new","name":"podB"}`,
    },
    {
      "(RFC7386) Merge type: add new field with intermediate keys for all original objects",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":".*"},"type":"merge","patch":{"metadata":{"label":"new"}}}`,
      `{"kind":"pod","metadata":{"label":"new"},"name":"podA"}{"kind":"pod","metadata":{"label":"new"},"name":"podB"}`,
    },
    {
      "(RFC7386) Merge type: override existing items for all original objects",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":".*"},"type":"merge","patch":{"kind":"new"}}`,
      `{"kind":"new","name":"podA"}{"kind":"new","name":"podB"}`,
    },
    {
      "(RFC7386) Merge type: override existing items for a specific object",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":"podA"},"type":"merge","patch":{"kind":"new"}}`,
      `{"kind":"new","name":"podA"}{"kind":"pod","name":"podB"}`,
    },
    {
      "(RFC6902) JSON type: add new field for all original objects",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":".*"}, "type":"json", "patch": [{"op": "add", "path": "/metadata", "value": "new"}]}`,
      `{"kind":"pod","metadata":"new","name":"podA"}{"kind":"pod","metadata":"new","name":"podB"}`,
    },
    {
      "(RFC6902) JSON type: remove field for a specific objects",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":"podB"}, "type":"json", "patch": [{"op": "remove", "path": "/kind"}]}`,
      `{"kind":"pod","name":"podA"}{"name":"podB"}`,
    },
    {
      "(RFC6902 & RFC7386) JSON and Merge types: combination of patches",
      `{"kind":"pod","name":"podA"}{"kind":"pod","name":"podB"}`,
      `{"filters":{"name":"podB"}, "type":"json", "patch": [{"op": "remove", "path": "/kind"}]}{"filters":{"name":".*"},"type":"merge","patch":{"metadata":"new"}}`,
      `{"kind":"pod","metadata":"new","name":"podA"}{"metadata":"new","name":"podB"}`,
    },
  }

  for _, s := range patchSuites {
    t.Log("suite:", s.description)

    kd := &kubernetes.KubeInstall{
      BuiltManifest: []byte(s.original),
    }
    err := kd.ProcessPatches([]byte(s.patch))
    if err != nil {
      t.Errorf(err.Error())
    }

    assert.Equal(t, s.expected, string(kd.BuiltManifest), "they should be equal")

    // if s.assert != string(kd.BuiltManifest) {
    //   t.Errorf("Test suite object %v doesn't match to the generated data %v", s.assert, string(kd.BuiltManifest))
    // }
  }
}
