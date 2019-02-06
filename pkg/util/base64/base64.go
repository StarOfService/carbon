package base64

import (
  "encoding/base64"

  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
)

func Encode(data []byte) string {
  log.Debug("Encoding data to Base64 format")
  return base64.StdEncoding.EncodeToString(data)
}

func Decode(data string) ([]byte, error) {
  log.Debug("Decoding data from Base64 format")
  resp, err := base64.StdEncoding.DecodeString(data)
  if err != nil {
    return nil, errors.Wrapf(err, "decoding string `%s`", data)
  }
  return resp, nil
}
