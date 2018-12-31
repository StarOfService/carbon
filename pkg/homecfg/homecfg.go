package homecfg

import (
  "path/filepath"
  "os"

  "github.com/pkg/errors"

  "github.com/starofservice/carbon/pkg/util/homedir"
)

const (
  HomeConfigDir = ".carbon"
  HomeConfigVarfile = "carbon.vars"
)

func HomeConfigPath() string {
  return filepath.Join(homedir.Path(), HomeConfigDir)
}

func HomeConfigVarfilePath() string {
  return filepath.Join(HomeConfigPath(), HomeConfigVarfile)
}

func InitHomeConfig() error {
  if err := InitHomeConfigDir(); err != nil {
    return err
  }

  if err := InitHomeConfigVarfile(); err != nil {
    return err
  }

  return nil
}

func InitHomeConfigDir() error {
  _, err := os.Stat(HomeConfigPath())
  if err == nil {
    return nil
  }

  if os.IsNotExist(err) {
    mkErr := os.Mkdir(HomeConfigPath(), 0700)
    if mkErr != nil {
      errors.Wrap(mkErr, "creating carbon directory at user home")
    }
    return nil
  }
  return errors.Wrap(err, "checking carbon directory at user home")
}

func InitHomeConfigVarfile() error {
  _, err := os.Stat(HomeConfigVarfilePath())
  if err == nil {
    return nil
  }

  if os.IsNotExist(err) {
    vf, vfErr := os.Create(HomeConfigVarfilePath())
    if vfErr != nil {
      errors.Wrap(vfErr, "creating carbon var file at user home")
    }

    vfErr = vf.Chmod(0600)
    if vfErr != nil {
      errors.Wrap(vfErr, "setting up permissions for carbon var file at user home")
    }
    vf.Close()

    return nil
  }
  return errors.Wrap(err, "checking carbon var file at user home")

}