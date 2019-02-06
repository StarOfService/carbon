package kubernetes

import (
  "github.com/rhysd/go-fakeio"
  log "github.com/sirupsen/logrus"
  cmddelete "k8s.io/kubernetes/pkg/kubectl/cmd/delete"
)

func Delete(manifest, ns string) error {
  log.Debug("Deleting Kubernetes manifests")

  f, ioStreams := KubeCmdFactory(ns)

  cmd := cmddelete.NewCmdDelete(f, ioStreams)

  o := cmddelete.DeleteOptions{IOStreams: ioStreams}
  o.FilenameOptions.Filenames = []string{"-"}
  o.IgnoreNotFound = true

  err := o.Complete(f, []string{}, cmd)
  if err != nil {
    return err
  } 

  log.Trace("Kubernetes manifests for being deleted: ", manifest)

  fake := fakeio.Stdin(manifest)
  defer fake.Restore()
  fake.CloseStdin()

  err = o.RunDelete()
  if err != nil {
    return err
  }

  return nil
}
