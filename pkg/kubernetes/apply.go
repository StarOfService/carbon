package kubernetes

import (
  "fmt"
  "os"
  "strings"

  "github.com/rhysd/go-fakeio"
  log "github.com/sirupsen/logrus"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/cli-runtime/pkg/genericclioptions"
  "k8s.io/client-go/discovery"
  cmdapply "k8s.io/kubernetes/pkg/kubectl/cmd/apply"
  cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"

  "github.com/starofservice/carbon/pkg/minikube"
)

func (self *KubeDeployment) Apply(defPWL bool, ns string) error {
  log.Debug("Applying kubernetes manifests")

  kubeConfigFlags := genericclioptions.NewConfigFlags()

  if minikube.Enabled {
    context := minikube.K8sContext
    kubeConfigFlags.Context = &context
  }

  matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

  f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
  ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

  cmd := cmdapply.NewCmdApply("kubectl", f, ioStreams)
  cmd.Flags().Set("context", "minikube")

  o := cmdapply.NewApplyOptions(ioStreams)

  err := o.Complete(f, cmd)
  if err != nil {
    return err
  } 

  o.DeleteOptions.FilenameOptions.Filenames = []string{"-"}
  o.Prune = true

  force := true
  o.DeleteFlags.Force = &force

  selector := fmt.Sprintf("carbon/component-name=%s", self.Variables.Pkg.Name)
  o.Selector = selector

  if !defPWL {
    allRes, err := getAllResources()
    if err == nil {
      for _, i := range allRes {
        o.PruneWhitelist = append(o.PruneWhitelist, i)
      }
    } else {
      log.Warn("I'm unable to discover kubernetes resources for prune operation. So I'll be using the default prune-whitelist from kubectl apply")
    }
  }

  if ns != "" {
    o.Namespace = ns
    o.EnforceNamespace = true
  }

  log.Trace("Final Kubernetes config for being applied: ", string(self.BuiltManifest))

  fake := fakeio.Stdin(string(self.BuiltManifest))
  defer fake.Restore()
  fake.CloseStdin()

  err = o.Run()
  if err != nil {
    msg := err.Error()
    msg = strings.Replace(msg, `error validating "STDIN": `, "", -1)
    msg = strings.Replace(msg, `; if you choose to ignore these errors, turn validation off with --validate=false`, "", -1)
    return fmt.Errorf(msg)
  }

  return nil
}


func getAllResources() ([]string, error) {

  kubeConfig, err := GetKubeConfig()
  if err != nil {
    log.Debugf("Failed to discover location of kubeconfig due to the error: %s", err.Error())
    return nil, err
  }

  discClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
  if err != nil {
    log.Debugf("Failed to create kubernetes discovery handler due to the error: %s", err.Error())
    return nil, err
  }

  apiResList, err := discClient.ServerResources()
  if err != nil {
    log.Debugf("Failed to receive kubernetes server resources due to the error: %s", err.Error())
    return nil, err
  }

  return procApiResList(apiResList), nil
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
