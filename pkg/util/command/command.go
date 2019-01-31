package command

import (
  "io"
  "os/exec"
  // "os"
  "strings"

  log "github.com/sirupsen/logrus"
)

func Run(cmd, cwd string, stdout, stderr io.Writer) error {
  log.Debugf("Running system comamnd: %s", cmd)
  parts := strings.Fields(cmd)
  head := parts[0]
  parts = parts[1:len(parts)]
  cmdHandler := exec.Command(head,parts...)
  if cwd != "" {
    cmdHandler.Dir = cwd
  }
  if stdout != nil {
    cmdHandler.Stdout = stdout
  }
  if stderr != nil {
    cmdHandler.Stderr = stderr
  }
  // cmdHandler.Stdout = os.Stdout
  // cmdHandler.Stderr = os.Stderr
  return cmdHandler.Run()
}
