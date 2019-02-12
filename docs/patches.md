## Patches
Patches are used on a package installation stage in order to apply custom changes, which are not expected by a package developer.

Besides obvious benefits of limitless modifications of Kubernetes manifests, this feature may be used as a key factor for separation of areas of responsibility between developers and infrastructure team. Now developers may omit some infrastructure-specific fields, which will be filled by CI platform during the installation.

Carbon supports two standards of patches:
- Merge Patch, RFC7386
- JSON Patch, RFC6902

Patches can be provided by `--patch` or `--patch-file` flags. Both flags can be used multiple times
`--patch` accepts JSON format
`--patch-file` accepts JSON or YAML format

Any Carbon patch consists of 3 fields:
- `filters`
- `type`
- `patch`

### Filters
Every filter rule is a key-value pair, where:
- *key* is a slash-separated path to a field containing a string as a value
- *value* is a [regular expression](https://github.com/google/re2/wiki/Syntax)

Filters section may contain multiple rules. A patch from the `patch` section is applied only for the resources which match to all filters

### Type
This field can have only one of two values:
- *merge* for RFC7386 merge patch
- *json* for RFC6902 json patch

### Patch
This field contains a patch corresponding to a time from the `type` section.
A comprehensive information for every type can be fond here:
- [Merge Patch, RFC7386](https://tools.ietf.org/html/rfc7386)
- [JSON Patch, RFC6902](https://tools.ietf.org/html/rfc6902)

### Examples
Add 'managed-by' label for all resources:
```
filters:
  kind: .*
type: merge
patch:
  metadata:
    labels:
      managed-by: carbon
```

Remove '/spec/scope' field for a 'CustomResourceDefinition' resource with a name prefixed by 'carbon':
```
{
  "filters": {
    "kind": "CustomResourceDefinition",
    "metadata/name": "carbon.*"
  },
  "type": "json",
  "patch": [
    {
      "op": "remove",
      "path": "/spec/scope"
    }
  ]
}
```
