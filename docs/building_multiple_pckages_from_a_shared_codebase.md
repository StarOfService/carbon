## Building multiple packages from a shared codebase
This approach isn't recommended. We recommend implementing multiple commands in order to start different modules/roles of your package but to build and distribute this codebase as a single package.

If you still need to build two or more Carbon packages from the same codebase, you have such a possibility. As I described at [Package config](docs/package_config.md), Carbon config has to be located in the root of a codebase, but `carbon build` has a parameter `-c/--config` which defines a path to a package config. Thus you may have multiple package configs with different names in the root of your codebase and use them for building of multiple Carbon packages.
