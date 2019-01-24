package kubernetes

import (
  "fmt"
  "path/filepath"

  restclient "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"

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
  } else {
    return clientcmd.BuildConfigFromFlags("", configPath)
  }
}

func kubeCinfigPath() (string, error) {
  if home := homedir.Path(); home != "" {
    return filepath.Join(home, ".kube", "config"), nil
  } else {
    return "", fmt.Errorf("Unable to discover a user home directory")
  }
}
