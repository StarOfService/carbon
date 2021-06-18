package latest

import (
	"encoding/json"
	"errors"

	"github.com/starofservice/vconf"
)

const Version string = "v1alpha1"

type KubeMetadata struct {
	APIVersion string            `json:"apiVersion"`
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Source     string            `json:"source"`
	Variables  map[string]string `json:"variables"`
	Patches    string            `json:"patches"`
	Namespace  string            `json:"namespace"`
	Manifest   string            `json:"manifest"`
	// dependencies
}

func NewKubeMetadata() vconf.ConfigInterface {
	return new(KubeMetadata)
}

func (c *KubeMetadata) GetVersion() string {
	return c.APIVersion
}

func (c *KubeMetadata) Parse(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}

	return nil
}

func (c *KubeMetadata) Upgrade() (vconf.ConfigInterface, error) {
	return nil, errors.New("not implemented yet")
}
