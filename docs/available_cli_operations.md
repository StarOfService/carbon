## Available CLI operations
- `carbon build` - build Carbon package. It builds a Docker image with Carbon metadata. As soon as it's built, you can operate with this Carbon package as you do with any Docker image: assign additional tags, run `docker inspect`, push to a Docker registry and so on.
- `carbon inspect dockerImage [dockerImage ...]` - this command exposes Carbon package metadata from the given Docker images. It can either local or remote images. If it's a private repository, you have to be authenticated and authorized in order to work with this repo (see `docker login`). The command shows Carbon *name*, *version* and *variables* for the requested packages.
- `carbon install dockerImage` - reads Carbon metadata from a given Docker image, builds Kubernetes manifest based on variables, applies patches if any and deploys the resulting manifest to Kubernetes.
- `carbon status [packageName [packageName ...]]` - if no packages are provided, Carbon lists all installed packages for a current Kubernetes context. When at least one package name is provided, Carbon shows detailed information for the Carbon package. A `--full` flag allows seeing even more extended information like the patches used on the installation stage and an eventual Kubernetes manifest applied to a Kubernetes cluster.
- `carbon uninstall packageName [packageName ...]` - uninstall specified packages
- `carbon version` - show version for your Carbon CLI tool and the latest supported version of the main config.

Almost all mentioned subcommands have additional flags which extend the described behavior. So don't hesitate to use CLI built-in help, for example: `carbon build -h` or `carbon help build`
