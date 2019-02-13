package variables

import (
  "regexp"

  "github.com/magiconair/properties"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
)

type Vars struct {
  Data map[string]string
}

func NewVars() *Vars {
  return &Vars{Data: make(map[string]string)}
}

func (self *Vars) ParseFiles(vf []string) error {
  log.Debug("Parsing variable files")

  l := &properties.Loader{
    Encoding: properties.UTF8,
    IgnoreMissing: false,
  }

  for _, i := range vf {
    log.Tracef("Parsing %s", i)
    p, err := l.LoadFile(i)
    if err != nil {
      return errors.Wrapf(err, "parsing variable file '%s'", i)
    }

    for k, v := range p.Map() {
      self.Data[k] = v
    }
  }

  return nil
}

func (self *Vars) ParseFlags(flags []string) error {
  log.Debug("Parsing variable flags")

  stringMapRegex := regexp.MustCompile("[=]")
  for _, v := range flags {
    log.Tracef("Parsing %s", v)
    parts := stringMapRegex.Split(v, 2)
    if len(parts) != 2 {
      return errors.Errorf("Expected KEY=VALUE format, but got '%s'", v)
    }
    self.Data[parts[0]] = parts[1]
  }

  return nil
}
