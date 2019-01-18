package kubernetes

import (
  "fmt"
  // "strings"
  "path/filepath"
  // "os"
  // "k8s.io/cli-runtime/pkg/genericclioptions"
  // cmdapply "k8s.io/kubernetes/pkg/kubectl/cmd/apply"
  // cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
  // "github.com/rhysd/go-fakeio"
  // "k8s.io/client-go/discovery"
  "k8s.io/client-go/tools/clientcmd"
  restclient "k8s.io/client-go/rest"
  // metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  // log "github.com/sirupsen/logrus"
  "github.com/starofservice/carbon/pkg/util/homedir"
)


func GetKubeConfig() (*restclient.Config, error) {
  var err error
  var config *restclient.Config
  if home := homedir.Path(); home != "" {
    config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
    if err != nil {
      return nil, err
    }
  } else {
    return nil, fmt.Errorf("Unable to discover user home directory")
  }
  return config, nil
}
