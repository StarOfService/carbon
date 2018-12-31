## Working with Minikube
In order to provide a simple development process, Carbon has native integration with Minikube. All that you have to do is to install Minikube and run Carbon commands with `-m/--minikube` flag.

In this mode `carbon build` builds an image at the Minikube Docker daemon, `carbon install` reads data from this daemon and installs Carbon package to the Minikube cluster. Almost all flags (except `carbon version` of course :) ) support Minikube mode.