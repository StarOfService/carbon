## Kubernetes manifest templates and variables
Kubernetes manifest templates are standard Go templates, so you can find the whole necessary information about the format at the documentation for `text/template` package: https://golang.org/pkg/text/template/

At the templates you can use two kinds of variables:
- Package properties which are filled by Carbon automatically:
- User-defined variables

### Package properties
Package properties are filled by Carbon automatically. Currently, these properties are exposed:
- `.Pkg.Name` - name of the Carbon package
- `.Pkg.Version` - version of the Carbon package
- `.Pkg.DockerName` - name of the Focker image which is used for the installation (for example _ubuntu_)
- `.Pkg.DockerTag` - tag suffix of the Docker image which is used for the installation (for example _latest_)

It's up to you how to use these properties, but here we want to give some recommendations:
- `.Pkg.Name` may be useful for naming Kubernetes resources and thus to avoid conflicts with other packages.
- combination of `.Pkg.DockerName` and `.Pkg.DockerTag` properties will let you avoid hardcoding of a Pod image

### User-defined variables
All user-defined variables must be defined at `variables` section of the main config. Each variable may have such parameters:
Parameter     | Required | Default | Description
--------------|----------|---------|------------
`name`        | true     | -       | variable name
`default`     | false    | ""      | Default value for the variable
`description` | false    | ""      | Description for the variable

*Variable name may contain alphanumeric symbols and undescores*. But we recommend to use only alphanumeric symbols with CamelCase format.

When you use these variables at a Kubernetes manifest template, they must be prefixed by `.Var.`, for example: `.Var.Environment`.

At the installation stage, variables can be taken from different sources and are applied in the following order:
1. Default values defined at the main config
2. `~/.carbon/carbon.vars` - variables file from a user home directory
3. Variables files provided by `--var-file` flag
4. Variables provided by `--var` flag

Variables file, in fact is a .property file (https://en.wikipedia.org/wiki/.properties). In most cases it's enough to use a simple key=value notation:
```
Cluster = core
Environment=Local
```
But if you have an advanced use-case, please follow the `.properties` format and remember that variables are disallowed to have spaces at name.
