package build

import (
  "context"
  "encoding/base64"
  "encoding/json"
  "os"
  "strings"

  clibuild "github.com/docker/cli/cli/command/image/build"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/archive"
  "github.com/docker/docker/pkg/idtools"
  "github.com/docker/docker/pkg/jsonmessage"
  "github.com/docker/docker/pkg/term"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/rootcfg"
)

type Options struct {
  Client *client.Client
  ContextPath string
  RootConfig *rootcfg.CarbonConfig
}

func NewOptions(cfg *rootcfg.CarbonConfig, ctxPath string) (*Options, error) {
  cli, err := client.NewEnvClient()
  if err != nil {
    return nil, errors.Wrap(err, "creating Docker client")
  }

  resp := &Options{
    Client: cli,
    ContextPath: ctxPath,
    RootConfig: cfg,
  }

  return resp, nil
}

// https://github.com/docker/cli/blob/master/cli/command/image/build.go#L40-L76
func (self *Options) Build(metadata map[string]string) error {
  log.Debug("Building docker image")

  excludes, err := clibuild.ReadDockerignore(self.ContextPath)
  if err != nil {
    return errors.Wrap(err, "reading dockerignore file")
  }

  excludes = clibuild.TrimBuildFilesFromExcludes(excludes, self.RootConfig.Data.Dockerfile, false)

  ctx, err := archive.TarWithOptions(self.ContextPath, &archive.TarOptions{
    ExcludePatterns: excludes,
    ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
  })
  if err != nil {
    return errors.Wrap(err, "creating Docker build context")
  }

  // https://github.com/docker/engine/blob/v18.09.0/api/types/client.go#L143-L190
  opt := types.ImageBuildOptions{
    Context:     ctx,
    Dockerfile:  self.RootConfig.Data.Dockerfile,
    ForceRemove: true,
    Labels:      metadata,
    NoCache:     true,
    PullParent:  true,
    Remove:      true,
    Tags:        self.tags(),
  }

  if suppressOutput() {
    opt.SuppressOutput = true  
  }

  response, err := self.Client.ImageBuild(context.Background(), ctx, opt)
  if err != nil {
    return errors.Wrap(err, "building Docker image")
  }

  defer response.Body.Close()

  termFd, isTerm := term.GetFdInfo(os.Stderr)
  jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stderr, termFd, isTerm, nil)

  return nil
}

func suppressOutput() bool {
  logLevel := log.GetLevel().String()
  switch logLevel {
  case "warning", "error", "fatal", "panic":
    return true
  }
  return false
}

func (self *Options) Push() error {
  termFd, isTerm := term.GetFdInfo(os.Stderr)
  for _, i := range self.tags() {
    meta := dockermeta.NewDockerMeta(i)
    username, password, err := meta.GetCredentials()
    if err != nil {
      return errors.Wrap(err, "getting registry credentials")
    }

    auth := types.AuthConfig{
      Username: username,
      Password: password,
    }
    authBytes, _ := json.Marshal(auth)
    authBase64 := base64.URLEncoding.EncodeToString(authBytes)

    opt := types.ImagePushOptions{
      RegistryAuth: authBase64,
    }

    response, err := self.Client.ImagePush(context.Background(), i, opt)
    if err != nil {
      return errors.Wrapf(err, "pushing `%s` docker images", i)
    }
    defer response.Close()

    jsonmessage.DisplayJSONMessagesStream(response, os.Stderr, termFd, isTerm, nil)
  }
  return nil
}

func (self *Options) tags() []string {
  var tags []string
  if len(self.RootConfig.Data.Artifacts) > 0 {
    for _, i := range self.RootConfig.Data.Artifacts {
      parts := strings.Split(i, ":")
      if len(parts) == 1 {
        tags = append(tags, buildTag(parts[0], self.RootConfig.Data.Version))
      } else {
        tags = append(tags, i)
      }
    }
  } else {
    tags = append(tags, buildTag(self.RootConfig.Data.Name, self.RootConfig.Data.Version))
  }
  return tags
}

func buildTag(repo, tag string) string {
  return strings.Join([]string{repo, tag}, ":")
}
