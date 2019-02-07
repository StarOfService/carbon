Carbon is a package manager for Kubernetes. It allows to operate with your application along with Kubernetes manifests as a holistic package.

You are able to do with a Carbon package:
- develop
- build
- assign version (docker tags)
- distribute (via docker registries)
- install to Kubernetes
- and many other maintenance operations

Carbon packages are based on Docker images and don't impose any additional requirements or restrictions. Thus you may use your favorite Docker registry in order to store, share and install it to your Kubernetes cluster.

## How Carbon works
Carbon package is a docker image with additional metadata stored at the image lables. The most important part of the metadata is a Kubernetes manifest templates, which are converted to Kubernetes manifests on the installation stage.

Usage of docker lables allows us to avoid downloading of the whole image for the installation. Carbon just reads a package metadata directly from a registry, builds a Kubernetes manifest based on the templates and applies this manifest to a Kubernetes cluster. Based on the manifest, Kubernetes downloads specified Docker images from a regestry and launches Pods.

## Content
- (Available CLI operations)[docs/available_cli_operations.md]
- (Structure of a Carbon package)[docs/structure_of_a_carbon_package.md]
- (Main config)[docs/main_config.md]
- (Kubernetes manifest templates and variables)[docs/kubernetes_manifest_templates_and_variables.md]
- (Patches)[docs/patches.md]
- (Docker image names and tags)[docs/docker_image_names_and_tags.md]
- (Hooks)[docs/hooks.md]
- (Building a Carbon package without an application code)[docs/building_a_carbon_package_without_an_application_code.md]
- (Building multiple packages from a shared codebase)[docs/building_multiple_pckages_from_a_shared_codebase.md]
- (Working with Minikube)[docs/working_with_minikube.md]
TODO: Fix links. and check other TODOs
TODO: check english wording

## Roadmap
- `carbon init` - create main Carbon configs for a new package
- `carbon verify` - verify Carbon package configuration before running `carbon build`
- `carbon config-upgrade` - upgrade main config to the latest version
- version constraints and dependencies - provide possibility to define dependencies among Carbon packages and to install a package only when all version constraints are met