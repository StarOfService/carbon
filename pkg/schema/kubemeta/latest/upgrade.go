package latest

import (
  "errors"
  
  "github.com/starofservice/carbon/pkg/schema/versioned"
)

func (c *KubeMetadata) Upgrade() (versioned.ConfigHandler, error) {
  return nil, errors.New("not implemented yet")
}