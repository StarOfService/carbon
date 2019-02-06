package kubernetes

import (
  "os"
  "path/filepath"

  "k8s.io/cli-runtime/pkg/genericclioptions"
  restclient "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
  "github.com/pkg/errors"

  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/util/homedir"
)

func GetKubeConfig() (*restclient.Config, error) {
  configPath, err := kubeCinfigPath()
  if err != nil {
    return nil, err
  }

  if minikube.Enabled {
    return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
      &clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
      &clientcmd.ConfigOverrides{CurrentContext: minikube.K8sContext},
    ).ClientConfig()
  }

  return clientcmd.BuildConfigFromFlags("", configPath)
}

func kubeCinfigPath() (string, error) {
  if home := homedir.Path(); home != "" {
    return filepath.Join(home, ".kube", "config"), nil
  }
  return "", errors.New("Unable to discover a user home directory")
}

func KubeCmdFactory(ns string) (cmdutil.Factory, genericclioptions.IOStreams) {
  kubeConfigFlags := genericclioptions.NewConfigFlags()

  kubeConfigFlags.Namespace = &ns

  if minikube.Enabled {
    context := minikube.K8sContext
    kubeConfigFlags.Context = &context
  }

  matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

  f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
  ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

  return f, ioStreams
}
