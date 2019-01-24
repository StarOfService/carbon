package version

import (
  "fmt"

  "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"
)

var VERSION string

func PrintVersion() {
  fmt.Println("Carbon utility version:", VERSION)
  fmt.Println("Latest apiVersion for carbon.yaml:", latest.Version)
}