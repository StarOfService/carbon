## Working with Minikube
In order to provide a simple development process, Carbon has a native integration with Minikube. All what you have to do is to install minikube and run carbon commands with `-m/--minikube` flag.

In this mode `carbon build` builds image at the Minikube docker daemon, `carbon install` reads data from this daemon and installs Carbon package to the Minikube cluster. Almost all flags (except `carbon version` of course :) ) support Minikube mode.