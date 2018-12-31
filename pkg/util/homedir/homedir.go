package homedir

import "os"

func Path() string {
  if h := os.Getenv("HOME"); h != "" {
    return h
  }
  return os.Getenv("USERPROFILE") // windows
}