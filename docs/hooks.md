## Hooks
Hooks is a way to run system calls on specific stages of a Carbon package lifecycle. It may be a direct system call or a shell script.

Every stage may have list of command to run. For example:
```
hooks:
  pre-build:
    - minikube version
    - hooks/pre-build1.sh
    - python hooks/pre-build2.py
  post-build:
    - hooks/post-build1.sh
    - ruby hooks/post-build2.rb
    - rm -rf ./cache
```

_Please remember that it's easy to add OS-specific hooks like I just did with `rm -rf ./cache` and thus to complicate working with the package for other members of your team, who works with a different OS family (e.g. Windows)._

Hooks are executed in the order they are defined at the list. If any hook call is failed, the whole Carbon command call is terminated.

Hooks do not support advanced commands like Bash one-line scripts. In this case, you have to create a script file and call it from a hook.

Currently, Carbon allows to run hooks on two stages:
- pre-build
- post-build

### Pre-Build hooks
These hooks are executed after the main config is read processed, but before Kubernetes manifest templates are processed.

### Post-Build hooks
These hooks are executed after a Docker image is built, but before it's pushed and removed (if these actions are requested by `carbon build` parameters).
