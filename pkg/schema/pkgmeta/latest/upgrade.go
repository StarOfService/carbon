package latest

import (
  "errors"
  
  "github.com/starofservice/carbon/pkg/schema/versioned"
)

func (c *PackageConfig) Upgrade() (versioned.ConfigHandler, error) {
  return nil, errors.New("not implemented yet")
}