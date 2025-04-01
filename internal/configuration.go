package internal

import (
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// Configuration is the amalgamation of various configurations that may be needed.
type Configuration struct {
	Environment string `mapstructure:"environment"`
}

func (c Configuration) Kubeconfig() (*rest.Config, error) {
	var kubeconfig *rest.Config

	var err error

	switch c.Environment {
	case "local":
		home := os.Getenv("HOME")
		kubeConfigFile := home + "/.kube/config"

		kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			return nil, fmt.Errorf("error building kubeconfig: %w", err)
		}
	default:
		kubeconfig, err = rest.InClusterConfig()
		if err != nil {
			return kubeconfig, err
		}
	}

	return kubeconfig, nil
}
