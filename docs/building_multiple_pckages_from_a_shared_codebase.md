## Building multiple packages from a shared codebase
This approach isn't recommended. We recommend to implement multiple commands in order to start different modules/roles of your package, but to build and distribute this codebase as a single package.

If you still need to build two or more Carbon packages from the same codebase, you have such possibility. As I described at "Main config" section, carbon config has to be located in the root of a codebase, but `carbon build` has parameter `-c/--config` which allows to define path to the main config. Thus you may have multiple main Carbon configs with different names in the root of your codebase and use them for building of multiple Carbon packages.
