## Docker image names and tags
You have several options to define a Docker image name and tag:
1. Don't do anything, in this case, the Docker image will be named by Carbon package name and version, for example: "carbontest:0.0.1"
2. Define names and tags at `artifacts` section of the main config
3. Define names and tags by `--tag` flag for `carbon build`. `--tag` flag can be used multiple times in order to assign multiple names/tags. This option has precedence over the `artifacts` section from the main config.

When you define a name without a tag, a Carbon package version will be used as a tag.

`carbon build` command has additional flags `--version-prefix` and `--version-suffix`. When any of this flag is used, all Docker image tags will be prefixed/suffixed by the provided strings. Here's one exception: prefix and suffix are not applied for the `latest` tag. Besides Docker tags, defined prefix and suffix will be also added to a Carbon package version defined at `carbon.yaml` file.

Let's consider an example. We have a Carbon package with version `0.1.0` and run the command:
```
$ carbon build \
  --tag region1.registry.example.com/carbontest \
  --tag region1.registry.example.com/carbontest:foo \
  --tag region1.registry.example.com/carbontest:latest \
  --tag region2.registry.example.com/carbontest \
  --tag region2.registry.example.com/carbontest:foo \
  --tag region2.registry.example.com/carbontest:latest \
  --version-prefix "hotfix-" \
  --version-suffix "-alpha"
```

This command will create a Docker image with such tags:
- region1.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region1.registry.example.com/carbontest:hotfix-foo-alpha
- region1.registry.example.com/carbontest:latest
- region2.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region2.registry.example.com/carbontest:hotfix-foo-alpha
- region2.registry.example.com/carbontest:latest
