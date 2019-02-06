package tojson_test

import (
  "testing"

  "github.com/stretchr/testify/assert"

  "github.com/starofservice/carbon/pkg/util/tojson"
)

func TestToJSON(t *testing.T) {
  suites := []struct {
    description string
    original    string
    expected    string
    valid       bool
  }{
    {
      "Valid JSON",
      `{"a":"b"}`,
      `{"a":"b"}`,
      true,
    },
    {
      "Valid YAML",
      `a: b`,
      `{"a":"b"}`,
      true,
    },
    {
      "Invalid YAML",
      `a b`,
      ``,
      false,
    },
    {
      "Invalid JSON",
      `{"a":b}`,
      ``,
      false,
    },
  }
  for _, s := range suites {
    t.Log("  suite: ", s.description)

    resp, err := tojson.ToJSON([]byte(s.original))
    if err != nil {
      if s.valid {
        t.Errorf("Failed to convert data '%s' to JSON due to the error %s", s.original, err.Error())
      } else {
        return  
      }
    }

    if !s.valid {
      t.Errorf("Assumed the suite data '%s' is invalid, but the it was successfully converted to '%s'", s.original, string(resp))
      return
    }

    // if string(resp) != s.assert {
    //   t.Errorf("Test suite object %v doesn't match to the converted data %v", s.assert, string(resp))  
    // }

    assert.Equal(t, s.expected, string(resp), "they should be equal")

  } 
}
