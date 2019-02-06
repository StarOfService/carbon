package cmd

func AppNamespace(ns string) string {
  if ns != "" {
    return ns
  }
  return "default"
}

func MetadataNamespace(metans, ns string) string {
  if metans != "" {
    return metans
  }
  if ns != "" {
    return ns
  }
  return "default"
}