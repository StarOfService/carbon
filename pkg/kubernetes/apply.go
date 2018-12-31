package kubernetes

import (
  "fmt"
  "strings"
  "path/filepath"
  "os"
  "k8s.io/cli-runtime/pkg/genericclioptions"
  cmdapply "k8s.io/kubernetes/pkg/kubectl/cmd/apply"
  cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
  "github.com/rhysd/go-fakeio"
  "k8s.io/client-go/discovery"
  "k8s.io/client-go/tools/clientcmd"
  restclient "k8s.io/client-go/rest"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  log "github.com/sirupsen/logrus"
)


func (self *KubeDeployment) Apply(defPWL bool) error {
  log.Debug("Applying kubernetes manifests")

  kubeConfigFlags := genericclioptions.NewConfigFlags()
  matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
  f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
  ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

  cmd := cmdapply.NewCmdApply("kubectl", f, ioStreams)
  cmd.Flags().Set("filename", "-")
  cmd.Flags().Set("force", "true")
  cmd.Flags().Set("prune", "true")

  selector := fmt.Sprintf("carbon/component-name=%s", self.Variables.Pkg.Name)
  cmd.Flags().Set("selector", selector)

  if !defPWL {
    allRes, err := getAllResources()
    if err == nil {
      for _, i := range allRes {
        cmd.Flags().Set("prune-whitelist", i)
      }
    } else {
      log.Warn("I'm unable to discover kubernetes resources for prune operation. So I'll be using the default prune-whitelist from kubectl apply")
    }
  }

  fake := fakeio.Stdin(string(self.BuiltManifest))
  defer fake.Restore()
  fake.CloseStdin()
  cmd.Run(cmd, []string{})
  
  // TODO handle failure

  return nil
}


func getAllResources() ([]string, error) {

  kubeConfig, err := getKubeConfig()
  if err != nil {
    // panic(err.Error())
    log.Debugf("Failed to discover location of kubeconfig due to the error: %s", err.Error())
    return nil, err
  }

  discClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
  if err != nil {
    // panic(err.Error())
    log.Debugf("Failed to create kubernetes discovery handler due to the error: %s", err.Error())
    return nil, err
  }

  apiResList, err := discClient.ServerResources()
  if err != nil {
    // panic(err.Error())
    log.Debugf("Failed to receive kubernetes server resources due to the error: %s", err.Error())
    return nil, err
  }

  return procApiResList(apiResList), nil
}


func getKubeConfig() (*restclient.Config, error) {
  var err error
  var config *restclient.Config
  if home := homeDir(); home != "" {
    config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
    if err != nil {
      return nil, err
    }
  } else {
    return nil, fmt.Errorf("Unable to discover user home directory")
  }
  return config, nil
}

func homeDir() string {
  if h := os.Getenv("HOME"); h != "" {
    return h
  }
  return os.Getenv("USERPROFILE") // windows
}


func procApiResList (apiResList []*metav1.APIResourceList) []string {
  crdResMap := make(map[string]struct{})
  genResMap := make(map[string]struct{})
  for _, rl := range apiResList {
    groupVersion := rl.GroupVersion
    if len(strings.Split(groupVersion, "/")) == 1 {
      groupVersion = "/" + groupVersion
    }
    for _, rh := range rl.APIResources {
      if skipResource(rh.Kind) {
        continue
      }

      r := fmt.Sprintf("%s/%s", groupVersion, rh.Kind)
      // if skipVersionedResource(r) {
      //   continue 
      // }

      if rh.Kind == "CustomResourceDefinition" {
        crdResMap[r] = struct{}{}
      } else {
        genResMap[r] = struct{}{}  
      }
    }
  }

  var resources []string
  for k, _ := range genResMap {
    resources = append(resources, k)
  }
  // CRDs must be checked at a very last order
  for k, _ := range crdResMap {
    resources = append(resources, k)
  }
  return resources
}


func skipResource(res string) bool {
  skipRes := []string{
    "AdmissionReview",
    "Binding",
    "ComponentStatus",
    "ControllerRevision",
    "DeploymentRollback",
    "Event",
    "Eviction",
    "LocalSubjectAccessReview",
    "PodTemplate",
    "ReplicationControllerDummy",
    "Scale",
    "SelfSubjectAccessReview",
    "SelfSubjectRulesReview",
    "SubjectAccessReview",
    "TokenReview",
    "VolumeAttachment",
  }
  for _, i := range skipRes {
    if res == i {
      return true
    }
  }
  return false
}


   // It doesn't reduce processing time significantly.
   // So I prefer to keep these versions for compatibility
   // with  old Kubernetes versions
// func skipVersionedResource(res string) bool {
//   skipRes := []string{
//     "extensions/v1beta1/PodSecurityPolicy",
//     "extensions/v1beta1/NetworkPolicy",
//     "extensions/v1beta1/Deployment",
//     "extensions/v1beta1/DaemonSet",
//     "extensions/v1beta1/ReplicaSet",
//     "apps/v1beta1/Deployment",
//     "apps/v1beta2/Deployment",
//     "storage.k8s.io/v1beta1/StorageClass",
//     "apiregistration.k8s.io/v1beta1/APIService",
//     "autoscaling/v2beta1/HorizontalPodAutoscaler",
//     "apps/v1beta2/DaemonSet",
//     "apps/v1beta2/ReplicaSet",
//     "apps/v1beta1/StatefulSet",
//     "apps/v1beta2/StatefulSet",
//   }
//   for _, i := range skipRes {
//     // fmt.Println(rh.Kind, i)
//     if res == i {
//       // fmt.Println("continue")
//       return true
//     }
//   }
//   return false
// }
