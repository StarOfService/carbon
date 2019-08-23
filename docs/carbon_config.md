## Carbon config
Carbon config allows configuring some aspects of the Carbon behavior. Currently, Carbon config supports only one parameter (besides apiVersion):

Parameter     | Required | Default value  | Description
--------------|----------|----------------|------------
apiVersion    | true     |                | version of the config format. The latest version is `v1alpha1`
carbonScope   | false    | cluster        | this parameter defines the scope of the Carbon metadata. Allowed values are: `cluster` and `namespace`.

Carbon config can be defined in multiple places:
1. `carbon-data` namespace
2. your active namespace

In the future, we're going to add the third option - locally on a workstation.

When Carbon config is defined at Kubernetes (currently it's a single option), it has to be defined at a ConfigMap with name `carbon-config` under the data key `config`. It can be defined using JSON or YAML format. This is an example of a correct ConfigMap with the Carbon config:
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: carbon-config
data:
  config: |
    apiVersion: v1alpha1
    carbonScope: namespace
```

As soon as Carbon config has only one parameter, we didn't implement merging of parameters from different sources. But it's what we're going to do in the future as well.

### Carbon scope
Carbon scope controlls different aspects of a package installation. Currently it supports two modes:
- cluster
- namespace

#### 'cluster' Carbon scope
- This scope is used by default.
- Metadata for installed packages is stored at the `carbon-data` namespace. More details about pagckages metadata can be found at [Installed package metadata](installed_packages_metadata.md).
- If a resource has predefined namespace, this resource will be intalled to the corresponding namespace. It can be any existing namespace.
- If a resource doesn't have a configured namespace, it will be installed to the namespace managed by `--namespace` command line argument (`default` by default).
- You are not allowed to install multiple copies of the same Carbon package to the same Kubernetes cluster.

#### 'namespace' Carbon scope
- Metadata for installed packages is stored at the same namespace where the package is installed. More details about pagckages metadata can be found at [Installed package metadata](installed_packages_metadata.md).
- If a resource has predefined namespace and this NS is different from the one, managed by `--namespace` command line argument (`default` by default), the predefined NS will be overriten and a warning message will be exposed.
- If a resource doesn't have a configured namespace, it will be installed to the namespace managed by `--namespace` command line argument (`default` by default).
- Carbon applies additional `carbon/component-namespace` lable for all resources in order to avoid conflicts when Cluster-scoped resources are used.
- You are alloed to install multiple copies of the same Carbon package to different namespaces. But it's your responsibility to provide uniq names for cluster-scoped resources.

#### Limitations
- Carbon doesn't delete resources from a namespace when none resource in the NS is applied. It's a limitation of kubectl. In the future we're going to handle cleaning-up process by ourselves and it such will resolve the issue.
- It's impossible to migrate from one Carbon scope to another at the moment.
