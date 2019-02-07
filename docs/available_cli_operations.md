## Available CLI operations
- `carbon build` - build Carbon package. It builds a Docker image with a Carbon metadata. As soon as it's built, you can operate with this Carbon package like you do with any Docker image: assign additional tags, run `docker inspect`, push to a docker registry and so on.
- `carbon inspect dockerImage [dockerImage ...]` - this command exposes Carbon package metadata from the given Docker images. It can either local or remote images. If it's a private repository, you have to be authenticated and authorized in order to working with this repo (see `docker login`). The command shows Carbon *name*, *version* and *variables* for the requested packages.
- `carbon install dockerImage` - reads Carbon metadata from a given Docker image, builds kubernetes manifest based on variables, applies patches if any and deployes the resulting manifest to Kubernetes.
- `carbon status [packageName [packageName ...]]` - if no packages are provided, Carbon lists all installed packages for a current Kubernetes context. When at least one package name is provided, Carbon shows a detailed information for the Carbon package. A `--full` flag allows to see even more extended information like the pathces used on the installation stage and an eventual Kubernetes manifest applied to a Kubernetes cluster.
- `carbon delete packageName [packageName ...]` - delete specified packages
- `carbon version` - show version for your capbon CLI tool and the latest supported version of the main config.

Almost all mentioned subcommands have additional flags which extend the described behavior. So don't hestitate to use CLI built-in help, for example: `carbon build -h` or `carbon help build`
