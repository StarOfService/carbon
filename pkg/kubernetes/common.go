package kubernetes

import (
  "os"
  "path/filepath"

  log "github.com/sirupsen/logrus"
  apicorev1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/cli-runtime/pkg/genericclioptions"
  typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  restclient "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
  "github.com/pkg/errors"

  "github.com/starofservice/carbon/pkg/minikube"
  "github.com/starofservice/carbon/pkg/util/homedir"
)

const (
  GlobalCarbonNamespace = "carbon-data"
)

var CurrentNamespace string

func KubeConfig() (*restclient.Config, error) {
  cliCmdCfg, err := clientCmdConfig()
  if err != nil {
    return nil, err
  }

  return cliCmdCfg.ClientConfig()
}

func SetNamespace(cliNS string) error {
  if cliNS != "" {
    CurrentNamespace = cliNS
    return nil
  }

  cliCmdCfg, err := clientCmdConfig()
  if err != nil {
    return err
  }

  ns, _, err := cliCmdCfg.Namespace()
  if err != nil {
    return err
  }

  CurrentNamespace = ns
  log.Debug("Kubernetes namespace: ", CurrentNamespace)
  return nil
}

func clientCmdConfig() (clientcmd.ClientConfig, error) {
  configPath, err := kubeConfigPath()
  if err != nil {
    return nil, err
  }

  var cliCmdCfg clientcmd.ClientConfig
  if minikube.Enabled {
    cliCmdCfg = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
      &clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
      &clientcmd.ConfigOverrides{CurrentContext: minikube.K8sContext},
    )
  } else {
    cliCmdCfg = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
      &clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
      &clientcmd.ConfigOverrides{},
    )
  }
  return cliCmdCfg, nil
}

func kubeConfigPath() (string, error) {
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

func GetNamespaceHandler() (typedcorev1.NamespaceInterface, error) {
  kubeConfig, err := KubeConfig()
  if err != nil {
    return nil, err
  }
  coreV1Client, err := typedcorev1.NewForConfig(kubeConfig)
  if err != nil {
    return nil, err
  }
  namespaceHandler := coreV1Client.Namespaces()
  return namespaceHandler, nil
}

func GetSecretHandler(namespace string) (typedcorev1.SecretInterface, error) {
  kubeConfig, err := KubeConfig()
  if err != nil {
    return nil, err
  }
  coreV1Client, err := typedcorev1.NewForConfig(kubeConfig)
  if err != nil {
    return nil, err
  }
  secretHandler := coreV1Client.Secrets(namespace)
  return secretHandler, nil
}

func GetConfigMapHandler(namespace string) (typedcorev1.ConfigMapInterface, error) {
  kubeConfig, err := KubeConfig()
  if err != nil {
    return nil, err
  }
  coreV1Client, err := typedcorev1.NewForConfig(kubeConfig)
  if err != nil {
    return nil, err
  }
  secretHandler := coreV1Client.ConfigMaps(namespace)
  return secretHandler, nil
}

func GetAllNamespaces() (*apicorev1.NamespaceList, error) {
  nsh, err := GetNamespaceHandler()
  if err != nil {
    return nil, err
  }

  return nsh.List(metav1.ListOptions{})
}

func CheckCarbonNamespace(ns string) (bool, error) {
  nsh, err := GetNamespaceHandler()
  if err != nil {
    return false, err
  }

  nsList, err := nsh.List(metav1.ListOptions{})
  if err != nil {
    return false, errors.Wrap(err, "listing namespaces")
  }
  for _, i := range nsList.Items {
    if i.ObjectMeta.Name == ns {
      return true, nil
    }
  }

  return false, nil
}

func CreateGlobalCarbonNamespace() error {
  ns := &apicorev1.Namespace{
    ObjectMeta: metav1.ObjectMeta{
      Name: GlobalCarbonNamespace,
    },
  }

  nsh, err := GetNamespaceHandler()
  if err != nil {
    return err
  }

  _, err = nsh.Create(ns)
  return err
}
