package test

import (
  "testing"

  "github.com/starofservice/carbon/pkg/docker/metadata"
)

func TestDockerRegistryCreds(t *testing.T) {
  h := metadata.NewDockerMeta(DockerMockTestImage)
  u, p, err := h.GetCredentials()
  if err != nil {
    t.Errorf(err.Error())
  }
  if u != DockerMockTestUser {
    t.Errorf("Docker resgistry user doesn't match. Expected: '%s', got: '%s'", DockerMockTestUser, u) 
  }
  if p != DockerMockTestPassword {
    t.Errorf("Docker resgistry password doesn't match. Expected: '%s', got: '%s'", DockerMockTestPassword, p)
  }
}

func TestDockerImageLabels(t *testing.T) {
  suite := []string{
    DockerMockSrcImg,
    DockerMockTestImage,
    "registry:latest",
  }

  for _, i := range suite {
    h := metadata.NewDockerMeta(DockerMockTestImage)
    _, err := h.GetLabels()
    if err != nil {
      t.Errorf("Failed to receive image lables for '%s' due to the error: '%s'", i, err.Error())
    }
  }
}