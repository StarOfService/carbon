package util

import "github.com/starofservice/flapper"

const (
  MetaPrefix = "c6"
  MetaDelimiter = "."
)

type VersionedConfig interface {
  GetVersion() string
  Parse(map[string]string) error
  Upgrade() (VersionedConfig, error)
}

func NewFlapper() (*flapper.Flapper, error) {
  return flapper.New(MetaPrefix, MetaDelimiter)
}