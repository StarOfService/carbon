package kubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/starofservice/carbon/pkg/kubernetes"
	kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
)

func TestKubeSetNamespace(t *testing.T) {
	patchSuites := []struct {
		description string
		original    string
		scope       string
		namespace   string
		expected    string
	}{
		{
			`Setting up missing namespace for a resource with namespaced Carbon scope`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
			`namespace`,
			`default2`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test","namespace":"default2"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
		},
		{
			`Overriding namespace for a resource with namespaced Carbon scope`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test","namespace":"default"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
			`namespace`,
			`default2`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test","namespace":"default2"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
		},
		{
			`Setting up missing namespace for a resource with clustered Carbon scope`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
			`cluster`,
			`default2`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test","namespace":"default2"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
		},
		{
			`Making sure namespace is preserved for a resource with clustered Carbon scope`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test","namespace":"default"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
			`cluster`,
			`default2`,
			`{"apiVersion":"v1","kind":"Service","metadata":{"name":"test","namespace":"default"},"spec":{"ports":[{"port":80,"protocol":"TCP","targetPort":80}],"selector":{"app":"test"}}}`,
		},
	}

	for _, s := range patchSuites {
		t.Log("suite:", s.description)

		kubecommon.SetNamespace(s.namespace)
		ki := &kubernetes.KubeInstall{
			BuiltManifest: []byte(s.original),
			// RawManifest: []byte(s.original),
			Scope: s.scope,
			// Variables: DepVars{
			//   Pkg: DepVarsPkg{
			//     Name: "test",
			//     Version: "latest",
			//     DockerName: "name",
			//     DockerTag: "latest",
			//   },
			//   Var: make(map[string]string),
			// },
		}

		err := ki.SetNamespace()
		if err != nil {
			t.Errorf(err.Error())
		}

		assert.Equal(t, s.expected, string(ki.BuiltManifest), "they should be equal")
	}
}
