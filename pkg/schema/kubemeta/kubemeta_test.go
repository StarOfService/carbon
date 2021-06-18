package kubemeta_test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/starofservice/carbon/pkg/kubernetes"
	kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
	"github.com/starofservice/carbon/pkg/schema/kubemeta"
	// "github.com/starofservice/carbon/pkg/test"
)

const testPkgName = "test-kubemetac-reate"

func TestMain(m *testing.M) {
	// log.Info("Starting Minikube")
	// err := test.MinikubeStart()
	// if err != nil {
	//   log.Error("Failed to start Minikube due ot the error: ", err.Error())
	// }

	err := kubecommon.SetNamespace("")
	if err != nil {
		log.Errorf("Failed to set current namespace due to the error %s", err.Error())
	}

	code := m.Run()

	os.Exit(code)
}

func TestKubemetaCreate(t *testing.T) {
	kms, err := kubemeta.GetAll()
	if err != nil {
		t.Errorf("Failed to get all kubemeta items due to the error %s", err.Error())
		return
	}
	initCount := len(kms)

	ki := &kubernetes.KubeInstall{
		Variables: kubernetes.DepVars{
			Pkg: kubernetes.DepVarsPkg{
				Name: testPkgName,
			},
		},
	}

	km, err := kubemeta.New(ki, []byte{})
	if err != nil {
		t.Errorf("Failed to create a new kubemeta item due to the error %s", err.Error())
		return
	}
	err = km.Apply()
	if err != nil {
		t.Errorf("Failed to save a new kubemeta item due to the error %s", err.Error())
		return
	}

	kms, err = kubemeta.GetAll()
	if err != nil {
		t.Errorf("Failed to get all kubemeta items due to the error %s", err.Error())
		return
	}

	assert.Equal(t, initCount+1, len(kms), "they should be equal")
}

func TestKubemetaIsInstalled(t *testing.T) {
	resp, err := kubemeta.IsInstalled(testPkgName)
	if err != nil {
		t.Errorf("Failed to check kubemeta item due to the error %s", err.Error())
		return
	}
	assert.Equal(t, resp, true, "they should be equal")
}

func TestKubemetaDelete(t *testing.T) {
	kms, err := kubemeta.GetAll()
	if err != nil {
		t.Errorf("Failed to get all kubemeta items due to the error %s", err.Error())
		return
	}
	initCount := len(kms)

	km, err := kubemeta.Get(testPkgName)
	if err != nil {
		t.Errorf("Failed to get kubemeta item due to the error %s", err.Error())
		return
	}

	err = km.Delete()
	if err != nil {
		t.Errorf("Failed to delete kubemeta item due to the error %s", err.Error())
		return
	}

	kms, err = kubemeta.GetAll()
	if err != nil {
		t.Errorf("Failed to get all kubemeta items due to the error %s", err.Error())
		return
	}

	assert.Equal(t, initCount-1, len(kms), "they should be equal")
}
