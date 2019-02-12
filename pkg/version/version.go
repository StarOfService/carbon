package version

import (
  "fmt"

  "github.com/starofservice/carbon/pkg/schema/pkgcfg/latest"
)

const DefaultVersion = "0.0.0"
var VERSION string

func PrintVersion() {
  v := DefaultVersion
  if VERSION != "" {
    v = VERSION
  }
  fmt.Println("Carbon utility version:", v)
  fmt.Println("Latest apiVersion for carbon.yaml:", latest.Version)
}