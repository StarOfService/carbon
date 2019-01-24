package main

import (
  "github.com/starofservice/carbon/cmd"
  "github.com/starofservice/carbon/pkg/version"
)

// VERSION is set during build
var VERSION = "0.0.0"

func main() {
  version.VERSION = VERSION
  cmd.Execute()
}
