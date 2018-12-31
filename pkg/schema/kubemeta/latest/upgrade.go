package latest

import (
  "errors"
  
  "github.com/starofservice/carbon/pkg/schema/versioned"
)

func (c *CarbonConfig) Upgrade() (versioned.VersionedConfig, error) {
  return nil, errors.New("not implemented yet")
}