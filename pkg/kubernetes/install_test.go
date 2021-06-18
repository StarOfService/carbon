package kubernetes_test

import (
	"testing"

	"github.com/starofservice/carbon/pkg/kubernetes"
	// "github.com/starofservice/carbon/pkg/test"
)

func TestPurgeList(t *testing.T) {
	// t.Log("Starting Minikube")
	// err := test.MinikubeStart()
	// if err != nil {
	//   t.Error("Failed to start Minikube due ot the error: ", err.Error())
	//   return
	// }

	suite := map[string]bool{
		"rbac.authorization.k8s.io/v1/ClusterRole":              false,
		"apiextensions.k8s.io/v1beta1/CustomResourceDefinition": false,
		"/v1/Namespace":                 false,
		"extensions/v1beta1/Deployment": false,
		"/v1/ServiceAccount":            false,
	}

	allRes, err := kubernetes.GetAllResources()
	if err != nil {
		t.Errorf("Failed to receive all resources list due to the error: %s", err.Error())
		return
	}

	for _, i := range allRes {
		for k := range suite {
			if k == i {
				suite[k] = true
			}
		}
	}

	for k, v := range suite {
		if !v {
			t.Errorf("Missing resource: %s", k)
		}
	}
}
