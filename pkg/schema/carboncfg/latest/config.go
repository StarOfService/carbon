package latest

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/starofservice/vconf"
)

const Version string = "v1alpha1"

type CarbonConfig struct {
	APIVersion  string `json:"apiVersion"`
	CarbonScope string `json:"carbonScope"`
}

func NewCarbonConfig() vconf.ConfigInterface {
	return new(CarbonConfig)
}

func (c *CarbonConfig) GetVersion() string {
	return c.APIVersion
}

func (c *CarbonConfig) Parse(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}

	return nil
}

func (c *CarbonConfig) Upgrade() (vconf.ConfigInterface, error) {
	return nil, errors.New("not implemented yet")
}
