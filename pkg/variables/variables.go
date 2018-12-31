package variables

import (
  // "fmt"
  "os"
  "regexp"


  "github.com/magiconair/properties"
  log "github.com/sirupsen/logrus"

  // pkgmetalatest "github.com/starofservice/carbon/pkg/schema/pkgmeta/latest"
)


type Vars struct {
  Data map[string]string
}

func NewVars() *Vars {
  // self := &Vars{Data: make(map[string]string)}

  // for _, v := range m.Variables {
  //   self.Data[v.Name] = v.Default
  // }

  return &Vars{Data: make(map[string]string)}
}

// func PropErrorHandler(err error) {
//   log.Fatal("Failed to parse variables file due to the error: %s", err.Error())
//   os.Exit(1)
// }

func (self *Vars) ParseVarFiles(vf []string) {
  log.Debug("Parsing variable files")
  // properties.ErrorHandler = PropErrorHandler
  // p := properties.LoadFiles(vf, properties.UTF8, false)
  // for k, v := range p.Map() {
  //   self.Data[k] = v
  // }
  l := &properties.Loader{
    Encoding: properties.UTF8,
    IgnoreMissing: false,
  }

  for _, i := range vf {
    log.Tracef("Parsing %s", i)
    p, err := l.LoadFile(i)
    if err != nil {
      log.Fatalf("Failed to parse a variables file '%s', due to the error: '%s'. Skipping it.", i, err.Error())
      os.Exit(1)
      // continue
    }

    for k, v := range p.Map() {
      self.Data[k] = v
    }
  }
}

func (self *Vars) ParseFlags(flags []string) {
  log.Debug("Parsing variable flags")

  stringMapRegex := regexp.MustCompile("[=]")
  for _, v := range flags {
    log.Tracef("Parsing %s", v)
    // kd.Variables.Var[v.Name] = v.Default
    parts := stringMapRegex.Split(v, 2)
    if len(parts) != 2 {
      log.Fatalf("Expected KEY=VALUE format, but got '%s'. Skipping this item.", v)
      os.Exit(1)
      // continue
    }

    self.Data[parts[0]] = parts[1]
  }
}
