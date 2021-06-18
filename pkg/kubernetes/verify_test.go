package kubernetes_test

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/starofservice/carbon/pkg/kubernetes"
)

func TestKubeTemplateVerify(t *testing.T) {
	verifySuites := []struct {
		description string
		data        string
		variables   kubernetes.DepVars
		valid       bool
	}{
		{
			"Basic test of one Package field and one variable",
			"key: {{.Pkg.Name}}, {{.Var.Name}}",
			kubernetes.DepVars{
				Pkg: kubernetes.DepVarsPkg{Name: "foo"},
				Var: map[string]string{"Name": "bar"},
			},
			true,
		},
		{
			"Testing a condition in the template",
			`key: {{if eq .Pkg.Name "foo"}} foo {{else}} bar {{end}}`,
			kubernetes.DepVars{
				Pkg: kubernetes.DepVarsPkg{Name: "foo"},
			},
			true,
		},
		{
			"Testing an excess variable",
			"key: {{.Var.Name}}",
			kubernetes.DepVars{
				Var: map[string]string{
					"Name":        "bar",
					"Environment": "test",
				},
			},
			true,
		},
		{
			"Testing a missing variable",
			"key: {{.Var.Name}}",
			kubernetes.DepVars{
				Var: map[string]string{"Environment": "test"},
			},
			false,
		},
		{
			"Testing invalid YAML",
			"{{.Pkg.Name}}, {{.Var.Name}}",
			kubernetes.DepVars{
				Pkg: kubernetes.DepVarsPkg{Name: "foo"},
				Var: map[string]string{"Name": "bar"},
			},
			false,
		},
	}

	for _, s := range verifySuites {
		t.Log("suite:", s.description)

		if s.valid {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.FatalLevel)
			defer log.SetLevel(log.InfoLevel)
		}

		tmpfile, err := ioutil.TempFile("", "carbon-test-k8s-verify")
		if err != nil {
			t.Errorf(err.Error())
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(s.data)); err != nil {
			t.Errorf(err.Error())
		}
		if err := tmpfile.Close(); err != nil {
			t.Errorf(err.Error())
		}

		kd := &kubernetes.KubeInstall{Variables: s.variables}

		err = kd.VerifyTpl(tmpfile.Name())
		if err != nil && s.valid {
			t.Errorf("Suite data: %v", s.data)
			t.Errorf("Assumed the suite data is valid, but got an error: %v", err.Error())
		} else if err == nil && !s.valid {
			t.Errorf("Suite data: %v", s.data)
			t.Errorf("Assumed the suite data is invalid, but verification is successfully passed")
		}
	}

}
