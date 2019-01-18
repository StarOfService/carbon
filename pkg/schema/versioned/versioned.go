package versioned

import (
  "fmt"
  "regexp"
  "github.com/pkg/errors"
  "github.com/sirupsen/logrus"

  "github.com/blang/semver"
)

var re = regexp.MustCompile(`^(?:[a-zA-Z]+/)?v(\d)(?:(alpha|beta)(\d))?$`)

type VersionedConfig interface {
  GetVersion() string
  Parse([]byte) error
  Upgrade() (VersionedConfig, error)
}


type SchemaHandler struct {
  CurrentVersion string
  LatestVersion string
  SchemaVersions map[string]func() VersionedConfig
}

func NewSchemaHandler(current, latest string) *SchemaHandler {
  return &SchemaHandler{
    CurrentVersion: current,
    LatestVersion: latest,
    SchemaVersions: make(map[string]func() VersionedConfig),
  }
}

func (self *SchemaHandler) RegVersion(version string, handler func() VersionedConfig) {
  self.SchemaVersions[version] = handler
}

func (self *SchemaHandler) GetLatestConfig(body []byte) (VersionedConfig, error) {
  factory, ok := self.SchemaVersions[self.CurrentVersion]
  if !ok {
    return nil, errors.Errorf("unknown api version: '%s'", self.CurrentVersion)
  }

  cfg := factory()
  var err error
  if err = cfg.Parse(body); err != nil {
    return nil, errors.Wrap(err, "unable to parse config")
  }

  if cfg.GetVersion() != self.LatestVersion {
    cfg, err = self.upgradeToLatest(cfg)
    if err != nil {
      return nil, err
    }
  }
  return cfg, nil
}

func (self *SchemaHandler) upgradeToLatest(vc VersionedConfig) (VersionedConfig, error) {
  var err error

  // first, check to make sure config version isn't too new
  currentSemver, err := semverParse(vc.GetVersion())
    if err != nil {
      return nil, errors.Wrap(err, "parsing api version") 
  }
  // currentSemver, err := semver.Parse(vc.GetVersion())
  // if err != nil {
  //   currentSemver, err = semverParse(vc.GetVersion())
  //   if err != nil {
  //     return nil, errors.Wrap(err, "parsing api version") 
  //   }
  //   // return nil, errors.Wrap(err, "parsing api version")
  // }


  latestSemver := semverMustParse(self.LatestVersion)
  // latestSemver, err := semver.Parse(self.LatestVersion)
  // if err != nil {
  //   latestSemver = semverMustParse(self.LatestVersion)
  // }

  if currentSemver.EQ(latestSemver) {
    return vc, nil
  }
  if currentSemver.GT(latestSemver) {
    return nil, fmt.Errorf("config version %s is too new for this version of skaffold: upgrade skaffold", vc.GetVersion())
  }

  logrus.Debugf("config version (%s) out of date: upgrading to latest (%s)", vc.GetVersion(), self.LatestVersion)

  for vc.GetVersion() != self.LatestVersion {
    vc, err = vc.Upgrade()
    if err != nil {
      return nil, errors.Wrapf(err, "transforming skaffold config")
    }
  }

  return vc, nil
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

func semverParse(v string) (semver.Version, error) {
  sv := semverNormalize(v)
  return semver.Parse(sv)
}



// func semverParse(v string) (semver.Version, error) {
//   res := re.FindStringSubmatch(v)
//   if res == nil {
//     return semver.Version{}, fmt.Errorf("%s is an invalid api version", v)
//   }
//   if res[2] == "" || res[3] == "" {
//     return semver.Parse(fmt.Sprintf("%s.0.0", res[1]))
//   }
//   return semver.Parse(fmt.Sprintf("%s.0.0-%s.%s", res[1], res[2], res[3]))
// }

func semverMustParse(s string) semver.Version {
  v, err := semverParse(s)
  if err != nil {
    panic(`semver: Parse(` + s + `): ` + err.Error())
  }
  return v
}
