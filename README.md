Carbon is a package manager for Kubernetes. It's designed to leverage standard docker registry API and doesn't impose any additional requirements or restrictions. Thus you may use your favorite docker registry in order to store, share and install Carbon packages to your Kubernetes cluster.

## How Carbon works
Carbon package is a docker image with additional metadata stored and the image lables. The most important part of the metadata is a Kubernetes manifest templates, which are converted to Kubernetes manifests on the installation stage.

Usage of docker lables allows to avoid downloading of the whole image for the installation. Carbon just reads image labels directly from a registry, builds a Kubernetes manifest based on the templates and applies this manifest to the Kubernetes cluster. Based on the resources definition, Kubernetes downloads a Docker image from the regestry and launches Pods.

## Available CLI operations
- `carbon build` - build Carbon package. It builds a Docker image with a Carbon metadata. Since it's built. you can operate with this Carbon package like you do with a common Docker image: assign additional tags, run `docker inspect`, push to a docker registry and so on.
- `carbon inspect dockerImage [dockerImage ...]` - this command exposes Carbon package metadata from the given docker images. It supports both local and remote docker images. If it's a private repository, you have to be authenticated and authorized to working with this repo (see `docker login`). The command shows Carbon *name*, *version* and *variables* for given images.
- `carbon install dockerImage` - reads Carbon metadata for a givem image, builds kubernetes manifest using default and provided variables, if any patches are provided, they are applied for the manifests, and this manifest with some Carbon metadata is deployed to Kubernetes.
- `carbon status [packageName [packageName ...]]` - if no packages are provided, Carbon lists all installed packages for a current Kubernetes context. When at least one package name is provided, Carbon shows detailed information for this package. By default it's package *name*, *version*, *source* used for an installation, *variable values* used for buiding the package manifest. When `--full` flag is provided, you also get a list of applied patches and a final Kubernetes manifest which applied to the Kubernetes cluster.
- `carbon delete packageName [packageName ...]` - delete specified packages
- `carbon version` - show version for your capbon CLI tool and the latest supported version of the main config.

Almost all mentioned subcommands have additional flags which extend the described behavior. So don't hestitate to use CLI built-in help, for example: `carbon build -h` or `carbon help build`

## Main config (TODO: check which fields are really required)
Main config is a YAML file located at the root of a package codebase. Every Carbon package must containe the main Carbon config.
This config supports such the fields:

Parameter | Required | Default value | Description
---|---|---|---
apiVersion | true | | version of the config format. The lastest version is `v1alpha1`
name | true | | name of the Carbon package
version | true | | version of the Carbon package
Dockerfile  | false | ./Dockerfile | path to a Dockerfile which will be used for building a docker iamge with your application
kubeManifests | true | | path to Kubernetes manifest templates. Current parameter may contain wildcard, for example `./k8s/*.yaml`.
artifacts | false | _name_:_version_ | list of docker tags to be assigned for the built image. This list may contain a full tag like `carbontest:latest` or just a name without a tag suffix. When a tag siffux isn't provided, it will the package version will be used as a siffux. Thus if `version` is set to `0.1.0` and artifacts has item `carbontest`, the image will get a tag `carbontest:0.1.0`
variables | false | | List of all variables used at a kubernetes manifest teamplates, must be defined here. More details you can find at ..... section (TODO)
hooks | false | | Here you can define lifecycle hooks, a system calls which have to be ran on a certan stages of a package life. More details can be fond at ...... section (TODO)

Here's an example of the main config:
```
apiVersion: v1alpha1
dockerfile: Dockerfile
kubeManifests: k8s/*.yaml
name: carbon-test
version: 0.0.1
artifacts:
  - foo/carbontest
  - foo/carbontest:latest
hooks:
  pre-build:
    - minikube version
  post-build:
    - hooks/build.sh
variables:
  - name: FullName
    default: carbon-test
  - name: Environment
    default: local
  - name: KubeNamespace
    default: default
  - name: Cluster
    default: core
    description: K8s cluster name
```

Currently the latest `apiVersion` is `v1alpha1`, but you shouldn't worry about migration to the new version when it's released. Carbon versioning is based on [VConf](https://github.com/StarOfService/vconf) library which allows to translate any old version up to the latest seamlessly for a user. The same approach is used for Carbon package metadata and installed packages metadata. In the future we're going to add `carbon config-upgrade` command which will let you to upgrade an old version of your main config up to the latest one by one command call.

## Kubernetes manifest templates and variables
Kubernetes manifest templates are standard Go templates, so you can find the whole necessarey information about the format at the documentation for `text/template` package: https://golang.org/pkg/text/template/

At the templates you can use two kinds of variables:
- Package properties which are filled by Carbon automatically:
- User-defined variables

### Package properties
Package properties are filled by Carbon automatically. Currently these properties are exposed:
- `.Pkg.Name` - name of the Carbon package
- `.Pkg.Version` - version of the Carbon package
- `.Pkg.DockerName` - name of the docker image which is used for the installation (for example _ubuntu_)
- `.Pkg.DockerTag` - tag suffix of the docker image which is used for the installation (for example _latest_)

It's up to you how to use these properties, but here we want to give some recommendations:
- `.Pkg.Name` may be useful for naming Kubernetes resources and thus to avoid conflicts with other packages.
- combination of `.Pkg.DockerName` and `.Pkg.DockerTag` properties will let you to avoid hardcoding of a Pod image

### User-defined variables
All user-defined variables must be defined at `variables` section of the main config. Each variable may have such parameters:

Parameter | Required | Default | Description
---|---|---|---
`name` | true | - | variable name
`default` | false | "" | Default value for the variable
`description` | false | "" | Description for the variable

TODO: write requirements for variable names!

When you use these variables at a kubernetes manifest template, they must be prefixed by `.Var.`, for example: `.Var.Environment`.

At the installation stage, variables can be takein from different sources and are applied inthe following order:
1. Default values defined at the main config
2. `~/carbon.vars` - variables file from a user home directory
3. Variables files provied by `--var-file` flag
4. Variables provied by `--var` flag

Variables file, in fact is a .property file (https://en.wikipedia.org/wiki/.properties). In most cases it's enough to use a simple key=value notation:
```
Cluster = core
Environment=Local
```
But if you have an advanced use-case, please follow the `.properties` format and remember that variables are disallowed to have spaces at name.

## Docker image names and tags
You have several optins to define docker image name and tag:
1. Don't do anything, in this case docker image will be named by Carbon package name and version, for example: "carbontest:0.0.1"
2. Define names and tags at `artifacts` section of the main config
3. Define names and tags by `--tag` flag for `carbon build`. `--tag` flag can be used multiple times in order to assigne multiple names/tags. This option has precedence over the `artifacts` section from the main config.

When you define a name without a tag, a Carbon package version will be used as a tag.
`carbon build` command has additional flags `--tag-prefix` and `--tag-suffix`. When any of this flag is used, all docker iamge tags will be prfixed/suffixed by the provided strings. Here's one exception: prefix and siffux are not applied for the `latest` tag.

Let's consider an example. We have a Carbon package with version `0.1.0` and run the command: 
```
$ carbon build \
  --tag region1.registry.example.com/carbontest \
  --tag region1.registry.example.com/carbontest:foo \
  --tag region1.registry.example.com/carbontest:latest \
  --tag region2.registry.example.com/carbontest \
  --tag region2.registry.example.com/carbontest:foo \
  --tag region2.registry.example.com/carbontest:latest \
  --tag-prefix "hotfix-" \
  --taf-suffix "-alpha"
```

This command will create a docker image with such tags:
- region1.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region1.registry.example.com/carbontest:hotfix-foo-alpha
- region1.registry.example.com/carbontest:latest
- region2.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region2.registry.example.com/carbontest:hotfix-foo-alpha
- region2.registry.example.com/carbontest:latest

## Hooks
Hooks is a way to run system calls on a specific stages of a Carbon package lifecycle. It may be a direct system call, or a shell script.

Currently Carbon allows to run hooks on two stages:
- pre-build
- post-build
Every stage may have list of command to run. For example:
```
hooks:
  pre-build:
    - minikube version
    - hooks/pre-build1.sh
    - python hooks/pre-build2.py
  post-build:
    - hooks/post-build1.sh
    - ruby hooks/post-build2.rb
    - rm -rf ./cache
```

_Please remember that it's easy to add OS-specific hooks like I just did with `rm -rf ./cache` and thus to complicate working with the package for other members of your team, who works with a different OS family (e.g. Windows)._

Hooks are executed in the order they are defined at the list. If any hook call is failed, the whole carbon command call is terminated.

Hooks do not support advanced commands like Bash one-line scripts. In this case you have to create a script file and call it from a hook.

Currently Carbon allows to run hooks on two stages:
- pre-build
- post-build

### Pre-Build hooks
These hooks are executed after the main config is read processed, but before Kubernetes manifest templates are processed.

### Post-Build hooks
These hooks are executed after a Docker image is built, but before it's pushed and removed (if these actions are  requested by `carbon build` parameters).

## Building a Carbon package without custom application code
Sometimes you may need to build a Carbon package without a real application code. For example, it may be necessary, when you want to use community Docker image (say, RabbitMQ) but with custom configuration. It's still acheavable with Carbon.

Taking into account that a Docker image is a foundation for a Carbon package, but we want to avoid senseless size bloating, we recommend to use `scratch` as a parent image. It's a special base image, which doesn't have any data. hHence, you'll get the most lightweight Docker image possible, without any useless data. Thus your full Dockerfile will be:
```
FROM scratch
```
It doesn't waight a single byte, but can be use for distribution of your Kubernetes manifests and allows you to use other amazing features of Carbon, like version constraints (comming soon).

## Building multiple packages from the same codebase
This approach isn't recommended. We recommend to implement multiple commands in order to start different modules/roles of your package, but to build and distribute this codebase as a single package.

If you still need to build two or more Carbon packages from the same codebase, you have such possibility. As I described at "Main config" section, carbon config has to be located in the root of a codebase, but `carbon build` has parameter `-c/--config` which allows to define path to the main config. Thus you may have multiple main Carbon configs with different names in the root of your codebase and use them for building of multiple Carbon packages

## Working with Minikube
In order to provide a simple development process, Carbon has a native integration with Minikube. All what you have to do is to install minikube and run carbon commands with `-m/--minikube` flag.
In this mode `carbon build` builds image at the Minikube docker daemon, `carbon install` reads data from this daemon and installs Carbon package to the Minikube cluster. Almost all flags (except `carbon version` of course :) ) support Minikube mode.

## Roadmap
- `carbon init` - create main Carbon configs for a new package
- `carbon verify` - verify Carbon package configuration before running `carbon build`
- `carbon config-upgrade` - upgrade main config to the latest version
- version constraints and dependencies - provide possibility to define dependencies among Carbon packages and to install a package only when all version constraints are met