## Get started with Carbon
### Install Carbon 
(TODO - link)

### Install Minikube
Install Minikube (https://kubernetes.io/docs/tasks/tools/install-minikube/) and make sure that Minikube is started: `minikube start` & `minikube status`

### Install Docker 
https://docs.docker.com/install/

### Clone Carbon git repository
Clone Carbon git repository (https://github.com/StarOfService/carbon) and change current working directory (CWD) to _./test_ folder

### Files overview 
At the CWD you may find 2 files:
- Dockerfile - is a standard Dockerfile which will be used for building a Docker image. The one your observing now countains only `FROM scratch`. It's because we don't have any real application which has to be built into the Docker image. When you're be working with a real application, you have to add this application into the Docker image. For more details please check official reference for Dockerfile (https://docs.docker.com/engine/reference/builder/)
- carbon.yaml - is a package config which contains Carbon-specific parameters related to your packace. More details can be found at [Package config](docs/package_config.md)

Also the CWD contains 2 folder:
- k8s - in this folder you can find Kubernetes manifests. The path for manifests if configured by `kubeManifests` parameter at the package config.
- hooks - here we put lifecycle hooks. If you want to learn more about this feature, please check [Hooks](docs/hooks.md) and `hooks` parameter at the package config

### Docker image template
Please open _k8s/deployment.yaml_ file and find _containers_ section at the _Deployment_ reource specification. Pay attention to the this line:
`image: '{{.Pkg.DockerName}}:{{.Pkg.DockerTag}}'`. As you can see, here's no any real Docker image name, we use template variables instead. We'll return to this point a bit later, on an installation stage.

### Carbon package building
We've checked content of the working directory and now we're ready to build our first Carbon package. Type: `carbon build -m`. `-m` flag directs Carbon to use Minikube VM which we installed on the 2nd step. The output will be like this:
```
$ carbon build -m
INFO[2019-02-13T10:31:02+03:00] Starting Carbon build
INFO[2019-02-13T10:31:02+03:00] Reading Carbon config
INFO[2019-02-13T10:31:02+03:00] Running pre-build hook
minikube version: v0.30.0
INFO[2019-02-13T10:31:02+03:00] Building Carbon package
Step 1/2 : FROM scratch
 --->
Step 2/2 : LABEL "carbon-package-metadata"='{"apiVersion":"v1alpha1","pkgName":"carbon-test","pkgVersion":"0.0.1","buildtime":1550043062,"mainConfigB64":"YXBpVmVyc2lvbjogdjFhbHBoYTEKZG9ja2VyZmlsZTogRG9ja2VyZmlsZQprdWJlTWFuaWZlc3RzOiBrOHMvKi55YW1sCm5hbWU6IGNhcmJvbi10ZXN0CnZlcnNpb246IDAuMC4xCmhvb2tzOgogIHByZS1idWlsZDoKICAgIC0gbWluaWt1YmUgdmVyc2lvbgogIHBvc3QtYnVpbGQ6CiAgICAtIGhvb2tzL2J1aWxkLnNoCnZhcmlhYmxlczoKICAtIG5hbWU6IEZ1bGxOYW1lCiAgICBkZWZhdWx0OiBjYXJib24tdGVzdAogIC0gbmFtZTogRW52aXJvbm1lbnQKICAgIGRlZmF1bHQ6IGxvY2FsCiAgLSBuYW1lOiBLdWJlTmFtZXNwYWNlCiAgICBkZWZhdWx0OiBkZWZhdWx0CiAgLSBuYW1lOiBDbHVzdGVyCiAgICBkZWZhdWx0OiBjb3JlCiAgICBkZXNjcmlwdGlvbjogSzhzIGNsdXN0ZXIgbmFtZQ==","kubeConfigB64":"Ci0tLQphcGlWZXJzaW9uOiBhcGlleHRlbnNpb25zLms4cy5pby92MWJldGExCmtpbmQ6IEN1c3RvbVJlc291cmNlRGVmaW5pdGlvbgptZXRhZGF0YToKICBuYW1lOiBjYXJib250ZXN0cy5zdGFyb2ZzZXJ2aWNlLmNvbQpzcGVjOgogIGdyb3VwOiBzdGFyb2ZzZXJ2aWNlLmNvbQogIHZlcnNpb246IHYxYWxwaGExCiAgc2NvcGU6IE5hbWVzcGFjZWQKICBuYW1lczoKICAgIHBsdXJhbDogY2FyYm9udGVzdHMKICAgIHNpbmd1bGFyOiBjYXJib250ZXN0CiAgICBraW5kOiBDYXJib25UZXN0CgotLS0KYXBpVmVyc2lvbjogZXh0ZW5zaW9ucy92MWJldGExCmtpbmQ6IERlcGxveW1lbnQKbWV0YWRhdGE6CiAgbmFtZTogJ3t7LlBrZy5OYW1lfX0nCiAgbGFiZWxzOgogICAgYXBwOiAne3suUGtnLk5hbWV9fScKc3BlYzoKICByZXBsaWNhczogMQogIHNlbGVjdG9yOgogICAgbWF0Y2hMYWJlbHM6CiAgICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCiAgdGVtcGxhdGU6CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxzOgogICAgICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCiAgICAgIGFubm90YXRpb25zOgogICAgICAgIGlhbS5hbWF6b25hd3MuY29tL3JvbGU6ICd7ey5WYXIuRnVsbE5hbWV9fScKICAgIHNwZWM6CiAgICAgIHNlcnZpY2VBY2NvdW50TmFtZTogJ3t7LlZhci5GdWxsTmFtZX19JwogICAgICBjb250YWluZXJzOgogICAgICAgIC0gaW1hZ2U6ICd7ey5Qa2cuRG9ja2VyTmFtZX19Ont7LlBrZy5Eb2NrZXJUYWd9fScKICAgICAgICAgIG5hbWU6ICd7ey5Qa2cuTmFtZX19JwogICAgICAgICAgY29tbWFuZDoKICAgICAgICAgICAgLSAvYXBwL21haW4KICAgICAgICAgICAgLSAtYW1zX3ByZWZpeD17ey5WYXIuRW52aXJvbm1lbnR9fS17ey5WYXIuQ2x1c3Rlcn19CiAgICAgICAgICAgIC0gLWFtc19kZGJ0YWJsZT17ey5WYXIuRW52aXJvbm1lbnR9fS17ey5Qa2cuTmFtZX19CiAgICAgICAgICAgIC0gLWFtc19kZGJyZWdpb249e3tpZiAob3IgKGVxIC5WYXIuRW52aXJvbm1lbnQgInNhbmRib3giKSAoZXEgLlZhci5FbnZpcm9ubWVudCAibG9jYWwiKSl9fWV1LXdlc3QtMXt7ZWxzZX19ZXUtY2VudHJhbC0xe3tlbmR9fQogICAgICAgICAgICAtIC1sb2d0b3N0ZGVycj10cnVlCi0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZXJ2aWNlCm1ldGFkYXRhOgogIG5hbWU6ICd7ey5Qa2cuTmFtZX19JwogIGxhYmVsczoKICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCnNwZWM6CiAgc2VsZWN0b3I6CiAgICBhcHA6ICd7ey5Qa2cuTmFtZX19JwogIHBvcnRzOgogIC0gcHJvdG9jb2w6IFRDUAogICAgcG9ydDogODAKICAgIHRhcmdldFBvcnQ6IDkzNzYKCi0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZXJ2aWNlQWNjb3VudAptZXRhZGF0YToKICBuYW1lOiAne3suVmFyLkZ1bGxOYW1lfX0nCiAgbGFiZWxzOgogICAgYXBwOiAne3suUGtnLk5hbWV9fScKLS0tCmFwaVZlcnNpb246IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8vdjEKa2luZDogQ2x1c3RlclJvbGUKbWV0YWRhdGE6CiAgbmFtZTogJ3t7LlBrZy5OYW1lfX0nCiAgbGFiZWxzOgogICAgYXBwOiAne3suUGtnLk5hbWV9fScKcnVsZXM6Ci0gYXBpR3JvdXBzOgogIC0gc3Rhcm9mc2VydmljZS5jb20KICByZXNvdXJjZXM6CiAgLSBhbXNwcm9kdWNlcnMKICAtIGFtc2NvbnN1bWVycwogIHZlcmJzOgogIC0gZ2V0CiAgLSBsaXN0CiAgLSB3YXRjaAotLS0KYXBpVmVyc2lvbjogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pby92MQpraW5kOiBDbHVzdGVyUm9sZUJpbmRpbmcKbWV0YWRhdGE6CiAgbmFtZTogJ3t7LlZhci5GdWxsTmFtZX19JwogIGxhYmVsczoKICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCnJvbGVSZWY6CiAgYXBpR3JvdXA6IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8KICBraW5kOiBDbHVzdGVyUm9sZQogIG5hbWU6ICd7ey5Qa2cuTmFtZX19JwpzdWJqZWN0czoKLSBraW5kOiBTZXJ2aWNlQWNjb3VudAogIG5hbWU6ICd7ey5WYXIuRnVsbE5hbWV9fScKICBuYW1lc3BhY2U6ICd7ey5WYXIuS3ViZU5hbWVzcGFjZX19Jwo=","variables":[{"name":"FullName","default":"carbon-test","description":""},{"name":"Environment","default":"local","description":""},{"name":"KubeNamespace","default":"default","description":""},{"name":"Cluster","default":"core","description":"K8s cluster name"}]}'
 ---> Running in 5c8c65c8275a
Removing intermediate container 5c8c65c8275a
 ---> 9675eecf6564
Successfully built 9675eecf6564
Successfully tagged carbon-test:0.0.1
INFO[2019-02-13T10:31:02+03:00] Running post-build hook
INFO[2019-02-13T10:31:02+03:00] Carbon package has been built successfully
```
You may notice that the significant part of the output is a standard output for `docker build`, but on the second step there's added a huge label. It's a Carbon package metadata which contains some parameters from the package config and base64-encoded Kubernetes manifest templates. Now you can see how Carbon extends Docker image in order to store Kubernetes manifests.

### Docker image inspection
You can do all standard Docker operations with the Carbon package, for example inspect.
Due to the `-m` flag, the image was built at the Minikube VM. In order to configure our local Docker client to use Minikube docker daemon, minikube provides a command `minikube docker-env`. Please type this command and follow the instructions for your OS.
When the Docker client is configured, we can check the docker image details:
```
$ docker inspect carbon-test:0.0.1
[
    {
        "Id": "sha256:9675eecf6564333640dd0ba2ae4b0c4094f370f5425d1b7b32e65809a7214494",
        "RepoTags": [
            "carbon-test:0.0.1"
        ],
        "RepoDigests": [],
        "Parent": "",
...
            "Labels": {
                "carbon-package-metadata": "{\"apiVersion\":\"v1alpha1\",\"pkgName\":\"carbon-test\",\"pkgVersion\":\"0.0.1\",\"buildtime\":1550043062,\"mainConfigB64\":\"YXBpVmVyc2lvbjogdjFhbHBoYTEKZG9ja2VyZmlsZTogRG9ja2VyZmlsZQprdWJlTWFuaWZlc3RzOiBrOHMvKi55YW1sCm5hbWU6IGNhcmJvbi10ZXN0CnZlcnNpb246IDAuMC4xCmhvb2tzOgogIHByZS1idWlsZDoKICAgIC0gbWluaWt1YmUgdmVyc2lvbgogIHBvc3QtYnVpbGQ6CiAgICAtIGhvb2tzL2J1aWxkLnNoCnZhcmlhYmxlczoKICAtIG5hbWU6IEZ1bGxOYW1lCiAgICBkZWZhdWx0OiBjYXJib24tdGVzdAogIC0gbmFtZTogRW52aXJvbm1lbnQKICAgIGRlZmF1bHQ6IGxvY2FsCiAgLSBuYW1lOiBLdWJlTmFtZXNwYWNlCiAgICBkZWZhdWx0OiBkZWZhdWx0CiAgLSBuYW1lOiBDbHVzdGVyCiAgICBkZWZhdWx0OiBjb3JlCiAgICBkZXNjcmlwdGlvbjogSzhzIGNsdXN0ZXIgbmFtZQ==\",\"kubeConfigB64\":\"Ci0tLQphcGlWZXJzaW9uOiBhcGlleHRlbnNpb25zLms4cy5pby92MWJldGExCmtpbmQ6IEN1c3RvbVJlc291cmNlRGVmaW5pdGlvbgptZXRhZGF0YToKICBuYW1lOiBjYXJib250ZXN0cy5zdGFyb2ZzZXJ2aWNlLmNvbQpzcGVjOgogIGdyb3VwOiBzdGFyb2ZzZXJ2aWNlLmNvbQogIHZlcnNpb246IHYxYWxwaGExCiAgc2NvcGU6IE5hbWVzcGFjZWQKICBuYW1lczoKICAgIHBsdXJhbDogY2FyYm9udGVzdHMKICAgIHNpbmd1bGFyOiBjYXJib250ZXN0CiAgICBraW5kOiBDYXJib25UZXN0CgotLS0KYXBpVmVyc2lvbjogZXh0ZW5zaW9ucy92MWJldGExCmtpbmQ6IERlcGxveW1lbnQKbWV0YWRhdGE6CiAgbmFtZTogJ3t7LlBrZy5OYW1lfX0nCiAgbGFiZWxzOgogICAgYXBwOiAne3suUGtnLk5hbWV9fScKc3BlYzoKICByZXBsaWNhczogMQogIHNlbGVjdG9yOgogICAgbWF0Y2hMYWJlbHM6CiAgICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCiAgdGVtcGxhdGU6CiAgICBtZXRhZGF0YToKICAgICAgbGFiZWxzOgogICAgICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCiAgICAgIGFubm90YXRpb25zOgogICAgICAgIGlhbS5hbWF6b25hd3MuY29tL3JvbGU6ICd7ey5WYXIuRnVsbE5hbWV9fScKICAgIHNwZWM6CiAgICAgIHNlcnZpY2VBY2NvdW50TmFtZTogJ3t7LlZhci5GdWxsTmFtZX19JwogICAgICBjb250YWluZXJzOgogICAgICAgIC0gaW1hZ2U6ICd7ey5Qa2cuRG9ja2VyTmFtZX19Ont7LlBrZy5Eb2NrZXJUYWd9fScKICAgICAgICAgIG5hbWU6ICd7ey5Qa2cuTmFtZX19JwogICAgICAgICAgY29tbWFuZDoKICAgICAgICAgICAgLSAvYXBwL21haW4KICAgICAgICAgICAgLSAtYW1zX3ByZWZpeD17ey5WYXIuRW52aXJvbm1lbnR9fS17ey5WYXIuQ2x1c3Rlcn19CiAgICAgICAgICAgIC0gLWFtc19kZGJ0YWJsZT17ey5WYXIuRW52aXJvbm1lbnR9fS17ey5Qa2cuTmFtZX19CiAgICAgICAgICAgIC0gLWFtc19kZGJyZWdpb249e3tpZiAob3IgKGVxIC5WYXIuRW52aXJvbm1lbnQgInNhbmRib3giKSAoZXEgLlZhci5FbnZpcm9ubWVudCAibG9jYWwiKSl9fWV1LXdlc3QtMXt7ZWxzZX19ZXUtY2VudHJhbC0xe3tlbmR9fQogICAgICAgICAgICAtIC1sb2d0b3N0ZGVycj10cnVlCi0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZXJ2aWNlCm1ldGFkYXRhOgogIG5hbWU6ICd7ey5Qa2cuTmFtZX19JwogIGxhYmVsczoKICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCnNwZWM6CiAgc2VsZWN0b3I6CiAgICBhcHA6ICd7ey5Qa2cuTmFtZX19JwogIHBvcnRzOgogIC0gcHJvdG9jb2w6IFRDUAogICAgcG9ydDogODAKICAgIHRhcmdldFBvcnQ6IDkzNzYKCi0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZXJ2aWNlQWNjb3VudAptZXRhZGF0YToKICBuYW1lOiAne3suVmFyLkZ1bGxOYW1lfX0nCiAgbGFiZWxzOgogICAgYXBwOiAne3suUGtnLk5hbWV9fScKLS0tCmFwaVZlcnNpb246IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8vdjEKa2luZDogQ2x1c3RlclJvbGUKbWV0YWRhdGE6CiAgbmFtZTogJ3t7LlBrZy5OYW1lfX0nCiAgbGFiZWxzOgogICAgYXBwOiAne3suUGtnLk5hbWV9fScKcnVsZXM6Ci0gYXBpR3JvdXBzOgogIC0gc3Rhcm9mc2VydmljZS5jb20KICByZXNvdXJjZXM6CiAgLSBhbXNwcm9kdWNlcnMKICAtIGFtc2NvbnN1bWVycwogIHZlcmJzOgogIC0gZ2V0CiAgLSBsaXN0CiAgLSB3YXRjaAotLS0KYXBpVmVyc2lvbjogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pby92MQpraW5kOiBDbHVzdGVyUm9sZUJpbmRpbmcKbWV0YWRhdGE6CiAgbmFtZTogJ3t7LlZhci5GdWxsTmFtZX19JwogIGxhYmVsczoKICAgIGFwcDogJ3t7LlBrZy5OYW1lfX0nCnJvbGVSZWY6CiAgYXBpR3JvdXA6IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8KICBraW5kOiBDbHVzdGVyUm9sZQogIG5hbWU6ICd7ey5Qa2cuTmFtZX19JwpzdWJqZWN0czoKLSBraW5kOiBTZXJ2aWNlQWNjb3VudAogIG5hbWU6ICd7ey5WYXIuRnVsbE5hbWV9fScKICBuYW1lc3BhY2U6ICd7ey5WYXIuS3ViZU5hbWVzcGFjZX19Jwo=\",\"variables\":[{\"name\":\"FullName\",\"default\":\"carbon-test\",\"description\":\"\"},{\"name\":\"Environment\",\"default\":\"local\",\"description\":\"\"},{\"name\":\"KubeNamespace\",\"default\":\"default\",\"description\":\"\"},{\"name\":\"Cluster\",\"default\":\"core\",\"description\":\"K8s cluster name\"}]}"
            }
...
        "RootFS": {
            "Type": "layers"
        },
        "Metadata": {
            "LastTagTime": "2019-02-13T07:31:03.381705852Z"
        }
    }
]
```
As you can see, it's a standard Docker image, which just contains an additional label. Hence you can operate with this image as any other Docker iamge: pull, push, store on a Docker registry, use other third-party tools.

### Carbon package inspection
Let's check the information for the Carbon package:
```
$ carbon inspect -m carbon-test:0.0.1
Name: carbon-test
Version: 0.0.1
Variables:
  NAME           DEFAULT      DESCRIPTION
  FullName       carbon-test
  Environment    local
  KubeNamespace  default
  Cluster        core         K8s cluster name
```
Here Carbon reads the package metadata from the docker label and exposes basic information about the package:
- package name
- package version
- package variables and default vaues for those variables

### Carbon package installation
Now let's install the package to a Kubernetes cluster:
```
$ carbon install -m carbon-test:0.0.1
INFO[2019-02-13T11:06:44+03:00] Starting Carbon install
INFO[2019-02-13T11:06:44+03:00] Getting Carbon package metadata
INFO[2019-02-13T11:06:44+03:00] Building Kubernetes configuration
INFO[2019-02-13T11:06:44+03:00] Applying patches
INFO[2019-02-13T11:06:44+03:00] Applying kubernetes configuration
customresourcedefinition.apiextensions.k8s.io/carbontests.starofservice.com created
deployment.extensions/carbon-test created
service/carbon-test created
serviceaccount/carbon-test created
clusterrole.rbac.authorization.k8s.io/carbon-test created
clusterrolebinding.rbac.authorization.k8s.io/carbon-test created
INFO[2019-02-13T11:06:53+03:00] Saving Carbon package metadata
INFO[2019-02-13T11:06:54+03:00] Carbon package has been installed successfully
```
At the [Docker image template](#docker-image-template) section I showed a template for the _image_ field. The content for those variables is taken from the provided arguments. In our case it's `carbon-test:0.0.1`. It works fine for Minikube, but for a real environment you want to push your image to a remote registry like docker-hub. When you are installing such Carbon package, you provide the remote image identifier. This identifier is passed to the templates and thus Kubernetes uses the same docker image which is used for the installation.
    
Feel free to use `kubectl` in order to check just created resources.

### Installed Carbon packages overview
Now we can check what Carbon packages are installed to the Kubernetes cluster:
```
$ carbon status -m
  NAMESPACE  NAME         VERSION  SOURCE
  default    carbon-test  0.0.1    docker.io/library/carbon-test:0.0.1
```

### Installed Carbon package details
If we need a detailed information for a specific package, we should provide its name as an argument to the `carbon status` command:
```
$ carbon status -m carbon-test
Namespace: default
Name: carbon-test
Version: 0.0.1
Source: docker.io/library/carbon-test:0.0.1
Variables:
  NAME           VALUE
  Cluster        core
  Environment    local
  FullName       carbon-test
  KubeNamespace  default
```
We didn't override variables during the installation, that's why here we can see the default values. In fact, current command shows `real` values which were used for the package installation.

`carbon status` has an additional `-f/--full` flag, which allows to see the maximum available information about the package, including the eventual Kubernetes manifests:
```
$ carbon status -m carbon-test -f
Namespace: default
Name: carbon-test
Version: 0.0.1
Source: docker.io/library/carbon-test:0.0.1
Variables:
  NAME           VALUE
  Cluster        core
  Environment    local
  FullName       carbon-test
  KubeNamespace  default
Patches:
Manifest: {"apiVersion":"apiextensions.k8s.io/v1beta1","kind":"CustomResourceDefinition","metadata":{"labels":{"carbon/component-name":"carbon-test","carbon/component-version":"0.0.1","managed-by":"carbon"},"name":"carbontests.starofservice.com"},"spec":{"group":"starofservice.com","names":{"kind":"CarbonTest","plural":"carbontests","singular":"carbontest"},"scope":"Namespaced","version":"v1alpha1"}}{"apiVersion":"extensions/v1beta1","kind":"Deployment","metadata":{"labels":{"app":"carbon-test","carbon/component-name":"carbon-test","carbon/component-version":"0.0.1","managed-by":"carbon"},"name":"carbon-test"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"carbon-test"}},"template":{"metadata":{"annotations":{"iam.amazonaws.com/role":"carbon-test"},"labels":{"app":"carbon-test"}},"spec":{"containers":[{"command":["/app/main","-ams_prefix=local-core","-ams_ddbtable=local-carbon-test","-ams_ddbregion=eu-west-1","-logtostderr=true"],"image":"docker.io/library/carbon-test:0.0.1","name":"carbon-test"}],"serviceAccountName":"carbon-test"}}}}{"apiVersion":"v1","kind":"Service","metadata":{"labels":{"app":"carbon-test","carbon/component-name":"carbon-test","carbon/component-version":"0.0.1","managed-by":"carbon"},"name":"carbon-test"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":9376}],"selector":{"app":"carbon-test"}}}{"apiVersion":"v1","kind":"ServiceAccount","metadata":{"labels":{"app":"carbon-test","carbon/component-name":"carbon-test","carbon/component-version":"0.0.1","managed-by":"carbon"},"name":"carbon-test"}}{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"ClusterRole","metadata":{"labels":{"app":"carbon-test","carbon/component-name":"carbon-test","carbon/component-version":"0.0.1","managed-by":"carbon"},"name":"carbon-test"},"rules":[{"apiGroups":["starofservice.com"],"resources":["amsproducers","amsconsumers"],"verbs":["get","list","watch"]}]}{"apiVersion":"rbac.authorization.k8s.io/v1","kind":"ClusterRoleBinding","metadata":{"labels":{"app":"carbon-test","carbon/component-name":"carbon-test","carbon/component-version":"0.0.1","managed-by":"carbon"},"name":"carbon-test"},"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"ClusterRole","name":"carbon-test"},"subjects":[{"kind":"ServiceAccount","name":"carbon-test","namespace":"default"}]}
```

### Installed Carbon package uninstallation
Finally let's clean-up the environment from this testing package:
```
$ carbon uninstall -m carbon-test
INFO[2019-02-13T11:15:36+03:00] Uninstalling Carbon packagecarbon-test
customresourcedefinition.apiextensions.k8s.io "carbontests.starofservice.com" deleted
deployment.extensions "carbon-test" deleted
service "carbon-test" deleted
serviceaccount "carbon-test" deleted
clusterrole.rbac.authorization.k8s.io "carbon-test" deleted
clusterrolebinding.rbac.authorization.k8s.io "carbon-test" deleted
INFO[2019-02-13T11:15:37+03:00] Carbon packages has been uninstalled successfully
```

### Summary
As we just showed, Carbon allows to operate (build, distribute and install) with your application and kubernetes manifest as a single package. Now you are able to build a Carbon package, assign version, test a specific version and be sure that when you will be installing the package to your production evnironment, it will be absolutely the same package (with the same application code and kubernetes manifests) which was tested before.

If you want to learn more details about its features and how to build a Carbon package, please return to the [documentation index](../README.md) and check pages dedicated to different aspects of the current tool.