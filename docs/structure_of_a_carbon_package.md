## Structure of a Carbon package
Every Carbon package must contain:
- *Carbon main config*: a YAML file, located at the root of a package codebase and describes different parameters of a Carbon package
- *Dockerfile*: a standard Dockerfile, which describes how to build a docker image with your application code
- *Kubernetes manifest templates*: Golang templates, which describe Kubernetes resources

Normally your Carbon package should also contains an application code, which is built into a Docker image usting Dockerfile. But it isn't a requirement (see [Building a Carbon package without an application code](./building_a_carbon_package_without_an_application_code.md) section)

*Thus you have a signle codebase which contains your application, Dockerfile and kubernetes manifests. This codebase can be developed, versioned, built, distributed and installed as a holistic package*.
