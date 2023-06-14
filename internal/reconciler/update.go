package reconciler

import (
	"context"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"

	v1Controller "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
)

func UpdateCollector(ctx context.Context, clientset *kubernetes.Clientset, currentCR *v1Controller.Collector) error {

	// Create names of resources being deleted which follows the naming convention of the release name in the create.go file
	// {collector-name}-{tenant-instance} ex: cisco-amp-collector-main
	releaseName := currentCR.Spec.Collector.Name + "-" + currentCR.Spec.Tenant.Instance
	collectorDeployment, err := clientset.AppsV1().Deployments(currentCR.Spec.Tenant.Reference).Get(ctx, releaseName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	collectorChart, err := getCollectorChart(currentCR)
	if err != nil {
		fmt.Printf("could not get collector chart: %v", err)

		return err
	}

	// Unmarshal the values file to use for the helm chart
	vals, err := getValues(currentCR.Spec.Collector.Configuration)
	if err != nil {
		fmt.Printf("could not unmarshal values file: %v", err)

		return err
	}

	settings := cli.New()
	settings.SetNamespace(collectorDeployment.Namespace)
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "memory", myDebug); err != nil {
		return fmt.Errorf("error initializing action config: %v", err)
	}

	updateAction := action.NewUpgrade(actionConfig)
	updateAction.Namespace = collectorDeployment.Namespace
	updateAction.CleanupOnFail = true
	updateAction.ResetValues = true
	updateAction.Version = "latest"

	collectorChart.Values["image"] = map[string]string{
		"repository": collectorDeployment.Spec.Template.Spec.Containers[0].Image,
	}

	collectorChart.Metadata.AppVersion = "latest"

	rendered, err := updateAction.Run(releaseName, collectorChart, vals)
	if err != nil {
		return fmt.Errorf("error updating collector: %v", err)
	}

	// Create a file to write the rendered template to for debugging purposes
	file, err := os.Create("rendered_template_update.yaml")
	if err != nil {
		fmt.Printf("error could not render chart: %v", err)

		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error could not render chart: %v", err)
			os.Exit(1)
		}
	}(file)

	// Writing the contents of the rendered chart configuration
	if _, err := file.WriteString(rendered.Manifest); err != nil {
		fmt.Println(err)

		panic(err)
	}

	fmt.Println("\nRendered template written to file.")

	return nil
}
