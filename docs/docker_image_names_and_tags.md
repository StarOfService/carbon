## Docker image names and tags
You have several options to define a Docker image name and tag:
1. Don't do anything, in this case, the Docker image will be named by Carbon package name and version, for example: "carbontest:0.0.1"
2. Define names and tags at `artifacts` section of the main config
3. Define names and tags by `--tag` flag for `carbon build`. `--tag` flag can be used multiple times in order to assign multiple names/tags. This option has precedence over the `artifacts` section from the main config.

When you define a name without a tag, a Carbon package version will be used as a tag.

Let's consider an example. We have a Carbon package with version `0.1.0` and run the command:
```
$ carbon build \
  --tag region1.registry.example.com/carbontest \
  --tag region1.registry.example.com/carbontest:foo \
  --tag region1.registry.example.com/carbontest:latest \
  --tag region2.registry.example.com/carbontest \
  --tag region2.registry.example.com/carbontest:foo \
  --tag region2.registry.example.com/carbontest:latest
```

This command will create a Docker image with such tags:
- region1.registry.example.com/carbontest:0.1.0
- region1.registry.example.com/carbontest:foo
- region1.registry.example.com/carbontest:latest
- region2.registry.example.com/carbontest:0.1.0
- region2.registry.example.com/carbontest:foo
- region2.registry.example.com/carbontest:latest

`carbon build` command also has `--version-prefix` and `--version-suffix` parameters which allow to extend Carbon package version. Here's an exaple of these parameters usage and how they affect docker image tag:
```
$ carbon build \
  --tag region1.registry.example.com/carbontest \
  --tag region1.registry.example.com/carbontest:foo \
  --tag region1.registry.example.com/carbontest:latest \
  --version-prefix "hotfix-" \
  --version-suffix "-alpha"
```

This command will create a Docker image with such tags:
- region1.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region1.registry.example.com/carbontest:foo
- region1.registry.example.com/carbontest:latest

