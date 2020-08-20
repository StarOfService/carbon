## Package config
Package config must be located in the root of a package codebase. By default, Carbon is looking for the config at `carbon.yaml` file. A package config can be defined using YAML or JSON format. Every Carbon package must contain the package config.

This config supports such the fields:

Parameter     | Required | Default value  | Description
--------------|----------|----------------|------------
apiVersion    | true     |                | version of the config format. The latest version is `v1alpha1`
name          | true     |                | name of the Carbon package. Can be overriden with `--name` parameter on a build stage
version       | true     |                | version of the Carbon package. Can be extended by `--version-prefix` and `--version-suffix` parameters on a buid stage
kubeManifests | true     |                | path to Kubernetes manifest templates. Current parameter may contain wildcard, for example `./k8s/*.yaml`.
dockerfile    | false    | Dockerfile     | path to a Dockerfile which will be used for building a Docker image with your application
artifacts     | false    | _name:version_ | list of Docker tags to be assigned for the built image. This list may contain a full tag like `carbontest:latest` or just a name without a tag suffix. When a tag suffix isn't provided, the package version will be used as a suffix. Thus if `version` is set to `0.1.0` and artifacts have item `carbontest`, the image will get a tag `carbontest:0.1.0`
variables     | false    |                | List of all variables used at a Kubernetes manifest templates must be defined here. More details you can find at [Kubernetes manifest templates and variables](./kubernetes_manifest_templates_and_variables.md) section
hooks         | false    |                | Here you can define lifecycle hooks, a system calls which have to run on a certain stage of a package life. More details can be found at [Hooks](./hooks.md) section

Here's an example of a package config:
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

Currently, the latest `apiVersion` is `v1alpha1`, but you shouldn't worry about migration to the new version when it's released. Carbon versioning is based on [VConf](https://github.com/StarOfService/vconf) library which allows translating any old version up to the latest seamlessly for a user. The same approach is used for Carbon package metadata and installed packages metadata. In the future, we're going to add `carbon config-upgrade` command which will let you upgrade an old version of your package config up to the latest one by one command call.
