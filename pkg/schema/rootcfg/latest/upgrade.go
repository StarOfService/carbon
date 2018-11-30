package latest

import (
  "errors"
  
  "github.com/starofservice/carbon/pkg/schema/rootcfg/util"
)

func (c *CarbonConfig) Upgrade() (util.VersionedConfig, error) {
  return nil, errors.New("not implemented yet")
}