#!/bin/bash

set -ex

KUBECTL_VERSION='v1.12.5'
MINIKUBE_VERSION='v0.33.0'

mount --make-rshared /

curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl
chmod +x kubectl
mv kubectl /usr/local/bin/

curl -Lo minikube https://storage.googleapis.com/minikube/releases/${MINIKUBE_VERSION}/minikube-linux-amd64
chmod +x minikube
mv minikube /usr/local/bin/

export CHANGE_MINIKUBE_NONE_USER=true
minikube start --vm-driver=none --bootstrapper=kubeadm --kubernetes-version=$KUBECTL_VERSION

minikube update-context

JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'
until kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1; done

#curl -L https://git.io/vp6lP | sudo sh
#make codestyle

export CARBON_TEST_MINIKUBE_NONE_DRIVER=true
make test

make build-all-os
