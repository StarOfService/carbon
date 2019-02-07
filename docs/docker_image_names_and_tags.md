## Docker image names and tags
You have several optins to define docker image name and tag:
1. Don't do anything, in this case docker image will be named by Carbon package name and version, for example: "carbontest:0.0.1"
2. Define names and tags at `artifacts` section of the main config
3. Define names and tags by `--tag` flag for `carbon build`. `--tag` flag can be used multiple times in order to assigne multiple names/tags. This option has precedence over the `artifacts` section from the main config.

When you define a name without a tag, a Carbon package version will be used as a tag.

`carbon build` command has additional flags `--tag-prefix` and `--tag-suffix`. When any of this flag is used, all docker iamge tags will be prfixed/suffixed by the provided strings. Here's one exception: prefix and siffux are not applied for the `latest` tag.

Let's consider an example. We have a Carbon package with version `0.1.0` and run the command:
```
$ carbon build \
  --tag region1.registry.example.com/carbontest \
  --tag region1.registry.example.com/carbontest:foo \
  --tag region1.registry.example.com/carbontest:latest \
  --tag region2.registry.example.com/carbontest \
  --tag region2.registry.example.com/carbontest:foo \
  --tag region2.registry.example.com/carbontest:latest \
  --tag-prefix "hotfix-" \
  --taf-suffix "-alpha"
```

This command will create a docker image with such tags:
- region1.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region1.registry.example.com/carbontest:hotfix-foo-alpha
- region1.registry.example.com/carbontest:latest
- region2.registry.example.com/carbontest:hotfix-0.1.0-alpha
- region2.registry.example.com/carbontest:hotfix-foo-alpha
- region2.registry.example.com/carbontest:latest
