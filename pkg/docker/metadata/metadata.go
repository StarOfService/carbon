package metadata  

import (
  "context"
  "strings"

  "github.com/containers/image/transports/alltransports"
  "github.com/containers/image/types"
  "github.com/docker/cli/cli/config"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/term"
  "github.com/docker/distribution/reference"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"
)

const (
  kubeImageOS = "linux"
)

type DockerMeta struct {
  image string
  ref types.ImageReference
}

func NewDockerMeta(image string) *DockerMeta {

  imagePrefixed := image
  if !strings.HasPrefix("docker://", image) {
    imagePrefixed = "docker://" + image
  }

  ref, err := alltransports.ParseImageName(imagePrefixed)
  if err != nil {
    panic(err)
  }

  self := &DockerMeta{
    image: image,
    ref: ref,
  }

  return self
}

func (self *DockerMeta) GetLabels() (map[string]string, error) {
  log.Debug("Getting dokcer image labels")

  log.Debug("Trying to receive the lables from a locally avaiable image")
  resp, err := self.getLocalImageLabels()
  if err == nil {
    return resp, nil
  }
  log.Debug("Got an error: %s", err.Error())

  sys := &types.SystemContext{
    OSChoice: kubeImageOS,
  }

  log.Debug("Trying to receive the lables without authentication for a public repo")
  resp, err = self.getRemoteMetaLabels(sys)
  if err == nil {
    return resp, nil
  }

  username, password, err := self.GetCredentials()
  if err != nil {
    return nil, err
  }

  sys.DockerAuthConfig = &types.DockerAuthConfig{
    Username: username,
    Password: password,
  }

  resp, err = self.getRemoteMetaLabels(sys)
  if err != nil {
    return nil, errors.Wrapf(err, "getting Carbon metadata for a repository '%s'", self.Name())
  }

  return resp, nil
}

func (self *DockerMeta) Domain() string {
  return reference.Domain(self.dockerReference())
}

func (self *DockerMeta) Name() string {
  return self.dockerReference().Name()
}

func (self *DockerMeta) Tag() string {
  return self.dockerReference().Tag()
}

func (self *DockerMeta) dockerReference() reference.NamedTagged {
  return self.ref.DockerReference().(reference.NamedTagged)
}

func (self *DockerMeta) getLocalImageLabels() (map[string]string, error) {
  cli, err := client.NewEnvClient()
  if err != nil {
    panic(err)
  }
  ctx := context.Background()

  element, _, err := cli.ImageInspectWithRaw(ctx, self.image)
  if err != nil {
    return nil, err
  }
  return element.Config.Labels, nil
}

func (self *DockerMeta) getRemoteMetaLabels(sys *types.SystemContext) (map[string]string, error) {
  ctx := context.Background()

  img, err := self.ref.NewImage(ctx, sys)
  if err != nil {
    log.Debugf("Failed to create an image handler due to the error: %s", err.Error())
    return nil, err
  }

  imgInspect, err := img.Inspect(ctx)
  if err != nil {
    log.Debugf("Failed to receive a docker image metadata due to the error: %s", err.Error())
    return nil, err
  }

  return imgInspect.Labels, nil
}

func (self *DockerMeta) GetCredentials() (string, string, error) {
  registry := self.Domain()

  _, _, stderr := term.StdStreams()
  dockerConfig := config.LoadDefaultConfigFile(stderr)
  creds, err := dockerConfig.GetAuthConfig(registry)
  if err != nil {
      return "", "", errors.Wrapf(err, "extracting Docker credentials for a repository '%s'", self.Name())
  }

  if len(creds.Username) == 0 || len(creds.Password) == 0 {
    return "", "", errors.Errorf("Got an empty docker username or password for a repository '%s'", self.Name())
  }

  return creds.Username, creds.Password, nil
}
