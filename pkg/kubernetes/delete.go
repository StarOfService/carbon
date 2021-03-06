package kubernetes

import (
  "github.com/rhysd/go-fakeio"
  log "github.com/sirupsen/logrus"
  cmddelete "k8s.io/kubernetes/pkg/kubectl/cmd/delete"
  kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
)

func Delete(manifest, ns string) error {
  log.Debug("Deleting Kubernetes manifests")

  f, ioStreams := kubecommon.KubeCmdFactory(ns)

  cmd := cmddelete.NewCmdDelete(f, ioStreams)

  // https://github.com/kubernetes/kubernetes/blob/v1.13.1/pkg/kubectl/cmd/delete/delete_flags.go#L30-L44
  o := cmddelete.DeleteOptions{IOStreams: ioStreams}
  o.Cascade = true
  o.FilenameOptions.Filenames = []string{"-"}
  o.IgnoreNotFound = true

  err := o.Complete(f, []string{}, cmd)
  if err != nil {
    return err
  }

  log.Trace("Kubernetes manifests for being deleted: ", manifest)

  fake := fakeio.StdinBytes([]byte{})
  defer fake.Restore()
  go func() {
    fake.Stdin(manifest)
    fake.CloseStdin()
  }()

  err = o.RunDelete()
  if err != nil {
    return err
  }

  return nil
}
