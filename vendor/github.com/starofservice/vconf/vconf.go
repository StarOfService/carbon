package vconf

import (
  "fmt"
  "regexp"

  "github.com/blang/semver"
  "github.com/pkg/errors"
)

var re = regexp.MustCompile(`^(?:[a-zA-Z]+/)?v(\d)(?:(alpha|beta)(\d))?$`)

// ConfigInterface defines interface for config version
// Every single config version must implement this interface
type ConfigInterface interface {
  GetVersion() string
  Parse([]byte) error
  Upgrade() (ConfigInterface, error)
}

// SchemaHandler is a management interface for the whole versioned config
type SchemaHandler struct {
  LatestVersion string
  SchemaVersions map[string]func() ConfigInterface
}

// NewSchemaHandler creates a new instance of a SchemaHandler
func NewSchemaHandler(latest string) *SchemaHandler {
  return &SchemaHandler{
    LatestVersion: latest,
    SchemaVersions: make(map[string]func() ConfigInterface),
  }
}

// RegVersion adds a version of your config to the SchemaVersions map.
// All available versions must be registered using this method.
// Order doesn't matter.
func (sh *SchemaHandler) RegVersion(version string, handler func() ConfigInterface) {
  sh.SchemaVersions[version] = handler
}

// GetLatestConfig is trying to find config version at the SchemaVersions map
// and to upgrade the given config up to the latest version.
func (sh *SchemaHandler) GetLatestConfig(configVersion string, body []byte) (ConfigInterface, error) {
  factory, ok := sh.SchemaVersions[configVersion]
  if !ok {
    return nil, errors.Errorf("unknown version: '%s'", configVersion)
  }

  cfg := factory()
  var err error
  if err = cfg.Parse(body); err != nil {
    return nil, errors.Wrap(err, "parsing source config")
  }

  if cfg.GetVersion() != sh.LatestVersion {
    cfg, err = sh.upgradeToLatest(cfg)
    if err != nil {
      return nil, err
    }
  }
  return cfg, nil
}

func (sh *SchemaHandler) upgradeToLatest(vc ConfigInterface) (ConfigInterface, error) {
  currentSemver, err := semverParse(vc.GetVersion())
    if err != nil {
      return nil, errors.Wrap(err, "converting current version to semver format") 
  }
  latestSemver := semverMustParse(sh.LatestVersion)

  if currentSemver.EQ(latestSemver) {
    return vc, nil
  }
  if currentSemver.GT(latestSemver) {
    return nil, errors.Errorf(
      "the current version '%s' is higher than the latest supported version '%s': utility upgrade is required",
      vc.GetVersion(),
      sh.LatestVersion,
    )
  }

  for vc.GetVersion() != sh.LatestVersion {
    iv := vc.GetVersion()
    vc, err = vc.Upgrade()
    if err != nil {
      return nil, errors.Wrapf(err, "running upgrade for version '%s'", iv)
    }
  }

  return vc, nil
}

func semverParse(v string) (semver.Version, error) {
  sv := semverNormalize(v)
  return semver.Parse(sv)
}

func semverMustParse(s string) semver.Version {
  v, err := semverParse(s)
  if err != nil {
    panic(`semver: Parse(` + s + `): ` + err.Error())
  }
  return v
}

func semverNormalize(v string) string {
  res := re.FindStringSubmatch(v)
  if res == nil {
    return v
  }
  if res[2] == "" || res[3] == "" {
    return fmt.Sprintf("%s.0.0", res[1])
  }
  return fmt.Sprintf("%s.0.0-%s.%s", res[1], res[2], res[3])
}
