## Carbon config
Carbon config allows configuring some aspects of the Carbon behavior. Currently, Carbon config supports only one parameter (besides apiVersion):

Parameter     | Required | Default value  | Description
--------------|----------|----------------|------------
apiVersion    | true     |                | version of the config format. The latest version is `v1alpha1`
scope         | false    | cluster        | this parameter defines the scope of the Carbon metadata. Allowed values are: `cluster` and `namespace`. More details can be found at [Installed package metadata](installed_packages_metadata.md)

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
