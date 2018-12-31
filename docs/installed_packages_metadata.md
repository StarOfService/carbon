## Installed package metadata
Carbon stores information about installed packages at Kubernetes secrets prefixed by `carbon-package-metadata-`. By default, it uses `carbon-data` namespace, but this behavior can be modified through `scope` parameter of a [Carbon config](carbon_config.md)

### Scope
Carbon scope defines where metadata for installed packages. Current parameter supports two values:
- cluster
- namespace

#### Cluster
This scope is used by default. With this scope, Carbon stores metadata for all installed packages at the `carbon-data` namespace. Thus you may install packages to the different namespaces and to use a central registry metadata in order to make sure which packages and versions are installed. It's convenient when components from different namespaces have to interact with each other.

#### Namespace
When namespace scope is configured, Carbon stores a package and the package metadata to the same namespace. It allows isolating Carbon metadata per namespace. Thus the current option is useful when you have isolated namespaces and components from different namespaces don't interact with each other.
