package metadata_test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/starofservice/carbon/pkg/docker/metadata"
	"github.com/starofservice/carbon/pkg/test"
)

func TestMain(m *testing.M) {
	log.Info("Starting Docker registry")

	err := test.DockerRegistryStart()
	if err != nil {
		test.DockerRegistryDelete()
		log.Fatal("Failed to start Docker registry due ot the error: ", err.Error())
	}

	code := m.Run()

	test.DockerRegistryDelete()
	os.Exit(code)
}

func TestDockerRegistryCreds(t *testing.T) {
	h, err := metadata.NewDockerMeta(test.DockerMockTestImage)
	if err != nil {
		t.Errorf(err.Error())
	}
	u, p, err := h.GetCredentials()
	if err != nil {
		t.Errorf(err.Error())
	}
	if u != test.DockerMockTestUser {
		t.Errorf("Docker resgistry user doesn't match. Expected: '%s', got: '%s'", test.DockerMockTestUser, u)
	}
	if p != test.DockerMockTestPassword {
		t.Errorf("Docker resgistry password doesn't match. Expected: '%s', got: '%s'", test.DockerMockTestPassword, p)
	}
}

func TestDockerImageLabels(t *testing.T) {
	suite := []string{
		test.DockerMockSrcImg,
		test.DockerMockTestImage,
		"registry:latest",
	}

	for _, i := range suite {
		t.Log("image:", i)

		h, err := metadata.NewDockerMeta(test.DockerMockTestImage)
		if err != nil {
			t.Errorf(err.Error())
		}
		_, err = h.GetLabels()
		if err != nil {
			t.Errorf("Failed to receive image lables for '%s' due to the error: '%s'", i, err.Error())
		}
	}
}
