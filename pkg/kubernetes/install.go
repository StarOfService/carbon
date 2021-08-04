package kubernetes

import (
  "bytes"
  "fmt"
  "strings"
  "text/template"

  "github.com/Masterminds/sprig"
  "github.com/pkg/errors"
  "github.com/rhysd/go-fakeio"
  log "github.com/sirupsen/logrus"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/discovery"
  _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
  cmdapply "k8s.io/kubernetes/pkg/kubectl/cmd/apply"

  kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
  "github.com/starofservice/carbon/pkg/util/tojson"
)

func (self *KubeInstall) UpdateVars(vars map[string]string) {
  log.Debug("Applying Carbon variables")

  for k, v := range vars {
    log.Tracef("%s: %s", k, v)
    if _, ok := self.Variables.Var[k]; ok {
      self.Variables.Var[k] = v
    } else {
      log.Debugf("Variable '%s' is defined, but it isn't supported by the current package", k)
    }
  }
}

func (self *KubeInstall) Build() error {
  log.Debug("Building Kubernetes manifest based on the template from the package and provided variables")

  tpl, err := template.New("kubeManifest").Option("missingkey=zero").Funcs(sprig.TxtFuncMap()).Parse(string(self.RawManifest))
  if err != nil {
    return errors.Wrap(err, "parsing Kubernetes manifests teamplate")
  }

  buf := &bytes.Buffer{}
  err = tpl.Execute(buf, self.Variables)
  if err != nil {
    return errors.Wrap(err, "building Kubernetes manifests")
  }

  self.BuiltManifest, err = tojson.ToJSON(buf.Bytes())
  if err != nil {
    return errors.Wrap(err, "converting Kubernetes manifests to JSON")
  }

  return nil
}

func (self *KubeInstall) SetAppLabels() error {
  log.Debug("Applying Carbon lables for Kubernetes manifests")

  var ops string
  if self.Scope == "cluster" {
    ops = self.clusterScopeLabelsPatch()
  } else if self.Scope == "namespace" {
    ops = self.nsScopeLabelsPatch()
  }

  patch, err := tojson.ToJSON([]byte(ops))
  if err != nil {
    log.Error("Most likely it's a bug of the Carbon tool. Please, create an issue for us and provide all possible details.")
    return errors.Wrap(err, "converting Kubernetes patch with Carbon labels to JSON")
  }
  if err := self.ProcessPatches(patch); err != nil {
    return err
  }
  return nil
}

func (self *KubeInstall) clusterScopeLabelsPatch() string {
  ops := fmt.Sprintf(`---
filters:
  kind: .*
type: merge
patch:
  metadata:
    labels:
      managed-by: carbon
      carbon/component-name: %s
      carbon/component-version: %s
`, self.Variables.Pkg.Name, self.Variables.Pkg.Version)
  return ops
}

func (self *KubeInstall) nsScopeLabelsPatch() string {
  ops := fmt.Sprintf(`---
filters:
  kind: .*
type: merge
patch:
  metadata:
    labels:
      managed-by: carbon
      carbon/component-name: %s
      carbon/component-version: %s
      carbon/component-namespace: %s
`, self.Variables.Pkg.Name, self.Variables.Pkg.Version, kubecommon.CurrentNamespace)
  return ops
}

func (self *KubeInstall) Apply(defPWL bool) error {
  log.Debug("Applying Kubernetes manifests")

  f, ioStreams := kubecommon.KubeCmdFactory(kubecommon.CurrentNamespace)

  cmd := cmdapply.NewCmdApply("kubectl", f, ioStreams)

  o := cmdapply.NewApplyOptions(ioStreams)

  err := o.Complete(f, cmd)
  if err != nil {
    return err
  }

  o.DeleteOptions.FilenameOptions.Filenames = []string{"-"}
  o.Prune = true

  force := true
  o.DeleteFlags.Force = &force

  var selector string
  if self.Scope == "cluster" {
    selector = fmt.Sprintf("carbon/component-name=%s", self.Variables.Pkg.Name)
  } else if self.Scope == "namespace" {
    selector = fmt.Sprintf("carbon/component-name=%s,carbon/component-namespace=%s",
                           self.Variables.Pkg.Name, kubecommon.CurrentNamespace)
  }
  o.Selector = selector

  if !defPWL {
    allRes, err := GetAllResources()
    if err == nil {
      for _, i := range allRes {
        o.PruneWhitelist = append(o.PruneWhitelist, i)
      }
    } else {
      log.Warn("I'm unable to discover Kubernetes resources for prune operation. So I'll be using the default prune-whitelist from kubectl apply")
    }
  }

  o.Namespace = kubecommon.CurrentNamespace
  // o.EnforceNamespace = true

  if self.Scope == "cluster" {
    o.EnforceNamespace = false
  } else if self.Scope == "namespace" {
    o.EnforceNamespace = true
  }

  log.Trace("Final Kubernetes manifests for being applied: ", string(self.BuiltManifest))

  fake := fakeio.StdinBytes([]byte{})
  defer fake.Restore()
  go func() {
    fake.StdinBytes(self.BuiltManifest)
    fake.CloseStdin()
  }()

  err = o.Run()
  if err != nil {
    msg := err.Error()
    msg = strings.Replace(msg, `error validating "STDIN": `, "", -1)
    msg = strings.Replace(msg, `; if you choose to ignore these errors, turn validation off with --validate=false`, "", -1)
    return errors.Errorf(msg)
  }

  return nil
}

func GetAllResources() ([]string, error) {

  kubeConfig, err := kubecommon.KubeConfig()
  if err != nil {
    log.Debugf("Failed to discover location of kubeconfig due to the error: %s", err.Error())
    return nil, err
  }

  discClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
  if err != nil {
    log.Debugf("Failed to create Kubernetes discovery handler due to the error: %s", err.Error())
    return nil, err
  }

  apiResList, err := discClient.ServerResources()
  if err != nil {
    log.Debugf("Failed to receive Kubernetes server resources due to the error: %s", err.Error())
    return nil, err
  }

  return procAPIResList(apiResList), nil
}

func procAPIResList (apiResList []*metav1.APIResourceList) []string {
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
  for k := range genResMap {
    resources = append(resources, k)
  }
  // CRDs must be checked at a very last order
  for k := range crdResMap {
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
    "ExternalMetricValueList",
    "LocalSubjectAccessReview",
    "NodeProxyOptions",
    "PodAttachOptions",
    "PodExecOptions",
    "PodPortForwardOptions",
    "PodProxyOptions",
    "PodTemplate",
    "ReplicationControllerDummy",
    "Scale",
    "SelfSubjectAccessReview",
    "SelfSubjectRulesReview",
    "ServiceProxyOptions",
    "SubjectAccessReview",
    "TokenRequest",
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
//     if res == i {
//       return true
//     }
//   }
//   return false
// }
