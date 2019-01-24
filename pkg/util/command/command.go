package command

import (
  "os/exec"
  "os"
  "strings"

  log "github.com/sirupsen/logrus"
)

func Run(cmd string) error {
  log.Debugf("Running system comamnd: %s", cmd)
  parts := strings.Fields(cmd)
  head := parts[0]
  parts = parts[1:len(parts)]
  cmdHandler := exec.Command(head,parts...)
  cmdHandler.Stdout = os.Stdout
  cmdHandler.Stderr = os.Stderr
  return cmdHandler.Run()
}
