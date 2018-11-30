package latest

import (
  "errors"
  
  "github.com/starofservice/carbon/pkg/schema/pkgmeta/util"
)

func (c *PackageConfig) Upgrade() (util.VersionedConfig, error) {
  return nil, errors.New("not implemented yet")
}