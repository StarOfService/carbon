package versioned

import (
  "fmt"
  "regexp"

  "github.com/blang/semver"
)

var re = regexp.MustCompile(`^skaffold/v(\d)(?:(alpha|beta)(\d))?$`)

// var schemaVersions = versions{
//   // {v1alpha1.Version, v1alpha1.NewCarbonConfig},
//   {latest.Version, latest.NewCarbonConfig},
// }

type VersionedConfig interface {
  GetVersion() string
  Parse([]byte) error
  Upgrade() (VersionedConfig, error)
}


type SchemaVersion struct {
  apiVersion string
  factory    func() VersionedConfig
}

type SchemaHandler struct {
  CurrentVersion string
  LatestVersion string
  SchemaVersions []SchemaVersion
}

func NewSchemaHandler(current, latest string) *SchemaHandler {
  return &SchemaHandler{
    CurrentVersion: current,
    LatestVersion: latest,
  }
}

func (self *SchemaHandler) RegVersion() (version string, handler func() VersionedConfig) {
  self.SchemaVersions = append(self.SchemaVersions, SchemaVersion{version, handler})
}

func (self *SchemaHandler) GetLatestConfig(body []byte) (VersionedConfig, error) {
  factory, present := self.find(self.CurrentVersion)
  if !present {
    return nil, errors.Errorf("unknown api version: '%s'", apiVersion.Version)
  }

  cfg := factory()
  if err := cfg.Parse(body); err != nil {
    return nil, errors.Wrap(err, "unable to parse config")
  }

  if cfg.GetVersion() != self.LatestVersion {
    cfg, err = upgradeToLatest(cfg)
    if err != nil {
      return nil, err
    }
  }
}


// Find search the constructor for a given api version.
func (self *SchemaHandler) find(apiVersion string) (func() VersionedConfig, bool) {
  for _, version := range *self.SchemaVersions {
    if version.apiVersion == apiVersion {
      return version.factory, true
    }
  }

  return nil, false
}


func (self *SchemaHandler) upgradeToLatest(vc VersionedConfig) (VersionedConfig, error) {
  var err error

  // first, check to make sure config version isn't too new
  version, err := semverParse(vc.GetVersion())
  if err != nil {
    return nil, errors.Wrap(err, "parsing api version")
  }

  semver := semverMustParse(self.LatestVersion)
  if version.EQ(semver) {
    return vc, nil
  }
  if version.GT(semver) {
    return nil, fmt.Errorf("config version %s is too new for this version of skaffold: upgrade skaffold", vc.GetVersion())
  }

  logrus.Warnf("config version (%s) out of date: upgrading to latest (%s)", vc.GetVersion(), self.LatestVersion)

  for vc.GetVersion() != self.LatestVersion {
    vc, err = vc.Upgrade()
    if err != nil {
      return nil, errors.Wrapf(err, "transforming skaffold config")
    }
  }

  return vc, nil
}

func semverParse(v string) (semver.Version, error) {
  res := re.FindStringSubmatch(v)
  if res == nil {
    return semver.Version{}, fmt.Errorf("%s is an invalid api version", v)
  }
  if res[2] == "" || res[3] == "" {
    return semver.Parse(fmt.Sprintf("%s.0.0", res[1]))
  }
  return semver.Parse(fmt.Sprintf("%s.0.0-%s.%s", res[1], res[2], res[3]))
}

func semverMustParse(s string) semver.Version {
  v, err := semverParse(s)
  if err != nil {
    panic(`semver: Parse(` + s + `): ` + err.Error())
  }
  return v
}
