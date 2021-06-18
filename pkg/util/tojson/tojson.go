package tojson

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"

	ghodssyaml "github.com/ghodss/yaml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	jsonPrefix = "{"
	yamlSep    = "---"
)

func ToJSON(data []byte) ([]byte, error) {
	log.Debug("Normalizing config data")
	log.Tracef("Input data: %s", string(data))

	if err := checkJSON(data); err == nil {
		return data, nil
	}

	if err := checkYAML(data); err != nil {
		log.Errorf("Document:\n%s", data)
		return nil, errors.Wrap(err, "verifying YAML document")
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
			log.Errorf("Document:\n%s", d)
			return []byte{}, errors.Wrap(err, "converting YAML document to JSON format")
		}

		jsonDocs = append(jsonDocs, jd)
	}

	resp := bytes.Join(jsonDocs, []byte(""))

	// Checking response format because plain text may be considered as a valid YAML format
	if err := checkJSON(resp); err != nil {
		log.Errorf("Document:\n%s", data)
		return []byte{}, errors.Wrap(err, "verifying JSON document after the YAML->JSON convertion")
	}

	log.Trace("Output data: ", string(resp))
	return resp, nil
}

// json.Valid considers multi-document as a wrong JSON format.
// That's why we need a custom function.
func checkJSON(data []byte) error {
	log.Debug("Trying to verify JSON document")
	trimData := bytes.TrimSpace(data)
	if !bytes.HasPrefix(trimData, []byte(jsonPrefix)) {
		return errors.New("Wrong JSON format")
	}

	dec := json.NewDecoder(bytes.NewReader(trimData))
	for {
		var d interface{}
		err := dec.Decode(&d)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Debug("Failed to unmarshal JSON data due to the error: ", err.Error())
			return err
		}
	}
}

func checkYAML(data []byte) error {
	log.Debug("Trying to verify YAML document")
	dec := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var d interface{}
		err := dec.Decode(&d)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Debug("Failed to unmarshal YAML data due to the error: ", err.Error())
			return err
		}
	}
}
