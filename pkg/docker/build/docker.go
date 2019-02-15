package build

import (
  "bytes"
  "context"
  "encoding/base64"
  "encoding/json"
  "io"
  "os"
  // "strconv"
  "strings"

  // "github.com/containerd/containerd"

  // "github.com/docker/cli/cli/command"
  // "github.com/docker/cli/cli/command/image"
  // clitypes "github.com/docker/cli/types"

  clibuild "github.com/docker/cli/cli/command/image/build"
  "github.com/docker/docker/pkg/archive"
  "github.com/docker/docker/pkg/idtools"

  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/jsonmessage"
  "github.com/docker/docker/pkg/term"
  "github.com/pkg/errors"
  log "github.com/sirupsen/logrus"

  dockermeta "github.com/starofservice/carbon/pkg/docker/metadata"
  "github.com/starofservice/carbon/pkg/schema/pkgcfg"
)

// const containerdSockPath = "/run/containerd/containerd.sock"

type Options struct {
  Client *client.Client
  ContextPath string
  RootConfig *pkgcfg.CarbonConfig
  Tags []string
}

func NewOptions(cfg *pkgcfg.CarbonConfig, ctxPath string) (*Options, error) {
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

func (self *Options) ExtendTags(cliTags []string, prefix string, suffix string) error {
  var selectTags []string
  if len(cliTags) > 0 {
    selectTags = cliTags
  } else if len(self.RootConfig.Data.Artifacts) > 0 {
    selectTags = self.RootConfig.Data.Artifacts
  } else {
    selectTags = append(selectTags, joinTag(self.RootConfig.Data.Name, self.RootConfig.Data.Version))
  }

  for _, i := range selectTags {
    im, err := dockermeta.NewDockerMeta(i)
    if err != nil {
      return err
    }
    name := im.Name()

    var tag string
    if i == name || "docker.io/" + i == name {
      tag = self.RootConfig.Data.Version
    } else {
      tag = im.Tag()
    }

    fullTag := tag
    if fullTag != "latest" {
      fullTag = joinTag(name, (prefix + tag + suffix))
    }
    self.Tags = append(self.Tags, fullTag)
  }
  return nil
}

// https://github.com/docker/cli/blob/master/cli/command/image/build.go#L40-L76
func (self *Options) Build(metadata map[string]string) error {
  log.Debug("Building Docker image")

  excludes, err := clibuild.ReadDockerignore(self.ContextPath)
  if err != nil {
    return errors.Wrap(err, "reading dockerignore file")
  }

  excludes = clibuild.TrimBuildFilesFromExcludes(excludes, self.RootConfig.Data.Dockerfile, false)

  ctx, err := archive.TarWithOptions(self.ContextPath, &archive.TarOptions{
    ExcludePatterns: excludes,
    ChownOpts:       &idtools.IDPair{UID: 0, GID: 0},
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
    Tags:        self.Tags,
  }

  response, err := self.Client.ImageBuild(context.Background(), ctx, opt)
  if err != nil {
    return errors.Wrap(err, "building Docker image")
  }
  defer response.Body.Close()

  displayJSONMsg(response.Body)

  return nil
}

// func (self *Options) Build(metadata map[string]string) error {

//   stdin, stdout, stderr := term.StdStreams()
//   dockerCli := command.NewDockerCli(stdin, stdout, stderr, contentTrustEnabled()) //, newContainerizedClient)

//   cmd := image.NewBuildCommand(dockerCli)
//   cmd.Flags().Set("file", self.RootConfig.Data.Dockerfile)
//   cmd.Flags().Set("force-rm", "true")
//   cmd.Flags().Set("no-cache", "true")
//   // cmd.Flags().Set("pull", "true") // TODO: do we really need this?
//   cmd.Flags().Set("rm", "true")
//   cmd.SetArgs([]string{"."})

//   for k, v := range metadata {
//     cmd.Flags().Set("label", k + "=" + v) 
//   }
   
//   for _, i := range self.Tags {
//     cmd.Flags().Set("tag", i) 
//   }

//   err := cmd.Execute()

//   if err != nil {
//     return errors.Wrap(err, "building docker image")
//   }


  
//   // cmd.SetOutput(ioutil.Discard)
//   // err := cmd.Execute()
//   return nil

// }

// func newContainerizedClient(sockPath string) (clitypes.ContainerizedClient, error) {
//   if sockPath == "" {
//     sockPath = containerdSockPath
//   }
//   cclient, err := containerd.New(sockPath)
//   if err != nil {
//     return nil, err
//   }
//   return &baseClient{
//     cclient: cclient,
//   }, nil
// }

func (self *Options) Push() error {
  for _, i := range self.Tags {
    meta, err := dockermeta.NewDockerMeta(i)
    if err != nil {
      return err
    }
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
      return errors.Wrapf(err, "pushing Docker image `%s`", i)
    }
    defer response.Close()

    displayJSONMsg(response)
  }
  return nil
}

func (self *Options) Remove() error {
  for _, i := range self.Tags {

    opt := types.ImageRemoveOptions{
      // Force: true,
      PruneChildren: true,
    }

    response, err := self.Client.ImageRemove(context.Background(), i, opt)
    if err != nil {
      return errors.Wrapf(err, "removing Docker image `%s`", i)
    }
    for _, i := range response {
      if i.Untagged != "" {
        log.Debug("Untagged: ", i.Untagged)
      }
      if i.Deleted != "" {
        log.Debug("Deleted: ", i.Deleted)
      }
    }
  }
  return nil
}

func joinTag(repo, tag string) string {
  return strings.Join([]string{repo, tag}, ":")
}

func suppressOutput() bool {
  logLevel := log.GetLevel().String()
  switch logLevel {
  case "warning", "error", "fatal", "panic":
    return true
  }
  return false
}

func displayJSONMsg(in io.Reader) {
  var out io.Writer
  if suppressOutput() {
    out = &bytes.Buffer{}
  } else {
    out = os.Stdout
  }

  termFd, isTerm := term.GetFdInfo(out)
  jsonmessage.DisplayJSONMessagesStream(in, out, termFd, isTerm, nil)
}

// // https://github.com/docker/cli/blob/v18.09.2/cmd/docker/docker.go#L197-L205
// func contentTrustEnabled() bool {
//   if e := os.Getenv("DOCKER_CONTENT_TRUST"); e != "" {
//     if t, err := strconv.ParseBool(e); t || err != nil {
//       // treat any other value as true
//       return true
//     }
//   }
//   return false
// }