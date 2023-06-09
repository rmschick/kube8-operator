package deploy

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/client-go/kubernetes"

	v1Controller "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
)

func UpdateCollector(ctx context.Context, clientset *kubernetes.Clientset, oldCR *v1Controller.Collector, currentCR *v1Controller.Collector) error {

	// Create names of resources being deleted which follows the naming convention of the release name in the create.go file
	// {collector-name}-{tenant-instance} ex: cisco-amp-collector-main
	// releaseName := currentCR.Spec.Collector.Name + "-" + currentCR.Spec.Tenant.Instance

	if oldCR.Spec.Collector.Configuration != currentCR.Spec.Collector.Configuration {
		settings := cli.New()
		actionConfig := new(action.Configuration)

		if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "memory", myDebug); err != nil {
			return fmt.Errorf("error initializing action config: %v", err)
		}

		// do something here
	}

	// do some more here

	return nil
}
