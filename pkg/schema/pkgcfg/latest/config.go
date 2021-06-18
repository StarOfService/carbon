package latest

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/starofservice/vconf"

	"github.com/starofservice/carbon/pkg/schema/pkgcfg/defval"
)

const Version string = "v1alpha1"

type CarbonConfigVariable struct {
	Name        string `json:"name"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type CarbonConfigHooks struct {
	PreBuild  []string `json:"pre-build"`
	PostBuild []string `json:"post-build"`
}

type CarbonConfig struct {
	APIVersion    string                 `json:"apiVersion"`
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	Dockerfile    string                 `json:"dockerfile"`
	KubeManifests string                 `json:"kubeManifests"`
	Artifacts     []string               `json:"artifacts"`
	Variables     []CarbonConfigVariable `json:"variables"`
	Hooks         CarbonConfigHooks      `json:"hooks"`
	// dependencies
}

func NewCarbonConfig() vconf.ConfigInterface {
	return new(CarbonConfig)
}

func (c *CarbonConfig) GetVersion() string {
	return c.APIVersion
}

func (c *CarbonConfig) Parse(contents []byte) error {
	if err := json.Unmarshal(contents, c); err != nil {
		return err
	}

	if c.Name == "" {
		return errors.New("'name' parameter isn't defined")
	}
	if c.Version == "" {
		return errors.New("'version' parameter isn't defined")
	}
	if c.KubeManifests == "" {
		return errors.New("'kubeManifests' parameter isn't defined")
	}

	if c.Dockerfile == "" {
		c.Dockerfile = defval.Dockerfile
	}

	return nil
}

func (c *CarbonConfig) Upgrade() (vconf.ConfigInterface, error) {
	return nil, errors.New("not implemented yet")
}
