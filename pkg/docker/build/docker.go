package build

import (
  "context"
  // "fmt"
  "strings"
  "os"

  log "github.com/sirupsen/logrus"

  clibuild "github.com/docker/cli/cli/command/image/build"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/archive"
  "github.com/docker/docker/pkg/idtools"
  "github.com/docker/docker/pkg/jsonmessage"
  "github.com/docker/docker/pkg/term"

  rootcfglatest "github.com/starofservice/carbon/pkg/schema/rootcfg/latest"

)

// https://github.com/docker/cli/blob/master/cli/command/image/build.go#L40-L76
// https://github.com/moby/moby/blob/v1.13.1/api/types/client.go#L142-L178
type BuildOptions struct {
  Tags      []string
  Labels    map[string]string
  BuildArgs map[string]*string
  // Quiet     bool
  NoCache   bool
  Squash    bool
}

func NewBuildOptions() *BuildOptions {
  return new(BuildOptions)
}

func (o *BuildOptions) Build(cfg *rootcfglatest.CarbonConfig, ctxPath string, metadata map[string]string) {
  log.Debug("Building docker image")

  excludes, err := clibuild.ReadDockerignore(ctxPath)
  if err != nil {
    log.Fatalf("Failed to read Dockerignore file due to the error: %s", err.Error())
    os.Exit(1)
  }

  // if err := build.ValidateContextDirectory(contextDir, excludes); err != nil {
  //   return errors.Errorf("error checking context: '%s'.", err)
  // }

  // relDockerfile string

  // // And canonicalize dockerfile name to a platform-independent one
  // relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

  excludes = clibuild.TrimBuildFilesFromExcludes(excludes, cfg.Dockerfile, false)

  ctx, err := archive.TarWithOptions(ctxPath, &archive.TarOptions{
    ExcludePatterns: excludes,
    ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
  })
  if err != nil {
    log.Fatalf("Failed to create Docker build context due to the error: %s", err.Error())
    os.Exit(1)
  }

  // https://github.com/docker/engine/blob/v18.09.0/api/types/client.go#L143-L190
  opt := types.ImageBuildOptions{
    // BuildArgs:   args,
    Context:     ctx,
    Dockerfile:  cfg.Dockerfile,
    ForceRemove: true,
    Labels:      metadata,
    NoCache:     true,
    PullParent:  true,
    Remove:      true,
    Tags:        buildOptionsTags(cfg),
  }

  if suppressOutput() {
    opt.SuppressOutput = true  
  }
    
  cli, err := client.NewEnvClient()
  if err != nil {
    log.Fatalf("Failed to create Docker client due to the error: %s", err.Error())
    os.Exit(1)
  }

  response, err := cli.ImageBuild(context.Background(), ctx, opt)
  if err != nil {
    // https://github.com/docker/cli/blob/master/cli/command/image/build.go#L405-L411
    log.Fatalf("Failed to build Docker image due to the error: %s", err.Error())
    os.Exit(1)
  }

  defer response.Body.Close()

  termFd, isTerm := term.GetFdInfo(os.Stderr)
  jsonmessage.DisplayJSONMessagesStream(response.Body, os.Stderr, termFd, isTerm, nil)
}

func suppressOutput() bool {
  logLevel := log.GetLevel().String()
  // if logLevel == "warning" || logLevel == "error" || logLevel == "fatal" || logLevel == "panic" {
  switch logLevel {
  case "warning", "error", "fatal", "panic":
    return true
  }
  return false
}

func buildOptionsTags(cfg *rootcfglatest.CarbonConfig) []string {
  var tags []string
  tags = append(tags, buildTag(cfg.Name, cfg.Version))
  tags = append(tags, buildTag(cfg.Name, "latest"))
  tags = append(tags, cfg.Artifacts...)

  return tags
}

func buildTag(repo, tag string) string {
  return strings.Join([]string{repo, tag}, ":")
}
