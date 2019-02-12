Carbon is a package manager for Kubernetes. It allows to operate with your application along with Kubernetes manifests as a holistic package.

You are able to do with a Carbon package:
- develop
- build
- assign version (Docker tags)
- distribute (via Docker registries)
- install to Kubernetes
- and many other maintenance operations

Carbon packages are based on Docker images and don't impose any additional requirements or restrictions. Thus you may use your favorite Docker registry in order to store, share and install it to your Kubernetes cluster.

## How Carbon works
Carbon package is a Docker image with additional metadata stored at the image labels. The most important part of the metadata is a Kubernetes manifest templates, which are converted to Kubernetes manifests on the installation stage.

Usage of Docker labels allows us to avoid downloading the whole image for the installation. Carbon just reads a package metadata directly from a registry, builds a Kubernetes manifest based on the templates and applies this manifest to a Kubernetes cluster. Based on the manifest, Kubernetes downloads specified Docker images from a registry and launches Pods.

## Content
- [Available CLI operations](docs/available_cli_operations.md)
- [Carbon config](docs/carbon_config.md)
- [Structure of a Carbon package](docs/structure_of_a_carbon_package.md)
- [Package config](docs/package_config.md)
- [Kubernetes manifest templates and variables](docs/kubernetes_manifest_templates_and_variables.md)
- [Patches](docs/patches.md)
- [Docker image names and tags](docs/docker_image_names_and_tags.md)
- [Hooks](docs/hooks.md)
- [Building a Carbon package without an application code](docs/building_a_carbon_package_without_an_application_code.md)
- [Building multiple packages from a shared codebase](docs/building_multiple_pckages_from_a_shared_codebase.md)
- [Installed package metadata](docs/installed_packages_metadata.md)
- [Working with Minikube](docs/working_with_minikube.md)

TODO: Fix links. and check other TODOs
TODO: check english wording

## Limitations (TODO not finished)
- All resources of a package will be deployed to the same namespace. If you need to deploy resources to different namespaces, we recommend splitting such package to different packages
- If a package has global resources (like ClusterRole or CustomResourceDefinition), it's your responsibility to avoid multiple installations of your package in different Namespaces, especially when you use `cluster` scope (it may be improved in the future).

## Roadmap
- `carbon init` - create package Carbon configs for a new package
- `carbon verify` - verify Carbon package configuration before running `carbon build`
- `carbon config-upgrade` - upgrade a package config to the latest version
- version constraints and dependencies - provide a possibility to define dependencies among Carbon packages and to install a package only when all version constraints are met