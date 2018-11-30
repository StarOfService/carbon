package util

type VersionedConfig interface {
  GetVersion() string
  Parse([]byte) error
  Upgrade() (VersionedConfig, error)
}
