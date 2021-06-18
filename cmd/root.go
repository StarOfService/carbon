package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/starofservice/carbon/pkg/homecfg"
	kubecommon "github.com/starofservice/carbon/pkg/kubernetes/common"
	"github.com/starofservice/carbon/pkg/minikube"
)

var logLevel string
var rootMinikube bool
var namespace string

var RootCmd = &cobra.Command{
	Use:   "carbon",
	Short: "Carbon is a Kubernetes package management command-line utility.",
	Long: `
Carbon is a command-line utlity designed to let you to operate with your application
and related Kubernetes manifests as a signle package.
It uses standard Docker images as a foundation, but adds Kubernetes manifest templates
to Docker image lables. Hence you can use already existing Docker ecosystem in order to
distribute and store your Carbon packages.
More details can be found here: https://github.com/StarOfService/carbon`,
	// Long: `Commands:

	// `
	//  Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		switch logLevel {
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		default:
			log.Fatal("Unsupported log level: ", logLevel)
		}

		if rootMinikube {
			err := minikube.CheckStatus()
			if err != nil {
				log.Fatalf("Failed to verify Minikube status due to the error: %s", err)
			}
			minikube.Enabled = true
			if err := minikube.SetDockerEnv(); err != nil {
				log.Fatalf("Failed to set up Docker environment variables for Minikube due to the error: %s", err)
			}
		}

		kubecommon.SetNamespace(namespace)

		err := homecfg.InitHomeConfig()
		if err != nil {
			log.Fatalf("Failed to initialize Carbon config directory at user home due to the error: %s", err)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set the logging level ('trace'|'debug'|'info'|'warn'|'error'|'fatal') (default 'info')")
	RootCmd.PersistentFlags().BoolVarP(&rootMinikube, "minikube", "m", false, "Use the local Minikube instance instead of remote repositories and Kubernetes clusters. May be useful for local development process. Disabled by default.")
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "If present, defineds a default Kubernetes namespace for installed resources. The behaviour of this parameter depends on the used Carbon scope")
}
